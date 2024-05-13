package configuration

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"selfWeb/src/configuration/structs"
)

var (
	FileConfig structs.Configuration
	FileVp     *viper.Viper
	Cache      local_cache.Cache
	Logger     *zap.Logger
)
