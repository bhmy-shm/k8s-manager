package controllers

import (
	"context"
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"manager/internal/service"
	"manager/wire"
	"net/http"
	"time"
)

type PodCtl struct {
	*wire.ServiceWire     `inject:"-"`
	*kubernetes.Clientset `inject:"-"`
}

func NewPodCtl() *PodCtl {
	return &PodCtl{}
}

func (p *PodCtl) Build(gofk *gofks.Gofk) {
	pods := gofk.Group("pods")
	pods.GET("/list", p.getAll)
	pods.GET("/containers", p.containers) //列出pods名称列表
	pods.GET("/followLog", p.followLog)   //阻塞刷新日志内容
}

func (p *PodCtl) Name() string {
	return "PodCtl"
}

func (p *PodCtl) getAll(c *gin.Context) {

	ns := c.DefaultQuery("namespace", "default")
	page := service.StrToInt(c.DefaultQuery("page", "1"), 1)
	size := service.StrToInt(c.DefaultQuery("size", "5"), 5)

	//SuccessResp(c, p.Context().PodService.ListByNs(ns))             //全部列表
	SuccessResp(c, p.Context().PodService.PagePods(ns, page, size)) //分页列表
}

func (p *PodCtl) containers(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	podName := c.DefaultQuery("name", "")

	list := p.Context().PodService.GetPodContainer(ns, podName)
	SuccessResp(c, list)
}

func (p *PodCtl) followLog(c *gin.Context) {

	ns := c.DefaultQuery("ns", "default")
	podName := c.DefaultQuery("podname", "")
	cName := c.DefaultQuery("cname", "")

	var tailLine int64 = 100
	opt := &v1.PodLogOptions{
		Follow:    true,
		Container: cName,
		TailLines: &tailLine, //从日志底部开始显示100行
	}

	//设置日志超时时间
	cc, cancel := context.WithTimeout(c, time.Minute*10)
	defer cancel()

	req := p.CoreV1().Pods(ns).GetLogs(podName, opt) //follow 持续监听
	reader, err := req.Stream(cc)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}
	defer reader.Close()

	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			break
		}

		w, err := c.Writer.Write(buf[0:n])
		if w == 0 || err != nil {
			break
		}
		c.Writer.(http.Flusher).Flush()
	}

}
