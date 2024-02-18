package maps

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"sync"
)

type ServiceMap struct {
	data *sync.Map // [ns string] []*v1.Service
}

func NewServiceMap() *ServiceMap {
	return &ServiceMap{
		data: new(sync.Map),
	}
}

func (s *ServiceMap) Get(ns string, name string) *corev1.Service {
	if items, ok := s.data.Load(ns); ok {
		for _, item := range items.([]*corev1.Service) {
			if item.Name == name {
				return item
			}
		}
	}
	return nil
}

func (s *ServiceMap) Add(svc *corev1.Service) {
	if list, ok := s.data.Load(svc.Namespace); ok {
		list = append(list.([]*corev1.Service), svc)
		s.data.Store(svc.Namespace, list)
	} else {
		s.data.Store(svc.Namespace, []*corev1.Service{svc})
	}
}

func (s *ServiceMap) Update(svc *corev1.Service) error {
	if list, ok := s.data.Load(svc.Namespace); ok {
		for i, range_pod := range list.([]*corev1.Service) {
			if range_pod.Name == svc.Name {
				list.([]*corev1.Service)[i] = svc
			}
		}
		return nil
	}
	return fmt.Errorf("service-%s not found", svc.Name)
}

func (s *ServiceMap) Delete(svc *corev1.Service) {
	if list, ok := s.data.Load(svc.Namespace); ok {
		for i, range_svc := range list.([]*corev1.Service) {
			if range_svc.Name == svc.Name {
				newList := append(list.([]*corev1.Service)[:i], list.([]*corev1.Service)[i+1:]...)
				s.data.Store(svc.Namespace, newList)
				break
			}
		}
	}
}

func (s *ServiceMap) ListAll(ns string) []*corev1.Service {
	if list, ok := s.data.Load(ns); ok {
		newList := list.([]*corev1.Service)
		sort.Slice(newList, func(i, j int) bool {
			return newList[i].CreationTimestamp.Time.After(newList[j].CreationTimestamp.Time)
		}) //  按时间倒排序1
		return newList
	}

	//返回空列表
	return []*corev1.Service{}
}
