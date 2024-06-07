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

var GetSpiderVersionQueue = make(chan int, 100)
var queue = make(chan int, 1)

func PostRunSpider(response *gin.RouterGroup) {
	response.POST("/upload/runSpider", func(context *gin.Context) {
		smt.DownloadSpiderFile(context)
	})
}

func GetRunSpiderVersion(r *gin.RouterGroup) {
	r.GET("/version/runSpider", func(c *gin.Context) {
		timeout, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancelFunc()
		select {
		case GetSpiderVersionQueue <- 1:
			defer func() {
				<-GetSpiderVersionQueue
			}()
			configuration.Logger.Info("GetRunSpiderVersion is done")
			smt.GetSpiderVersion(c)
		case <-timeout.Done():
			configuration.Logger.Info("GetRunSpiderVersion MaxConnection Wait Connect TimeOut...")
			response.ReturnContextNoBody(c, http.StatusRequestTimeout, "GetRunSpiderVersion MaxConnection Wait Connect TimeOut...")
		}

	})
}

func DownloadSpider(r *gin.RouterGroup) {
	r.GET("/download/runSpider", func(c *gin.Context) {
		// 队列超时3秒
		timeout, cancelFunc := context.WithTimeout(c.Request.Context(), 3*time.Second)
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
