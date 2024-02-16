package wire

import (
	"manager/internal/maps"
)

type K8sMaps struct {
	dm *maps.DeploymentMap
	pm *maps.PodMap
	nm *maps.NsMapStruct
}

func newK8sMaps() *K8sMaps {
	return &K8sMaps{
		dm: &maps.DeploymentMap{},
		pm: &maps.PodMap{},
		nm: &maps.NsMapStruct{},
	}
}

func (this *K8sMaps) InitDepMap() *maps.DeploymentMap {
	return this.dm
}

func (this *K8sMaps) InitPodMap() *maps.PodMap {
	return this.pm
}

func (this *K8sMaps) InitNsMap() *maps.NsMapStruct {
	return this.nm
}
