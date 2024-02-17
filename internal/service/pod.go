package service

import (
	"github.com/gin-gonic/gin"
	"manager/internal/maps"
	"manager/model"
	"time"
)

type PodService struct {
	Common   *CommonService `inject:"-"`
	PodMap   *maps.PodMap   `inject:"-"`
	EventMap *maps.EventMap `inject:"-"`
}

func Pod() *PodService {
	return &PodService{}
}

func (p *PodService) PagePods(ns string, page, size int) *ItemsPage {
	pods := p.ListByNs(ns).([]*model.Pod)
	readyCount := 0 //就绪的pod数量
	allCount := 0   //总数量
	ipods := make([]interface{}, len(pods))
	for i, pod := range pods {
		allCount++
		ipods[i] = pod
		if pod.IsReady {
			readyCount++
		}
	}
	return p.Common.PageResource(page, size, ipods).
		SetExt(gin.H{"ReadyNum": readyCount, "AllNum": allCount})
}

func (p *PodService) ListByNs(ns string) interface{} {

	podList := p.PodMap.ListByNs(ns)
	ret := make([]*model.Pod, len(podList))

	if len(podList) > 0 {
		for i, pod := range podList {
			ret[i] = &model.Pod{
				Name:       pod.Name,
				NameSpace:  pod.Namespace,
				Images:     p.Common.getImagesByPod(pod.Spec.Containers),
				NodeName:   pod.Spec.NodeName,
				Phase:      string(pod.Status.Phase),    // 阶段
				IsReady:    p.Common.getPodIsReady(pod), //是否就绪
				IP:         []string{pod.Status.PodIP, pod.Status.HostIP},
				Message:    p.EventMap.GetMessage(pod.Namespace, "Pod", pod.Name),
				CreateTime: pod.CreationTimestamp.Format(time.DateTime),
			}
		}
	}

	return ret
}
