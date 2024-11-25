package main

import (
	"cc/server/app"
	"cc/server/configs"
	"fmt"
	"os"
)

func main() {
	r := app.InitRouter()
	err := InitDir(configs.DataPrefix)
	if err != nil {
		panic(err)
	}
	err = r.Run(":9000")
	if err != nil {
		panic(err)
	}
}

func InitDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.Mkdir(dirName, 0755) // 0755是目录权限设置
		if err != nil {
			// 创建目录失败
			fmt.Printf("Failed to create directory: %s\n", err)
			return err
		}
	}
	return nil
}
