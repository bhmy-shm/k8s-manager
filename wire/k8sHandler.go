package wire

import "manager/internal/handler"

type K8sHandler struct {
	dh *handler.DepHandler
	ph *handler.PodHandler
	nh *handler.NsHandler
}

func newK8sHandler() *K8sHandler {
	return &K8sHandler{}
}

func (this *K8sHandler) DepHandlers() *handler.DepHandler {
	return this.dh
}

func (this *K8sHandler) PodHandler() *handler.PodHandler {
	return this.ph
}

func (this *K8sHandler) NsHandler() *handler.NsHandler {
	return this.nh
}
