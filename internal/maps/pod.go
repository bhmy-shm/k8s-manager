package maps

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sort"
	"sync"
)

type PodMap struct {
	data sync.Map // namespaces : corev1.Pod
}

func (P *PodMap) ListByNs(ns string) []*corev1.Pod {
	if list, ok := P.data.Load(ns); ok {
		ret := list.([]*corev1.Pod)
		sort.Slice(ret, func(i, j int) bool {
			return ret[i].CreationTimestamp.Time.Before(ret[j].CreationTimestamp.Time)
		})
		return ret
	}
	return nil
}

func (P *PodMap) Get(ns string, podName string) *corev1.Pod {
	if list, ok := P.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			if pod.Name == podName {
				return pod
			}
		}
	}
	return nil
}

func (P *PodMap) Add(pod *corev1.Pod) {
	if list, ok := P.data.Load(pod.Namespace); ok {
		list = append(list.([]*corev1.Pod), pod)
		P.data.Store(pod.Namespace, list)
	} else {
		P.data.Store(pod.Namespace, []*corev1.Pod{pod})
	}
}

func (P *PodMap) Update(pod *corev1.Pod) error {
	if list, ok := P.data.Load(pod.Namespace); ok {
		for i, range_pod := range list.([]*corev1.Pod) {
			if range_pod.Name == pod.Name {
				list.([]*corev1.Pod)[i] = pod
			}
		}
		return nil
	}
	return fmt.Errorf("Pod-%s not found", pod.Name)
}

func (P *PodMap) Delete(pod *corev1.Pod) {
	if list, ok := P.data.Load(pod.Namespace); ok {
		for i, range_pod := range list.([]*corev1.Pod) {
			if range_pod.Name == pod.Name {
				newList := append(list.([]*corev1.Pod)[:i], list.([]*corev1.Pod)[i+1:]...)
				P.data.Store(pod.Namespace, newList)
				break
			}
		}
	}
}

// ListByLabels 根据标签获取 POD列表
func (P *PodMap) ListByLabels(ns string, labels []map[string]string) ([]*corev1.Pod, error) {
	ret := make([]*corev1.Pod, 0)
	if list, ok := P.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			for _, label := range labels {
				if reflect.DeepEqual(pod.Labels, label) { //标签完全匹配
					ret = append(ret, pod)
				}
			}
		}
		return ret, nil
	}
	return nil, fmt.Errorf("pods not found ")
}

// 根据节点名称 获取pods数量
func (P *PodMap) GetNum(nodeName string) (num int) {
	P.data.Range(func(key, value interface{}) bool {
		list := value.([]*corev1.Pod)
		for _, pod := range list {
			if pod.Spec.NodeName == nodeName {
				num++
			}
		}
		return true
	})
	return
}
