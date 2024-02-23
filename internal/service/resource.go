package service

import (
	"context"
	"github.com/bhmy-shm/gofks/core/gormx"
	"github.com/bhmy-shm/gofks/core/logx"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"manager/model"
)

type ResourceService struct {
	*kubernetes.Clientset `inject:"-"`
	Conf                  *rest.Config `inject:"-"`
}

func K8sResource() *ResourceService {
	return &ResourceService{}
}

func (r *ResourceService) ListByNs(storage model.K8sResourceModel, opts ...gormx.SqlOptions) ([]runtime.Object, error) {

	log.Println("resourceService svc isExist")

	resList := []*model.K8sResource{}
	err := storage.Query(context.Background(), resList, opts...)
	if err != nil {
		return make([]runtime.Object, 0), err
	}

	//解析json
	objList := make([]runtime.Object, len(resList))
	for i, res := range resList {
		obj := &unstructured.Unstructured{}
		if err = obj.UnmarshalJSON([]byte(res.Object)); err != nil {
			logx.Error(err)
			continue
		}
		objList[i] = obj
	}
	return objList, nil
}

func (r *ResourceService) RemoteExecutor(namespace, name, resource, container string, command []string) remotecommand.Executor {

	option := &v1.PodExecOptions{
		Container: container,
		Command:   command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}
	req := r.CoreV1().RESTClient().Post().
		Resource(resource). //resource = pods
		Namespace(namespace).
		Name(name). //name nginx-7875f55f56-5qqbr"
		SubResource("exec").
		Param("color", "false").
		VersionedParams(
			option,
			scheme.ParameterCodec,
		)

	exec, err := remotecommand.NewSPDYExecutor(r.Conf, "POST", req.URL())
	if err != nil {
		panic(err)
	}
	return exec
}
