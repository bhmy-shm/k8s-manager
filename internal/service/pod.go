package service

import "manager/internal/maps"

type PodService struct {
	PodMap *maps.PodMap `inject:"-"`
}

func Pod() *PodService {
	return &PodService{}
}

func (this *PodService) ListByNs(ns string) interface{} {
	return this.PodMap.ListByNs(ns)
}
