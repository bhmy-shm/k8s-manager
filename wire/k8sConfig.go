package wire

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"manager/internal/handler"
)

type K8sConfig struct {
	DepHandler     *handler.DepHandler     `inject:"-"`
	PodHandler     *handler.PodHandler     `inject:"-"`
	NsHandler      *handler.NsHandler      `inject:"-"`
	ServiceHandler *handler.ServiceHandler `inject:"-"`
	EventHandler   *handler.EventHandler   `inject:"-"`
	IngressHandler *handler.IngressHandler `inject:"-"`
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func k8sRestConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "resources/config")
	if err != nil {
		logx.Errorf("build config is failed:%v", err)
		errorx.Fatal(err)
	}
	config.Insecure = false
	return config
}

func (this *K8sConfig) InitClient() *kubernetes.Clientset {
	c, err := kubernetes.NewForConfig(k8sRestConfig())
	errorx.Fatal(err, "new k8s client for config")
	return c
}

func (this *K8sConfig) InitInformer() informers.SharedInformerFactory {
	fact := informers.NewSharedInformerFactory(this.InitClient(), 0)

	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(this.DepHandler)

	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(this.PodHandler)

	nsInformer := fact.Core().V1().Namespaces()
	nsInformer.Informer().AddEventHandler(this.NsHandler)

	serviceInformer := fact.Core().V1().Services() //监听service-svc
	serviceInformer.Informer().AddEventHandler(this.ServiceHandler)

	eventInformer := fact.Core().V1().Events() //监听event
	eventInformer.Informer().AddEventHandler(this.EventHandler)

	ingressInformer := fact.Networking().V1().Ingresses() //监听ingress
	ingressInformer.Informer().AddEventHandler(this.IngressHandler)

	fact.Start(wait.NeverStop)
	return fact
}
