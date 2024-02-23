package handler

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"manager/internal/maps"
	"manager/internal/service"
	"manager/internal/wscore"
)

type NodeHandler struct {
	NodeMap     *maps.NodeMap        `inject:"-"`
	NodeService *service.NodeService `inject:"-"`
}

func (this *NodeHandler) OnAdd(obj interface{}) {
	this.NodeMap.Add(obj.(*corev1.Node))

	wscore.ClientMap.SendAll(
		gin.H{
			"type": "node",
			"result": gin.H{"ns": "node",
				"data": this.NodeService.ListAllNodes()},
		},
	)
}
func (this *NodeHandler) OnUpdate(oldObj, newObj interface{}) {
	//重点： 只要update返回true 才会发送 。否则不发送
	if this.NodeMap.Update(newObj.(*corev1.Node)) {
		wscore.ClientMap.SendAll(
			gin.H{
				"type": "node",
				"result": gin.H{"ns": "node",
					"data": this.NodeService.ListAllNodes()},
			},
		)
	}
}
func (this *NodeHandler) OnDelete(obj interface{}) {
	this.NodeMap.Delete(obj.(*corev1.Node))

	wscore.ClientMap.SendAll(
		gin.H{
			"type": "node",
			"result": gin.H{"ns": "node",
				"data": this.NodeService.ListAllNodes()},
		},
	)
}
