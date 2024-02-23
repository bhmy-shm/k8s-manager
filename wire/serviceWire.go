package wire

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"k8s.io/client-go/kubernetes"
	"manager/svc"
)

type ServiceWire struct {
	ctx                   *svc.ServiceContext
	*kubernetes.Clientset `inject:"-"`
}

func NewServiceWire(c *gofkConf.Config) *ServiceWire {
	return &ServiceWire{
		ctx: svc.NewServiceContext(c),
	}
}

func (w *ServiceWire) Context() *svc.ServiceContext {
	return w.ctx
}
