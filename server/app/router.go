package app

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	r := gin.New()

	r.PUT("/config", UpdateConfig)
	r.GET("/config", GetConfig)
	return r
}
