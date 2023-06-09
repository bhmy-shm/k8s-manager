package configs

import "view/internal/service"

type ServiceConfig struct{}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}

func (*ServiceConfig) CommonService() *service.CommonService {
	return service.NewCommonService()
}

func (*ServiceConfig) DeploymentService() *service.DeploymentService {
	return service.Deployment()
}

func (*ServiceConfig) PodService() *service.PodService {
	return service.Pod()
}
