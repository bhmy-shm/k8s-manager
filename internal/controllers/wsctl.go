package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/tools/remotecommand"
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
	ws.GET("/podConn", w.PodConnect)
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

func (w *WsCtl) PodConnect(c *gin.Context) {
	wsClient, err := wscore.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	//解析ws url 的 query
	ns := c.Query("namespace")
	if len(ns) <= 0 {
		InternalResp(c, RespField("reason", "namespace param err"))
		return
	}

	name := c.Query("name")
	if len(name) <= 0 {
		InternalResp(c, RespField("reason", "name param err"))
		return
	}

	resource := c.Query("resource")
	if len(resource) <= 0 {
		InternalResp(c, RespField("reason", "resource param err"))
		return
	}

	container := c.Query("container")
	if len(container) <= 0 {
		InternalResp(c, RespField("reason", "container param err"))
		return
	}

	shellClient := wscore.NewWsShellClient(wsClient)
	err = w.Context().ResourceService.RemoteExecutor(ns, name, resource, container, []string{"sh"}).
		Stream(remotecommand.StreamOptions{
			Stdin:  shellClient,
			Stdout: shellClient,
			Stderr: shellClient,
			Tty:    true,
		})
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	SuccessResp(c, "ok")
}
