package startwork

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"selfWeb/src/configuration"
	"selfWeb/src/tools/FileUtils"
	"selfWeb/src/tools/LoggerUtils"
)

func Logger() (logger *zap.Logger) {
	//文件Dir不存在
	if existsPath, _ := FileUtils.PathExists(configuration.FileConfig.Zap.Director); !existsPath {
		fmt.Printf("create %v directory\n", configuration.FileConfig.Zap.Director)
		_ = os.Mkdir(configuration.FileConfig.Zap.Director, os.ModePerm)
	}

	cores := LoggerUtils.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if configuration.FileConfig.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
