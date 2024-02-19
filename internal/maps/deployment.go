package maps

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"sort"
	"sync"
)

type DeploymentMap struct {
	data sync.Map // map: [namespaces : v1.Deployments]
}

func (d *DeploymentMap) Add(dep *v1.Deployment) {

	//判断当前命名空间下是否能够找到 deployments 数据
	if list, ok := d.data.Load(dep.Namespace); ok {
		list = append(list.([]*v1.Deployment), dep)
		d.data.Store(dep.Namespace, list)
	} else {
		//如果没有找到直接写入
		d.data.Store(dep.Namespace, []*v1.Deployment{dep})
	}
}

func (d *DeploymentMap) Update(dep *v1.Deployment) error {
	if list, ok := d.data.Load(dep.Namespace); ok {
		//遍历所有得数据，找到相同得进行替换
		for i, rangeDep := range list.([]*v1.Deployment) {
			if rangeDep.Name == dep.Name {
				list.([]*v1.Deployment)[i] = dep
				break
			}
		}
		return nil
	}
	return fmt.Errorf("deployment-%s not found", dep.Name)
}

func (d *DeploymentMap) Delete(dep *v1.Deployment) {
	if list, ok := d.data.Load(dep.Namespace); ok {
		for i, rangeDep := range list.([]*v1.Deployment) {
			if rangeDep.Name == dep.Name {
				newList := append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				d.data.Store(dep.Namespace, newList)
				break
			}
		}
	}
}

func (d *DeploymentMap) ListByNs(ns string) ([]*v1.Deployment, error) {
	if list, ok := d.data.Load(ns); ok {
		ret := list.([]*v1.Deployment)
		sort.Slice(ret, func(i, j int) bool {
			return ret[i].CreationTimestamp.Time.Before(ret[j].CreationTimestamp.Time)
		})
		return ret, nil
	}
	return nil, fmt.Errorf("record not found")
}

func (d *DeploymentMap) GetDeployment(ns string, depName string) (*v1.Deployment, error) {
	if list, ok := d.data.Load(ns); ok {
		for _, item := range list.([]*v1.Deployment) {
			if item.Name == depName {
				return item, nil
			}
		}
	}
	return nil, fmt.Errorf("record not found")
}
