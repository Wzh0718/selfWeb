package startwork

import (
	"fmt"
	"go.uber.org/zap"
	"selfWeb/src/configuration"
	"selfWeb/src/configuration/router"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	Router := router.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", configuration.FileConfig.System.Addr)

	s := initServer(address, Router)

	configuration.Logger.Info("Web server run success on ", zap.String("address", address))

	configuration.Logger.Error(s.ListenAndServe().Error())

}
