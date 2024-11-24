package app

import (
	"cc/server/common"
	"cc/server/domain"
	"cc/server/infra"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
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
		log.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	domain.NewConfigEpoll().Notify(config)
	c.JSON(200, gin.H{"message": "success"})
}

func GetConfig(c *gin.Context) {
	//获取get请求中的key参数
	key := c.Query("key")
	config, err := infra.GetConfig(key)
	resp := common.GetConfigResponse{
		Message: "success",
	}
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			resp.Message = "not found"
			c.AbortWithStatusJSON(404, resp)
			return
		}
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	callback := c.Query("callback")
	domain.NewConfigEpoll().Add(key, &domain.Listener{
		CallbackURL: callback,
	})
	resp.Config = config
	c.JSON(200, resp)
}
