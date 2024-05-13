package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"selfWeb/src/configuration"
	"selfWeb/src/tools/HttpsUtil/middleware"
	"time"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery())
	Router.Use(gin.Logger())

	//启动http
	Router.StaticFS(configuration.FileConfig.Local.StorePath, justFilesFilesystem{http.Dir(configuration.FileConfig.Local.StorePath)})
	configuration.Logger.Info("使用 Http配置, ", zap.String("File Path is ", configuration.FileConfig.Local.StorePath))

	//指定规则
	Router.Use(middleware.CorsByRules())
	configuration.Logger.Info("使用 middleware.CorsByRules()配置文件")

	//统一规范路由组前缀: 主路由层 -->
	RouterGroup := Router.Group(configuration.FileConfig.System.RouterPrefix)
	ChatRouterGroup := RouterGroup.Group("/chat")
	GetChatGpt(ChatRouterGroup)
	GetChatGptDepositScreen(ChatRouterGroup)
	{
		// 监测存活状态
		RouterGroup.GET("/heath", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"time": time.Now(),
				"path": configuration.FileConfig.AutoCode.Root,
			})
		})
	}

	configuration.Logger.Info("router register success")
	return Router
}
