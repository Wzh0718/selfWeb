package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetChatGpt 跳转到自己服务的ChatGpt中
func GetChatGpt(response *gin.RouterGroup) {
	response.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://47.119.161.182:7000/#/auth")
	})
}

func GetChatGptDepositScreen(response *gin.RouterGroup) {
	response.GET("/apikey", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://apikey.checo.icu")
	})
}
