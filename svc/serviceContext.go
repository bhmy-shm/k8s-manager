package svc

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"manager/internal/service"
	"manager/model"
)

type ServiceContext struct {
	serviceList
	serviceModel
}

type serviceList struct {
	CommonService     *service.CommonService     `inject:"-"`
	NameSpaceService  *service.NamespaceService  `inject:"-"`
	PodService        *service.PodService        `inject:"-"`
	DeploymentService *service.DeploymentService `inject:"-"`
	IngressService    *service.IngressService    `inject:"-"`
	SvcService        *service.SvcService        `inject:"-"`
	ResourceService   *service.ResourceService   `inject:"-"`
	NodeService       *service.NodeService       `inject:"-"`
}

type serviceModel struct {
	ResourceModel model.K8sResourceModel
}

func NewServiceContext(c *gofkConf.Config) *ServiceContext {
	return &ServiceContext{
		serviceModel: serviceModel{
			ResourceModel: model.NewK8sResourceModel(c),
		},
	}
}
