package startwork

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"selfWeb/src/configuration"
	"selfWeb/src/configuration/structs"
	"selfWeb/src/tools/TimeUtils"
)

// ChoseConfig 选择配置文件，自动加载Root路径
func ChoseConfig(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "chose config file")
		flag.Parse()
		if config == "" {
			if configEnv := os.Getenv(structs.ConfigEnv); configEnv == "" {
				//默认只走一个文件，可以后期修改 路径为resource --> yaml --> xxx.yaml
				switch gin.Mode() {
				case gin.DebugMode:
				case gin.TestMode:
				case gin.ReleaseMode:
					config = structs.ConfigFile
					fmt.Printf("正在使用gin模式的%s环境名称, config路径为%s\n", gin.Mode(), config)
				}
			} else {
				config = configEnv
				fmt.Printf("正在使用gin模式的%s环境名称, config路径为%s\n", gin.Mode(), config)
			}
		} else {
			fmt.Printf("正在使用命令行的-c参数传递的值, config的路径为%s\n", config)
		}

	} else {
		config = path[0]
		fmt.Printf("正在使用func Viper()传递的值, config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&configuration.FileConfig); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&configuration.FileConfig); err != nil {
		panic(err)
	}
	// 更新当前root路径
	configuration.FileConfig.AutoCode.Root, _ = filepath.Abs("../..")

	// 设置本地缓存时间
	InitCache()

	return v
}

func InitCache() {

	expiresTime, err := TimeUtils.ParseDuration(configuration.FileConfig.JWT.ExpiresTime)

	if err != nil {
		panic(err)
	}

	configuration.Cache = local_cache.NewCache(
		local_cache.SetDefaultExpire(expiresTime),
	)

}
