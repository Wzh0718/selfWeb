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

	configuration.Logger.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`Libre Web is Starting`, address)
	configuration.Logger.Error(s.ListenAndServe().Error())

}
