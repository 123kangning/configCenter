package main

import (
	"cc/client/app"
	"cc/client/configs"
	"os"
)

func main() {
	r := app.InitRouter()
	err := os.RemoveAll(configs.DataPrefix)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(configs.DataPrefix, 0777)
	if err != nil {
		panic(err)
	}
	err = r.Run(":9001")
	if err != nil {
		panic(err)
	}
}
