package svc

import (
	"manager/internal/service"
)

type ServiceContext struct {
	CommonService     *service.CommonService     `inject:"-"`
	NameSpaceService  *service.NamespaceService  `inject:"-"`
	PodService        *service.PodService        `inject:"-"`
	DeploymentService *service.DeploymentService `inject:"-"`
	IngressService    *service.IngressService    `inject:"-"`
	SvcService        *service.SvcService        `inject:"-"`
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{}
}
