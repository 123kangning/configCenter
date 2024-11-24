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
	mutex sync.RWMutex
	//涉及到先读出，再修改，再写入的操作，得加锁，sync.Map不能满足需求
	ToKeys map[string]*Config
	// 多个客户端同时修改一个配置，回调的过程中可能会造成多个客户端配置不一致的问题
	//之后加上一个序列号，让所有配置的修改串行化，客户端可以通过序列号判断是否是最新的配置
}

var configEpoll *ConfigEpoll

func init() {
	configEpoll = &ConfigEpoll{
		ToKeys: make(map[string]*Config),
	}
}

// NewConfigEpoll 饿汉单例模式
func NewConfigEpoll() *ConfigEpoll {
	return configEpoll
}

func (c *ConfigEpoll) Add(key string, observer *Listener) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.ToKeys[key]; !ok {
		c.ToKeys[key] = &Config{}
	}
	c.ToKeys[key].Add(observer)
}

func (c *ConfigEpoll) Notify(config *common.Config) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if _, ok := c.ToKeys[config.Key]; !ok {
		return
	}
	c.ToKeys[config.Key].Notify(config)
}

// Config 配置，即被观察者
type Config struct {
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
