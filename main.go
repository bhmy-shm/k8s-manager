package main

import (
	"github.com/bhmy-shm/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"manager/internal/controllers"
	"manager/internal/middlewares"
	"manager/wire"
)

func main() {

	conf := gofkConf.Load()

	gofks.Ignite("/v1", middlewares.OnRequest()).
		WireApply(
			wire.NewK8sMaps(),
			wire.NewK8sHandler(),
			wire.NewK8sConfig(),
			wire.NewServiceConfig(),
			wire.NewServiceWire(conf),
		).
		Mount(
			controllers.NewDeploymentCtl(),
			controllers.NewNamespaceCtl(),
			controllers.NewPodCtl(),
			controllers.NewSvcCtl(),
			controllers.NewIngressCtl(),
			controllers.NewWsCtl(),
			controllers.NewResourceCtl(),
		).
		Launch()
}
