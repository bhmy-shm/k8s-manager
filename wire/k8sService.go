package wire

import "manager/internal/service"

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

func (*ServiceConfig) NamespaceService() *service.NamespaceService { return service.NewNamespace() }

func (*ServiceConfig) IngressService() *service.IngressService { return service.Ingress() }

func (*ServiceConfig) SvcService() *service.SvcService { return service.Svc() }

func (*ServiceConfig) NodeService() *service.NodeService {
	return service.NewNodeService()
}

// storage 存储

func (*ServiceConfig) ResourceService() *service.ResourceService {
	return service.K8sResource()
}
