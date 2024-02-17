package handler

import (
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	"manager/internal/maps"
	"manager/internal/service"
	"manager/internal/wscore"
	"manager/model"
)

type DepHandler struct {
	DepMap     *maps.DeploymentMap        `inject:"-"`
	DepService *service.DeploymentService `inject:"-"`
}

func (d *DepHandler) OnAdd(obj interface{}) {

	var (
		err  error
		list []*model.Deployment
		ns   = obj.(*v1.Deployment).Namespace
	)

	deploy := obj.(*v1.Deployment)
	d.DepMap.Add(deploy)
	list, err = d.DepService.List(ns)
	if err == nil {
		wscore.ClientMap.SendAll(
			gin.H{
				"type":   "deployments",
				"result": gin.H{"ns": ns, "data": list},
			},
		)
	}
}

func (d *DepHandler) OnUpdate(oldObj, newObj interface{}) {
	var (
		err  error
		list []*model.Deployment
		ns   = newObj.(*v1.Deployment).Namespace
	)

	err = d.DepMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		logx.Error(err)
	} else {
		list, err = d.DepService.List(ns)
		if err == nil {
			wscore.ClientMap.SendAll(
				gin.H{
					"type":   "deployments",
					"result": gin.H{"ns": ns, "data": list},
				},
			)
		}
	}
}

func (d *DepHandler) OnDelete(obj interface{}) {
	if deploy, ok := obj.(*v1.Deployment); ok {

		var (
			err  error
			list []*model.Deployment
			ns   = obj.(*v1.Deployment).Namespace
		)

		d.DepMap.Delete(deploy)
		list, err = d.DepService.List(ns)
		if err == nil {
			wscore.ClientMap.SendAll(
				gin.H{
					"type":   "deployments",
					"result": gin.H{"ns": ns, "data": list},
				},
			)
		}
	}
}
