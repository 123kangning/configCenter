package domain

import (
	"bytes"
	"cc/server/common"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

// Observer 观察者接口
type Observer interface {
	Update(config *common.Config)
}

// Observable 被观察者接口
type Observable interface {
	Add(observer Observer)
	Del(observer Observer)
	Notify(config *common.Config)
}

// ConfigEpoll 多个客户端观察者监听一个配置
type ConfigEpoll struct {
	ToKeys sync.Map //TODO 并发控制之后加一下
}

var configEpoll *ConfigEpoll

func init() {
	configEpoll = &ConfigEpoll{
		ToKeys: sync.Map{},
	}
}

// NewConfigEpoll 饿汉单例模式
func NewConfigEpoll() *ConfigEpoll {
	return configEpoll
}

func (c *ConfigEpoll) Add(key string, observer *Listener) {
	if _, ok := c.ToKeys.Load(key); !ok {
		c.ToKeys.Store(key, &Config{
			Key: key,
		})
	}
	config, _ := c.ToKeys.Load(key)
	config.(*Config).Add(observer)
	c.ToKeys.Store(key, config)
}

func (c *ConfigEpoll) Notify(config *common.Config) {
	if _, ok := c.ToKeys.Load(config.Key); !ok {
		return
	}
	c1, _ := c.ToKeys.Load(config.Key)
	c1.(*Config).Notify(config)
}

// Config 配置，即被观察者
type Config struct {
	Key       string
	observers []Observer
}

// Add 添加观察者
func (c *Config) Add(observer Observer) {
	c.observers = append(c.observers, observer)
}

// Del 删除观察者
func (c *Config) Del(observer Observer) {
	for i, v := range c.observers {
		if v == observer {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

// Notify 通知观察者
func (c *Config) Notify(config *common.Config) {
	for _, v := range c.observers {
		v.Update(config)
	}
}

type Listener struct {
	CallbackURL string
}

func (l *Listener) Update(config *common.Config) {
	// 发送请求到CallbackURL
	log.Println("发送请求到CallbackURL:", l.CallbackURL)
	requestData, _ := json.Marshal(config)
	req, err := http.NewRequest("PUT", l.CallbackURL, bytes.NewBuffer(requestData))
	if err != nil {
		log.Println("Failed to create request: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request: ", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body: ", err)
		return
	}

	var answer map[string]interface{}
	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Println("Failed to unmarshal response body: ", err)
		return
	}
	if answer["message"] != "success" {
		log.Println("Failed to callback update config: ", answer)
		return
	}
	return
}
