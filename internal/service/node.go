package service

import (
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"manager/internal/maps"
	"manager/internal/types"
)

type NodeService struct {
	NodeMap *maps.NodeMap `inject:"-"`
	PodMap  *maps.PodMap  `inject:"-"`

	Common *CommonService       `inject:"-"`
	Metric *versioned.Clientset `inject:"-"`
}

func NewNodeService() *NodeService {
	return &NodeService{}
}

func (n *NodeService) ListAllNodes() []*types.NodeModel {
	list := n.NodeMap.ListAll()
	ret := make([]*types.NodeModel, len(list))
	for i, node := range list {
		nodeUsage := n.Common.GetNodeUsage(n.Metric, node)

		ret[i] = &types.NodeModel{
			Name:     node.Name,
			IP:       node.Status.Addresses[0].Address,
			HostName: node.Status.Addresses[1].Address,
			Labels:   n.Common.FilterLabels(node.Labels),
			Capacity: types.NewNodeCapacity(
				node.Status.Capacity.Cpu().Value(),
				node.Status.Capacity.Memory().Value(),
				node.Status.Capacity.Pods().Value(),
			),
			Usage: types.NewNodeUsage(
				n.PodMap.GetNum(node.Name),
				nodeUsage[0],
				nodeUsage[1],
			),
			CreateTime: node.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
	}
	return ret
}
