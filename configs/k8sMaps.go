package configs

import (
	"k8s.io/client-go/kubernetes"
	"view/internal/maps"
)

type K8sMaps struct {
	K8sClient *kubernetes.Clientset `inject:"-"`
}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

func (this *K8sMaps) InitDepMap() *maps.DeploymentMap {
	return &maps.DeploymentMap{}
}

func (this *K8sMaps) InitPodMap() *maps.PodMap {
	return &maps.PodMap{}
}

func (this *K8sMaps) InitNsMap() *maps.NsMapStruct {
	return &maps.NsMapStruct{}
}
