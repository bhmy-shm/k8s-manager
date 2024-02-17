package svc

import (
	"manager/internal/service"
)

type ServiceContext struct {
	NameSpaceService  *service.NamespaceService  `inject:"-"`
	PodService        *service.PodService        `inject:"-"`
	DeploymentService *service.DeploymentService `inject:"-"`
	CommonService     *service.CommonService     `inject:"-"`
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{}
}
