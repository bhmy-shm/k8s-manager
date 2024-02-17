package handler

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"manager/internal/maps"
)

type EventHandler struct {
	EventMap *maps.EventMap `inject:"-"`
}

func (e *EventHandler) storeData(obj interface{}, isDelete bool) {
	if event, ok := obj.(*corev1.Event); ok {
		key := fmt.Sprintf("%s_%s_%s", event.Namespace, event.InvolvedObject.Kind, event.InvolvedObject.Name)
		if !isDelete {
			e.EventMap.Data.Store(key, event)
		} else {
			e.EventMap.Data.Delete(key)
		}
	}
}
func (e *EventHandler) OnAdd(obj interface{}) {
	e.storeData(obj, false)
}

func (e *EventHandler) OnUpdate(oldObj, newObj interface{}) {
	e.storeData(newObj, false)
}

func (e *EventHandler) OnDelete(obj interface{}) {
	e.storeData(obj, true)
}
