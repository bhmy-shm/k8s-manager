package wire

import (
	"fmt"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"manager/internal/handler"
	"os"
)

type K8sConfig struct {
	DepHandler      *handler.DepHandler      `inject:"-"`
	PodHandler      *handler.PodHandler      `inject:"-"`
	NsHandler       *handler.NsHandler       `inject:"-"`
	ServiceHandler  *handler.ServiceHandler  `inject:"-"`
	EventHandler    *handler.EventHandler    `inject:"-"`
	IngressHandler  *handler.IngressHandler  `inject:"-"`
	ResourceHandler *handler.ResourceHandler `inject:"-"`
	NodeHandler     *handler.NodeHandler     `inject:"-"`
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func (conf *K8sConfig) k8sRestConfigDefault() *rest.Config {
	// 取用户目录   Linux： ~   /home/xxx
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultConfigPath := fmt.Sprintf("%s/.kube/config", home)

	config, err := clientcmd.BuildConfigFromFlags("", defaultConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func K8sRestConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "resources/config")
	if err != nil {
		logx.Errorf("build config is failed:%v", err)
		errorx.Fatal(err)
	}
	config.Insecure = false
	return config
}

func (conf *K8sConfig) InitClient() *kubernetes.Clientset {
	c, err := kubernetes.NewForConfig(K8sRestConfig())
	errorx.Fatal(err, "new k8s client for config")
	return c
}

func (conf *K8sConfig) InitDynamicClient() dynamic.Interface {
	client, err := dynamic.NewForConfig(K8sRestConfig())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

//  ---- 内置k8s 资源 ----

func (conf *K8sConfig) InitInformer() informers.SharedInformerFactory {
	fact := informers.NewSharedInformerFactory(conf.InitClient(), 0)

	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(conf.DepHandler)

	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(conf.PodHandler)

	nsInformer := fact.Core().V1().Namespaces()
	nsInformer.Informer().AddEventHandler(conf.NsHandler)

	serviceInformer := fact.Core().V1().Services() //监听service-svc
	serviceInformer.Informer().AddEventHandler(conf.ServiceHandler)

	eventInformer := fact.Core().V1().Events() //监听event
	eventInformer.Informer().AddEventHandler(conf.EventHandler)

	ingressInformer := fact.Networking().V1().Ingresses() //监听ingress
	ingressInformer.Informer().AddEventHandler(conf.IngressHandler)

	NodeInformer := fact.Core().V1().Nodes() //监听node节点
	NodeInformer.Informer().AddEventHandler(conf.NodeHandler)

	fact.Start(wait.NeverStop)
	return fact
}

//  ---- 第三方k8s 资源 ----

func (conf *K8sConfig) InitWatchFactory() dynamicinformer.DynamicSharedInformerFactory {
	dynamicClient := conf.InitDynamicClient()

	/*
		defaultResync: 该参数传0
		表示不启用重新同步（resync）功能，SharedIndexInformer 将不会定期从 Kubernetes API 服务器中拉取对象的最新状态，
		仅在对象发生变化时才会获取最新数据
	*/
	watch := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 0)

	//监听pods变化记录数据库
	gvrPods := watch.ForResource(schema.GroupVersionResource{Version: "v1", Resource: "pods"})
	gvrPods.Informer().AddEventHandler(conf.ResourceHandler)

	//监听configMap变化记录数据库
	configMaps := watch.ForResource(schema.GroupVersionResource{Version: "v1", Resource: "configmaps"})
	configMaps.Informer().AddEventHandler(conf.ResourceHandler)

	watch.Start(wait.NeverStop)
	return watch
}
