package main

import (
	"go.uber.org/zap"
	"selfWeb/src/configuration"
	"selfWeb/src/startwork"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	//加载配置文件
	configuration.FileVp = startwork.ChoseConfig()
	// 加载日志
	configuration.Logger = startwork.Logger()
	zap.ReplaceGlobals(configuration.Logger)
	startwork.RunServer()
}
