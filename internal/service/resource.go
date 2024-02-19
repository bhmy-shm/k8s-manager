package service

import (
	"context"
	"github.com/bhmy-shm/gofks/core/gormx"
	"github.com/bhmy-shm/gofks/core/logx"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"manager/model"
)

type ResourceService struct {
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
