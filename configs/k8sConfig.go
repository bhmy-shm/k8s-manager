package configs

import (
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"view/internal/handler"
)

type K8sConfig struct {
	DepHandler *handler.DepHandler `inject:"-"`
	PodHandler *handler.PodHandler `inject:"-"`
	NsHandler  *handler.NsHandler  `inject:"-"`
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func k8sRestConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "resources/config")
	errorx.Fatal(err, "build config is failed")
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

	fact.Start(wait.NeverStop)
	return fact
}
