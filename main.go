package main

import (
	"github.com/bhmy-shm/gofks"
	"manager/internal/controllers"
	"manager/internal/middlewares"
	"manager/wire"
)

func main() {

	gofks.Ignite("/v1", middlewares.OnRequest()).
		WireApply(
			wire.NewK8sHandler(),
			wire.NewK8sConfig(),
			wire.NewK8sMaps(),
			wire.NewServiceConfig(),
			wire.NewServiceWire(),
		).
		Mount(
			controllers.NewDeploymentCtl(),
			controllers.NewNamespaceCtl(),
			controllers.NewPodCtl(),
			controllers.NewSvcCtl(),
			controllers.NewIngressCtl(),
			controllers.NewWsCtl(),
		).
		Launch()
}
