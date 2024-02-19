package wire

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"manager/svc"
)

type ServiceWire struct {
	ctx *svc.ServiceContext
}

func NewServiceWire(c *gofkConf.Config) *ServiceWire {
	return &ServiceWire{
		ctx: svc.NewServiceContext(c),
	}
}

func (w *ServiceWire) Context() *svc.ServiceContext {
	return w.ctx
}
