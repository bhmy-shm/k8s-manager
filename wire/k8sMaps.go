package wire

import (
	"manager/internal/maps"
)

type K8sMaps struct{}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

func (this *K8sMaps) InitDepMap() *maps.DeploymentMap {
	return &maps.DeploymentMap{}
}

func (this *K8sMaps) InitPodMap() *maps.PodMap {
	return &maps.PodMap{}
}

func (this *K8sMaps) InitNsMap() *maps.NamespaceMap {
	return &maps.NamespaceMap{}
}

func (this *K8sMaps) InitEventMap() *maps.EventMap {
	return &maps.EventMap{}
}
