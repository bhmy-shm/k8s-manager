package main

import (
	"fmt"
	"github.com/bhmy-shm/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"k8s.io/apimachinery/pkg/api/meta"
	"manager/internal/controllers"
	"manager/internal/handler"
	"manager/internal/middlewares"
	"manager/wire"
)

func main() {

	conf := gofkConf.Load()

	go func() {
		for {
			get_obj, _ := handler.RateLimitQue.Get()
			if rm_obj, ok := get_obj.(*handler.RateLimitResource); ok {
				fmt.Print("类型是:", rm_obj.Type)
				obj, err := meta.Accessor(rm_obj.Resource)
				if err == nil {
					fmt.Printf(" 资源是:%s/%s \n", obj.GetNamespace(),
						obj.GetName())
				}
			}
		}
	}()

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
