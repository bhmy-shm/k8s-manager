package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/gin-gonic/gin"
	"manager/internal/wscore"
	"manager/wire"
	"net/http"
)

type WsCtl struct {
	*wire.ServiceWire `inject:"-"`
}

func NewWsCtl() *WsCtl {
	return &WsCtl{}
}

func (w *WsCtl) Build(gofk *gofks.Gofk) {
	ws := gofk.Group("ws")
	ws.GET("/conn", w.Connect)
}

func (w *WsCtl) Name() string {
	return "WsCtl"
}

func (w *WsCtl) Connect(c *gin.Context) {
	conn, err := wscore.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		logx.Error(err)
		return
	} else {
		wscore.ClientMap.Store(conn)
		SuccessResp(c, "ok")
	}
}
