package LoggerUtils

import (
	"go.uber.org/zap/zapcore"
	"os"
	"selfWeb/src/configuration"
	"selfWeb/src/tools/LoggerUtils/Cutter"
)

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct{}

// GetWriteSyncer 获取 zapcore.WriteSyncer, 写入本地日志并输出
func (r *fileRotatelogs) GetWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriter := Cutter.NewCutter(configuration.FileConfig.Zap.Director, level, Cutter.WithCutterFormat("2006-01-02"))
	if configuration.FileConfig.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}
	return zapcore.AddSync(fileWriter)
}
