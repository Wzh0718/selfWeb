package startwork

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"selfWeb/src/tools/Cache"
	"time"
)

func initServer(address string, router *gin.Engine) server {
	// 启动缓存清理器，设定定期清除过期缓存的数据
	Cache.StartCacheCleaner(1 * time.Hour)
	return &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  1800 * time.Second,
		WriteTimeout: 1800 * time.Second,
		//MaxHeaderBytes: 1 << 20,
	}
}
