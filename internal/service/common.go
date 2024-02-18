package service

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
	"time"
)

type CommonService struct{}

func NewCommonService() *CommonService {
	return &CommonService{}
}

func (c *CommonService) getImages(dep *v1.Deployment) string {
	return c.getImagesByPod(dep.Spec.Template.Spec.Containers)
}

func (c *CommonService) getImagesByPod(containers []corev1.Container) string {
	images := containers[0].Image
	if imgLen := len(containers); imgLen > 1 {
		images += fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}
	return images
}

func (c *CommonService) getIsComplete(dep *v1.Deployment) bool {
	return dep.Status.Replicas == dep.Status.AvailableReplicas
}

func (c *CommonService) getMessage(dep *v1.Deployment) string {
	for _, condition := range dep.Status.Conditions {
		if condition.Type == "Available" && condition.Status != "True" {
			return condition.Message
		}
	}
	return ""
}

func (c *CommonService) getCreateTime(dep *v1.Deployment) string {
	return dep.CreationTimestamp.Format(time.DateTime)
}

// ---------- pod ---------

/*
在Kubernetes中，Spec.ReadinessGates是一个Pod规范（PodSpec）中的字段，用于定义Pod的就绪条件。
每个元素都是一个PodReadinessGate对象，该对象定义了一个条件和一个期望的状态。

Pod的就绪条件用于指示Pod是否已准备好接收流量。当Pod的所有就绪条件都满足时，Kubernetes认为该Pod已准备好，并将其标记为就绪状态。
这是在服务发现、负载均衡和滚动升级等场景中非常有用的。

PodReadinessGate对象包含两个字段：
1. ConditionType：指定了一个条件类型，用于与Pod的状态条件进行匹配。常见的条件类型包括PodScheduled、ContainersReady和Initialized等。
2. ConditionStatus：指定了期望的条件状态。可以是True、False或Unknown。
*/
func (c *CommonService) getPodIsReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != "Running" {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Status != "True" {
			return false
		}
	}
	for _, rg := range pod.Spec.ReadinessGates {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == rg.ConditionType && condition.Status != "True" {
				return false
			}
		}
	}
	return false
}

type ItemsPage struct {
	Total   int           //一共多少条
	Current int           //当前页
	Size    int           // 页尺寸
	PageNum int           //一共多少页
	Data    []interface{} //数据
	Ext     interface{}   //扩展信息，方便插入值 给前端用
}

func (this *ItemsPage) SetExt(ext interface{}) *ItemsPage {
	this.Ext = ext
	return this
}

// --------------------- 分页 -----------------------

func StrToInt(str string, def int) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return ret
}

func (c *CommonService) PageResource(current, size int, list []interface{}) *ItemsPage {
	total := len(list)
	if size == 0 || size > total {
		size = 5 //默认 每页5个
	}
	if current <= 0 {
		current = 1
	}
	pageInfo := &ItemsPage{Total: total, Size: size}
	//计算总页数
	pageNum := 1
	if pageInfo.Total > size {
		pageNum = pageInfo.Total / size
		if pageInfo.Total%size != 0 {
			pageNum++
		}
	}
	if current > pageNum {
		current = 1
	}
	pageInfo.Current = current       //重新赋值Current ----当前页
	newSet := make([]interface{}, 0) //构建一个新的 切片

	if current*size > pageInfo.Total {
		newSet = append(newSet, list[(current-1)*size:]...)
	} else {
		//fmt.Println((current-1)*size,":",size)
		//  1 ,2,3,4,5,6
		// [) 左闭右开
		// 0,2   list[0:2]
		newSet = append(newSet, list[(current-1)*size:(current-1)*size+size]...)
	}
	//重新整理赋值
	pageInfo.Data = newSet
	pageInfo.PageNum = pageNum
	return pageInfo
}

// ------------ ingress 解析标签 ----------------

// 解析标签
func (c *CommonService) parseAnnotations(str string) map[string]string {
	replace := []string{"\t", " ", "\n", "\r\n"}
	for _, r := range replace {
		str = strings.ReplaceAll(str, r, "")
	}

	ret := make(map[string]string)
	list := strings.Split(str, ";")

	for _, item := range list {
		ss := strings.Split(item, ":")
		if len(str) == 2 {
			ret[ss[0]] = ss[1]
		}
	}
	return ret

}
