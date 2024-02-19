package maps

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type MetaResMapper struct {
	Client *kubernetes.Clientset
	Mapper meta.RESTMapper
}

func NewMetaRespMapper() *MetaResMapper {
	mapper := &MetaResMapper{}
	mapper.RestMapper()
	return mapper
}

// RestMapper 所有api groupResource
// 用于初始化和返回一个 RESTMapper，用于将 API 资源映射到 REST 路径，以便在后续的操作中能够根据 API 资源信息进行 REST 请求。
// 通过获取 API 组资源信息并创建 RESTMapper，可以帮助在 Kubernetes 集群中进行资源操作时更方便地进行资源路径映射和管理
func (res *MetaResMapper) RestMapper() *MetaResMapper {

	res.Client, _ = kubernetes.NewForConfig(buildConfig())

	gr, err := restmapper.GetAPIGroupResources(res.Client.Discovery())
	if err != nil {
		errorx.Fatal(err)
	}

	res.Mapper = restmapper.NewDiscoveryRESTMapper(gr)
	return res
}

func buildConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "resources/config")
	if err != nil {
		logx.Errorf("build config is failed:%v", err)
		errorx.Fatal(err)
	}
	config.Insecure = false
	return config
}
