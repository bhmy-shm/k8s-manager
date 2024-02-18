package maps

import (
	"fmt"
	"k8s.io/api/networking/v1"
	"sort"
	"sync"
)

type IngressMap struct {
	data *sync.Map // [ns string] []*v1.Ingress
}

func NewIngressMap() *IngressMap {
	return &IngressMap{
		data: new(sync.Map),
	}
}

func (ig *IngressMap) Get(ns string, name string) *v1.Ingress {
	if items, ok := ig.data.Load(ns); ok {
		for _, item := range items.([]*v1.Ingress) {
			if item.Name == name {
				return item
			}
		}
	}
	return nil
}

func (ig *IngressMap) Add(ingress *v1.Ingress) {
	if list, ok := ig.data.Load(ingress.Namespace); ok {
		list = append(list.([]*v1.Ingress), ingress)
		ig.data.Store(ingress.Namespace, list)
	} else {
		ig.data.Store(ingress.Namespace, []*v1.Ingress{ingress})
	}
}

func (ig *IngressMap) Update(ingress *v1.Ingress) error {
	if list, ok := ig.data.Load(ingress.Namespace); ok {
		for i, rangePod := range list.([]*v1.Ingress) {
			if rangePod.Name == ingress.Name {
				list.([]*v1.Ingress)[i] = ingress
			}
		}
		return nil
	}
	return fmt.Errorf("ingress-%s not found", ingress.Name)
}

func (ig *IngressMap) Delete(ingress *v1.Ingress) {
	if list, ok := ig.data.Load(ingress.Namespace); ok {
		for i, rangeIngress := range list.([]*v1.Ingress) {
			if rangeIngress.Name == ingress.Name {
				newList := append(list.([]*v1.Ingress)[:i], list.([]*v1.Ingress)[i+1:]...)
				ig.data.Store(ingress.Namespace, newList)
				break
			}
		}
	}
}

func (ig *IngressMap) ListByNs(ns string) []*v1.Ingress {

	if list, ok := ig.data.Load(ns); ok {

		newList := list.([]*v1.Ingress)
		sort.Slice(newList, func(i, j int) bool {
			return newList[i].CreationTimestamp.Time.After(newList[j].CreationTimestamp.Time)
		}) //  按时间倒排序

		return newList
	}
	return []*v1.Ingress{} //返回空列表
}
