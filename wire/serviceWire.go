package wire

import (
	"manager/svc"
)

type ServiceWire struct {
	ctx *svc.ServiceContext
}

func NewServiceWire() *ServiceWire {
	return &ServiceWire{
		ctx: svc.NewServiceContext(),
	}
}

func (w *ServiceWire) Context() *svc.ServiceContext {
	return w.ctx
}
