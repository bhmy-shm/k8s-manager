package service

import "manager/internal/maps"

type NamespaceService struct {
	NsMap *maps.NamespaceMap `inject:"-"`
}

func NewNamespace() *NamespaceService {
	return &NamespaceService{}
}

func (n *NamespaceService) ListAll() interface{} {
	return n.NsMap.ListAll()
}
