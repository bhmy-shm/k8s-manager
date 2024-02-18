package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	"manager/model"
	"manager/wire"
	"net/http"
)

type IngressCtl struct {
	*wire.ServiceWire `inject:"-"`
}

func NewIngressCtl() *IngressCtl {
	return &IngressCtl{}
}

func (*IngressCtl) Name() string {
	return "IngressCtl"
}

func (ig *IngressCtl) Build(gofk *gofks.Gofk) {
	ingressGroup := gofk.Group("ingress")
	ingressGroup.Handle(http.MethodGet, "/list", ig.ListAll)
	ingressGroup.Handle(http.MethodGet, "/add", ig.AddIngress)
	ingressGroup.Handle(http.MethodPost, "/remove", ig.RemoveIngress)
}

func (ig *IngressCtl) AddIngress(c *gin.Context) {
	params := &model.IngressPost{}
	err := c.ShouldBindJSON(params)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	err = ig.ServiceWire.Context().IngressService.AddPostIngress(params)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	SuccessResp(c, "ok")
}

func (ig *IngressCtl) RemoveIngress(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	name := c.DefaultQuery("name", "")

	err := ig.ServiceWire.Context().IngressService.RemoveIngress(ns, name)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	SuccessResp(c, "ok")
}

func (ig *IngressCtl) ListAll(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")

	list, err := ig.ServiceWire.Context().IngressService.ListByNs(ns)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	SuccessResp(c, list)
}
