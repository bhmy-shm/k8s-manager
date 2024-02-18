package handler

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"log"
	"manager/internal/maps"
	"manager/internal/wscore"
)

type ServiceHandler struct {
	SvcMap *maps.ServiceMap `inject:"-"`
}

func (svc *ServiceHandler) OnAdd(obj interface{}) {
	svc.SvcMap.Add(obj.(*corev1.Service))
	ns := obj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "service",
			"result": gin.H{"ns": ns,
				"data": svc.SvcMap.ListAll(ns)},
		},
	)
}
func (svc *ServiceHandler) OnUpdate(oldObj, newObj interface{}) {
	err := svc.SvcMap.Update(newObj.(*corev1.Service))
	if err != nil {
		log.Println(err)
		return
	}
	ns := newObj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "service",
			"result": gin.H{"ns": ns,
				"data": svc.SvcMap.ListAll(ns)},
		},
	)
}
func (svc *ServiceHandler) OnDelete(obj interface{}) {
	svc.SvcMap.Delete(obj.(*corev1.Service))
	ns := obj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "service",
			"result": gin.H{"ns": ns,
				"data": svc.SvcMap.ListAll(ns)},
		},
	)
}
