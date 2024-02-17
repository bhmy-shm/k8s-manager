package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	"manager/wire"
)

type NamespaceCtl struct {
	*wire.ServiceWire `inject:"-"`
}

func NewNamespaceCtl() *NamespaceCtl {
	return &NamespaceCtl{}
}

func (n *NamespaceCtl) Build(gofk *gofks.Gofk) {
	ns := gofk.Group("namespace")
	ns.GET("/list", n.ListAll)
}

func (n *NamespaceCtl) Name() string {
	return "NsCtl"
}

func (n *NamespaceCtl) ListAll(c *gin.Context) {
	SuccessResp(c, n.Context().NameSpaceService.ListAll())
}
