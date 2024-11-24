package infra

import (
	"cc/server/common"
	"cc/server/configs"
	"fmt"
	"os"
)

func UpdateConfig(config *common.Config) error {
	key := config.Key
	newValue := []byte(config.Value)
	path := fmt.Sprintf("%s%s", configs.DataPrefix, key)
	fmt.Println(os.Getwd())
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	_, err = file.Write(newValue)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(key string) (*common.Config, error) {
	config, err := getLocalConfig(key)
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
