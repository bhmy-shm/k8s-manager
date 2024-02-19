package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	"manager/model"
	"manager/wire"
)

type ResourceController struct {
	*wire.ServiceWire `inject:"-"`
}

func NewResourceCtl() *ResourceController {
	return &ResourceController{}
}

func (r *ResourceController) Name() string {
	return "ResourceController"
}

func (r *ResourceController) Build(gofk *gofks.Gofk) {

	resourceGroup := gofk.Group("resource")
	resourceGroup.Handle("GET", "/list", r.List)
}

func (r *ResourceController) List(ctx *gin.Context) {
	gvrQuery := ctx.Query("gvr")
	ns := ctx.DefaultQuery("namespace", "default")

	gvr := IntoGvr(gvrQuery)
	if gvr.Empty() {
		InternalResp(ctx, RespField("reason", "gvr query is empty"))
		return
	}

	resList, err := r.Context().ResourceService.ListByNs(
		r.Context().ResourceModel,
		model.WithNamespace(ns),
		model.WithGroup(gvr.Group),
		model.WithVersion(gvr.Version),
		model.WithResource(gvr.Resource),
	)
	if err != nil {
		InternalResp(ctx, RespField("reason", ""))
		return
	}

	SuccessResp(ctx, resList)
}
