package main

import (
	"cc/server/app"
)

func main() {
	err := app.InitRouter().Run(":9000")
	if err != nil {
		panic(err)
	}
}
