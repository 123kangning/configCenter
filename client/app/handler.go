package app

import (
	"cc/client/common"
	"cc/client/infra"
	"github.com/gin-gonic/gin"
)

func UpdateConfig(c *gin.Context) {
	config := &common.Config{}
	err := c.ShouldBind(config)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	err = infra.UpdateConfig(config)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func GetConfig(c *gin.Context) {
	//获取get请求中的key参数
	key := c.Query("key")
	config, err := infra.GetConfig(key)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, common.GetConfigResponse{
		Message: "success",
		Config:  config,
	})
}

func Callback(c *gin.Context) {
	config := &common.Config{}
	err := c.ShouldBind(config)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	err = infra.UpdateExistConfig(config)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
