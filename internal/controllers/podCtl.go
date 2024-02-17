package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	"manager/internal/service"
	"manager/wire"
)

type PodCtl struct {
	*wire.ServiceWire `inject:"-"`
}

func NewPodCtl() *PodCtl {
	return &PodCtl{}
}

func (p *PodCtl) Build(gofk *gofks.Gofk) {
	pods := gofk.Group("pods")
	pods.GET("/list", p.GetAll)
}

func (p *PodCtl) Name() string {
	return "PodCtl"
}

func (p *PodCtl) GetAll(c *gin.Context) {

	ns := c.DefaultQuery("namespace", "default")
	page := service.StrToInt(c.DefaultQuery("page", "1"), 1)
	size := service.StrToInt(c.DefaultQuery("size", "5"), 5)

	//SuccessResp(c, p.Context().PodService.ListByNs(ns))             //全部列表
	SuccessResp(c, p.Context().PodService.PagePods(ns, page, size)) //分页列表
}
