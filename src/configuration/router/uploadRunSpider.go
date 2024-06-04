package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"selfWeb/src/configuration"
	"selfWeb/src/main/service/smt"
	"selfWeb/src/tools/HttpsUtil/response"
	"time"
)

var queue = make(chan int, 1)

func PostRunSpider(response *gin.RouterGroup) {
	response.POST("/upload/runSpider", func(context *gin.Context) {
		smt.DownloadSpiderFile(context)
	})
}

func GetRunSpiderVersion(response *gin.RouterGroup) {
	response.GET("/version/runSpider", func(context *gin.Context) {
		smt.GetSpiderVersion(context)
	})
}

func DownloadSpider(r *gin.RouterGroup) {
	r.GET("/download/runSpider", func(c *gin.Context) {
		// 队列超时5秒
		timeout, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancelFunc()
		select {
		case queue <- 1:
			defer func() {
				<-queue
			}()
			configuration.Logger.Info("Server is done")
			smt.DownLoadSpiderFile(c)
		case <-timeout.Done():
			response.ReturnContextNoBody(c, http.StatusRequestTimeout, "Server is busy. Please try again later.")
			configuration.Logger.Info("Server is busy. Please try again later.")
		}

	})
}
