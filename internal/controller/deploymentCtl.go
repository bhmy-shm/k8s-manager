package controllers

import (
	"github.com/bhmy-shm/gofks/gofk"
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"view/internal/service"
)

type DeploymentCtl struct {
	Cli        *kubernetes.Clientset      `inject:"-"`
	DepService *service.DeploymentService `inject:"-"` //服务层
}

func NewDeploymentCtl() *DeploymentCtl {
	return &DeploymentCtl{}
}

func (this *DeploymentCtl) GetList(c *gin.Context) {
	ns, ok := c.GetQuery("namespace")
	if !ok {
		gofk.InternalResp(c, errorx.BadRequest("namespace", "get query is failed"))
		return
	}

	list, err := this.DepService.List(ns)
	if err != nil {
		gofk.InternalResp(c, errorx.InternalServer(err.Error(), this.Name()))
		return
	}
	gofk.Successful(c, list)
}

func (this *DeploymentCtl) Build(gofk *gofk.Gofk) {
	//路由
	deploy := gofk.Group("/deploy")
	deploy.GET("/list", this.GetList)
}

func (*DeploymentCtl) Name() string {
	return "deploymentCtl"
}
