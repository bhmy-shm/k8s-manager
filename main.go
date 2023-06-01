package main

import (
	"github.com/bhmy-shm/gofks/gofk"
	"view/configs"
	"view/internal/controller"
	"view/internal/middlewares"
)

func main() {

	gofk.Ignite(
		middlewares.OnRequest(), //全局中间件
	).
		Config(
			configs.NewK8sHandler(),    //1
			configs.NewK8sConfig(),     //2
			configs.NewK8sMaps(),       //3
			configs.NewServiceConfig(), //4
		).
		Mount(
			//controllers.NewDynamic(), //todo
			controllers.NewWsCtl(),
			controllers.NewDeploymentCtl(),
			controllers.NewPodCtl(),
			controllers.NewNsCtl(),
		).
		Watcher().Launch()
}
