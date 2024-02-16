package handler

import (
	corev1 "k8s.io/api/core/v1"
	"manager/internal/maps"
)

type NsHandler struct {
	NsMap *maps.NsMapStruct `inject:"-"`
}

func (this *NsHandler) OnAdd(obj interface{}) {
	this.NsMap.Add(obj.(*corev1.Namespace))
}

func (this *NsHandler) OnUpdate(oldObj, newObj interface{}) {
	this.NsMap.Update(newObj.(*corev1.Namespace))

}

func (this *NsHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*corev1.Namespace); ok {
		this.NsMap.Delete(d)
	}
}
