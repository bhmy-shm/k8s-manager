package controllers

import (
	"github.com/bhmy-shm/gofks/gofk"
	"github.com/gin-gonic/gin"
	"view/internal/maps"
)

type NsCtl struct {
	NsMap *maps.NsMapStruct `inject:"-"`
}

func NewNsCtl() *NsCtl {
	return &NsCtl{}
}

func (this *NsCtl) ListAll(c *gin.Context) {
	gofk.Successful(c, this.NsMap.ListAll())
}

func (this *NsCtl) Build(gofk *gofk.Gofk) {
	ns := gofk.Group("ns")
	ns.GET("/nsList", this.ListAll)
}

func (*NsCtl) Name() string {
	return "NsCtl"
}
