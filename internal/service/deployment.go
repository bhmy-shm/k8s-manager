package service

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"manager/internal/maps"
	"manager/internal/types"
)

type DeploymentService struct {
	Common *CommonService      `inject:"-"`
	DepMap *maps.DeploymentMap `inject:"-"` //数据层
}

func Deployment() *DeploymentService {
	return &DeploymentService{}
}

func (d *DeploymentService) List(ns string) ([]*types.Deployment, error) {
	depList, err := d.DepMap.ListByNs(ns)
	if err != nil {
		logx.Error("deployment list By Namespace failed:", err)
		return nil, errorx.Wrap(err, "deployment list By Namespace failed")
	}

	var res []*types.Deployment
	for _, item := range depList {
		res = append(res, &types.Deployment{
			NameSpace:  item.Namespace,
			Name:       item.Name,
			Replicas:   [3]int32{item.Status.Replicas, item.Status.AvailableReplicas, item.Status.UnavailableReplicas},
			Images:     d.Common.getImages(item),
			IsComplete: d.Common.getIsComplete(item),
			Message:    d.Common.getMessage(item),
			CreateTime: d.Common.getCreateTime(item),
		})
	}
	return res, nil
}
