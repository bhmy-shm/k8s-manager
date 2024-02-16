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
			wire.NewServiceWire(),
		).
		Mount(
			controllers.NewDeploymentCtl(),
		).
		Launch()
}
