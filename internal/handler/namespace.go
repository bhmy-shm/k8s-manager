package handler

import (
	corev1 "k8s.io/api/core/v1"
	"manager/internal/maps"
)

type NsHandler struct {
	NsMap *maps.NamespaceMap `inject:"-"`
}

func (n *NsHandler) OnAdd(obj interface{}) {
	n.NsMap.Add(obj.(*corev1.Namespace))
}

func (n *NsHandler) OnUpdate(oldObj, newObj interface{}) {
	n.NsMap.Update(newObj.(*corev1.Namespace))
}

func (n *NsHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*corev1.Namespace); ok {
		n.NsMap.Delete(d)
	}
}
