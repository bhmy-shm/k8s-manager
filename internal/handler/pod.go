package handler

import (
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"manager/internal/maps"
	"manager/internal/service"
	"manager/internal/wscore"
)

type PodHandler struct {
	PodMap     *maps.PodMap        `inject:"-"`
	PodService *service.PodService `inject:"-"`
}

func (p *PodHandler) OnAdd(obj interface{}) {
	p.PodMap.Add(obj.(*corev1.Pod))
	ns := obj.(*corev1.Pod).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":   "pods",
			"result": gin.H{"ns": ns, "data": p.PodService.ListByNs(ns)}, //todo 分页显示
		},
	)
}

func (p *PodHandler) OnUpdate(oldObj, newObj interface{}) {
	err := p.PodMap.Update(newObj.(*corev1.Pod))
	if err != nil {
		logx.Error(err)
	} else {
		ns := newObj.(*corev1.Pod).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":   "pods",
				"result": gin.H{"ns": ns, "data": p.PodService.ListByNs(ns)},
			},
		)
	}
}

func (p *PodHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*corev1.Pod); ok {
		p.PodMap.Delete(d)
		ns := obj.(*corev1.Pod).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":   "pods",
				"result": gin.H{"ns": ns, "data": p.PodService.ListByNs(ns)},
			},
		)
	}
}
