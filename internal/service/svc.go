package service

import (
	"manager/internal/maps"
	"manager/model"
)

type SvcService struct {
	Common *CommonService   `inject:"-"`
	SvcMap *maps.ServiceMap `inject:"-"` //数据层
}

func Svc() *SvcService {
	return &SvcService{}
}
func (s *SvcService) ListByNs(ns string) ([]*model.ServiceModel, error) {

	depList := s.SvcMap.ListAll(ns)
	ret := make([]*model.ServiceModel, len(depList))

	for i, item := range depList {

		ret[i] = &model.ServiceModel{
			Name:       item.Name,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace:  item.Namespace,
		}
	}
	return ret, nil
}
