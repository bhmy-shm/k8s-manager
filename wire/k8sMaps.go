package wire

import (
	"manager/internal/maps"
)

type K8sMaps struct{}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

func (ks *K8sMaps) InitDepMap() *maps.DeploymentMap {
	return maps.NewDeploymentMap()
}

func (ks *K8sMaps) InitPodMap() *maps.PodMap {
	return maps.NewPodMap()
}

func (ks *K8sMaps) InitNsMap() *maps.NamespaceMap {
	return maps.NewNamespaceMap()
}

func (ks *K8sMaps) InitEventMap() *maps.EventMap {
	return maps.NewEventMap()
}

func (ks *K8sMaps) InitIngressMap() *maps.IngressMap {
	return maps.NewIngressMap()
}

func (ks *K8sMaps) InitSvcMap() *maps.ServiceMap {
	return maps.NewServiceMap()
}
