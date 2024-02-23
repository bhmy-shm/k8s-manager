package service

import (
	"github.com/gin-gonic/gin"
	"manager/internal/maps"
	"manager/internal/types"
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
	pods := p.ListByNs(ns).([]*types.Pod)
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
	ret := make([]*types.Pod, len(podList))

	if len(podList) > 0 {
		for i, pod := range podList {
			ret[i] = &types.Pod{
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

func (p *PodService) GetPodContainer(ns, podName string) []*types.ContainerModel {
	ret := make([]*types.ContainerModel, 0)
	pod := p.PodMap.Get(ns, podName)
	if pod != nil {
		for _, c := range pod.Spec.Containers {
			ret = append(ret, &types.ContainerModel{
				Name: c.Name,
			})
		}
	}
	return ret
}
