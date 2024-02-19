package wire

import "manager/internal/handler"

type K8sHandler struct {
}

func NewK8sHandler() *K8sHandler {
	return &K8sHandler{}
}

func (h *K8sHandler) DepHandlers() *handler.DepHandler {
	return &handler.DepHandler{}
}

func (h *K8sHandler) PodHandler() *handler.PodHandler {
	return &handler.PodHandler{}
}

func (h *K8sHandler) NsHandler() *handler.NsHandler {
	return &handler.NsHandler{}
}

func (h *K8sHandler) EventHandler() *handler.EventHandler {
	return &handler.EventHandler{}
}

func (h *K8sHandler) IngressHandler() *handler.IngressHandler { return &handler.IngressHandler{} }

func (h *K8sHandler) SvcHandler() *handler.ServiceHandler {
	return &handler.ServiceHandler{}
}

func (h *K8sHandler) ResourceHandler() *handler.ResourceHandler {
	return &handler.ResourceHandler{}
}
