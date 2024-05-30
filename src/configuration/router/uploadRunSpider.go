package router

import (
	"github.com/gin-gonic/gin"
	"selfWeb/src/main/service/smt"
)

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

func DownloadSpider(response *gin.RouterGroup) {
	response.GET("/download/runSpider", func(context *gin.Context) {
		smt.DownLoadSpiderFile(context)
	})
}
