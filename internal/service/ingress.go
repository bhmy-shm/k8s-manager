package service

import (
	"context"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"manager/internal/maps"
	"manager/internal/types"
	"strconv"
)

type IngressService struct {
	Common     *CommonService        `inject:"-"`
	IngressMap *maps.IngressMap      `inject:"-"` //数据层
	Client     *kubernetes.Clientset `inject:"-"`
}

func Ingress() *IngressService {
	return &IngressService{}
}

func (ig *IngressService) getIngressOptions(t IngressCode, item *v1.Ingress) bool {
	if _, ok := item.Annotations[findTage(t)]; ok {
		return true
	}
	return false
}

func (ig *IngressService) ListByNs(ns string) ([]*types.IngressModel, error) {

	depList := ig.IngressMap.ListByNs(ns)
	ret := make([]*types.IngressModel, len(depList))

	for i, item := range depList {

		ret[i] = &types.IngressModel{
			Name:       item.Name,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace:  item.Namespace,
			Host:       item.Spec.Rules[0].Host,
			Options: types.IngressOptions{
				IsCros:    ig.getIngressOptions(IngressCodeCROS, item),
				IsRewrite: ig.getIngressOptions(IngressCodeRewrite, item),
			},
		}
	}
	return ret, nil
}

func (ig *IngressService) AddPostIngress(post *types.IngressPost) error {

	var (
		className    = "nginx"
		ingressRules = []v1.IngressRule{}
	)

	//拼凑rule规则
	for _, r := range post.Rules {
		httpRuleValue := &v1.HTTPIngressRuleValue{}
		rulePaths := make([]v1.HTTPIngressPath, 0)
		for _, pathCfg := range r.Paths {
			port, err := strconv.Atoi(pathCfg.Port)
			if err != nil {
				return err
			}
			rulePaths = append(rulePaths, v1.HTTPIngressPath{
				Path: pathCfg.Path,
				Backend: v1.IngressBackend{
					Service: &v1.IngressServiceBackend{
						Name: pathCfg.SvcName,
						Port: v1.ServiceBackendPort{
							Number: int32(port),
						},
					},
				},
			})
		}
		httpRuleValue.Paths = rulePaths
		rule := v1.IngressRule{
			Host: r.Host,
			IngressRuleValue: v1.IngressRuleValue{
				HTTP: httpRuleValue,
			},
		}
		ingressRules = append(ingressRules, rule)
	}

	// 凑 Ingress对象
	ingress := &v1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        post.Name,
			Namespace:   post.Namespace,
			Annotations: ig.Common.parseAnnotations(post.Annotations), //todo 需要改进，不要做复杂的字符串解析
		},
		Spec: v1.IngressSpec{
			IngressClassName: &className,
			Rules:            ingressRules,
		},
	}

	//调用创建ingress接口
	_, err := ig.Client.NetworkingV1().
		Ingresses(post.Namespace).
		Create(context.Background(), ingress, metav1.CreateOptions{})

	return err
}

func (ig *IngressService) RemoveIngress(namespace, ingress string) error {
	return ig.Client.NetworkingV1().
		Ingresses(namespace).
		Delete(context.Background(), ingress, metav1.DeleteOptions{})
}
