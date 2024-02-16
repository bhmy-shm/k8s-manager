package wire

import "manager/svc"

//注入顺序必须固定：
//wire.NewK8sHandler()
//wire.NewK8sConfig()
//wire.NewK8sMaps()
//wire.NewServiceSvc()

type ServiceWire struct {
	ctx     *svc.ServiceContext
	handler *K8sHandler
	maps    *K8sMaps
	conf    *K8sConfig
}

func NewServiceWire() *ServiceWire {
	return &ServiceWire{
		handler: newK8sHandler(),
		conf:    newK8sConfig(),
		maps:    newK8sMaps(),
		ctx:     svc.NewServiceContext(),
	}
}

func (w *ServiceWire) Handler() *K8sHandler {
	return w.handler
}

func (w *ServiceWire) Maps() *K8sMaps {
	return w.maps
}

func (w *ServiceWire) Conf() *K8sConfig {
	return w.conf
}

func (w *ServiceWire) Context() *svc.ServiceContext {
	return w.ctx
}
