package handler

import (
	"github.com/gin-gonic/gin"
	"k8s.io/api/networking/v1"
	"log"
	"manager/internal/maps"
	"manager/internal/wscore"
)

type IngressHandler struct {
	IngressMap *maps.IngressMap `inject:"-"`
}

func (ig *IngressHandler) OnAdd(obj interface{}) {
	ig.IngressMap.Add(obj.(*v1.Ingress))
	ns := obj.(*v1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "ingress",
			"result": gin.H{"ns": ns,
				"data": ig.IngressMap.ListByNs(ns)},
		},
	)
}

func (ig *IngressHandler) OnUpdate(oldObj, newObj interface{}) {
	log.Println("ingress update ")
	err := ig.IngressMap.Update(newObj.(*v1.Ingress))
	if err != nil {
		log.Println(err)
		return
	}
	ns := newObj.(*v1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "ingress",
			"result": gin.H{"ns": ns,
				"data": ig.IngressMap.ListByNs(ns)},
		},
	)
}

func (ig *IngressHandler) OnDelete(obj interface{}) {

	log.Println("ingress delete ")

	ig.IngressMap.Delete(obj.(*v1.Ingress))
	ns := obj.(*v1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type": "ingress",
			"result": gin.H{"ns": ns,
				"data": ig.IngressMap.ListByNs(ns)},
		},
	)
}
