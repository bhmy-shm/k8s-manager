package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	"manager/wire"
	"net/http"
)

type SvcCtl struct {
	*wire.ServiceWire `inject:"-"`
}

func NewSvcCtl() *SvcCtl {
	return &SvcCtl{}
}

func (svc *SvcCtl) Build(gofk *gofks.Gofk) {
	svcGroup := gofk.Group("svc")
	svcGroup.Handle(http.MethodGet, "/list", svc.ListAll)
}

func (*SvcCtl) Name() string {
	return "SvcCtl"
}

func (svc *SvcCtl) ListAll(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")

	respList, err := svc.ServiceWire.Context().SvcService.ListByNs(ns)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	SuccessResp(c, respList)
}
