package infra

import (
	"cc/client/common"
	"cc/client/configs"
	"fmt"
	"log"
	"os"
)

func UpdateConfig(config *common.Config) error {
	err := updateConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func UpdateExistConfig(config *common.Config) error {
	fmt.Println(os.Getwd())
	file, err := os.OpenFile(fmt.Sprintf("%s%s", configs.DataPrefix, config.Key), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(config.Value))
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(key string) (*common.Config, error) {
	config, err := getLocalConfig(key)
	if err == nil {
		log.Println("从本地找到了配置")
		return config, nil
	}
	log.Println("尝试从远程找到配置")
	config, err = getRemoteConfig(key)
	if err == nil {
		err = UpdateExistConfig(config)
	}
	return config, err
}

func getLocalConfig(key string) (*common.Config, error) {
	file, err := os.OpenFile(fmt.Sprintf("%s%s", configs.DataPrefix, key), os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	//将不确定大小的文件内容读取到内存中
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	value := make([]byte, fileInfo.Size())
	_, err = file.Read(value)
	if err != nil {
		return nil, err
	}
	config := &common.Config{
		Key:   key,
		Value: string(value),
	}
	return config, nil
}
