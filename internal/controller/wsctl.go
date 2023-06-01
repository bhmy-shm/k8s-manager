package controllers

import (
	"github.com/bhmy-shm/gofks/gofk"
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"github.com/gin-gonic/gin"
	"view/internal/wscore"
)

type WsCtl struct{}

func NewWsCtl() *WsCtl {
	return &WsCtl{}
}

func (this *WsCtl) Connect(c *gin.Context) {
	client, err := wscore.Upgrader.Upgrade(c.Writer, c.Request, nil) //升级
	if err != nil {
		gofk.InternalResp(c, errorx.InternalServer(err.Error(), "wscore upgrader"))
		return
	} else {
		wscore.ClientMap.Store(client)
		gofk.Successful(c, "successful")
	}
}

func (this *WsCtl) Build(gofk *gofk.Gofk) {
	ws := gofk.Group("ws")
	ws.GET("/ws", this.Connect)
}

func (this *WsCtl) Name() string {
	return "WsCtl"
}
