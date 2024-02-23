package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	"manager/wire"
)

type NodeCtl struct {
	*wire.ServiceWire `inject:"-"`
}

func NewNodeCtl() *NodeCtl {
	return &NodeCtl{}
}

func (n *NodeCtl) Build(gofk *gofks.Gofk) {
	nodeGroup := gofk.Group("node")
	nodeGroup.Handle("GET", "/list", n.ListAll)
}

func (n *NodeCtl) Name() string {
	return "NodeCtl"
}

func (n *NodeCtl) ListAll(c *gin.Context) {
	list := n.Context().NodeService.ListAllNodes()
	SuccessResp(c, list)
}
