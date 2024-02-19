package wire

import (
	"manager/internal/maps"
)

type K8sMaps struct{}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

func (ks *K8sMaps) InitDepMap() *maps.DeploymentMap {
	return &maps.DeploymentMap{}
}

func (ks *K8sMaps) InitPodMap() *maps.PodMap {
	return &maps.PodMap{}
}

func (ks *K8sMaps) InitNsMap() *maps.NamespaceMap {
	return &maps.NamespaceMap{}
}

func (ks *K8sMaps) InitEventMap() *maps.EventMap {
	return &maps.EventMap{}
}

func (ks *K8sMaps) InitIngressMap() *maps.IngressMap {
	return &maps.IngressMap{}
}

func (ks *K8sMaps) InitSvcMap() *maps.ServiceMap {
	return &maps.ServiceMap{}
}

func (ks *K8sMaps) InitResourceMap() *maps.MetaResMapper {
	return maps.NewMetaRespMapper()
}
