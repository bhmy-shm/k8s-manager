package maps

import (
	corev1 "k8s.io/api/core/v1"
	"manager/model"
	"sort"
	"sync"
)

type NamespaceMap struct {
	data *sync.Map // [key string] []*corev1.Namespace    key=>namespace的名称
}

func NewNamespaceMap() *NamespaceMap {
	return &NamespaceMap{
		data: new(sync.Map),
	}
}

func (n *NamespaceMap) Get(ns string) *corev1.Namespace {
	if item, ok := n.data.Load(ns); ok {
		return item.(*corev1.Namespace)
	}
	return nil
}

func (n *NamespaceMap) Add(ns *corev1.Namespace) {
	n.data.Store(ns.Name, ns)
}

func (n *NamespaceMap) Update(ns *corev1.Namespace) {
	n.data.Store(ns.Name, ns)
}

func (n *NamespaceMap) Delete(ns *corev1.Namespace) {
	n.data.Delete(ns.Name)
}

// ListAll 显示所有的 namespace
func (n *NamespaceMap) ListAll() *model.NamespaceModel {
	ret := make([]*model.NsModel, 0)

	total, items := convertToMapItems(n.data)
	sort.Sort(items)
	for _, item := range items {
		ret = append(ret, &model.NsModel{Name: item.key})
	}

	return &model.NamespaceModel{
		Total: total,
		List:  ret,
	}
}
