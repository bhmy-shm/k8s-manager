package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"manager/internal/types"
	"manager/wire"
	"net/http"
	"sigs.k8s.io/yaml"
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

	//
	ingressGroup.Handle(http.MethodGet, "/detail", ig.DetailIngress)
	ingressGroup.Handle(http.MethodGet, "/yaml", ig.DetailIngressYaml)
	ingressGroup.Handle(http.MethodGet, "/list", ig.ListAll)
	ingressGroup.Handle(http.MethodGet, "/add", ig.AddIngress)
	ingressGroup.Handle(http.MethodPost, "/remove", ig.RemoveIngress)

	//
	ingressGroup.Handle(http.MethodPost, "/auth", ig.BasicAuthFileOrMap) //创建 basic-auth 所需要密文

}

func (ig *IngressCtl) AddIngress(c *gin.Context) {
	params := &types.IngressPost{}
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

func (ig *IngressCtl) DetailIngress(c *gin.Context) {
	ns := c.Param("namespace")
	name := c.Param("name")

	ingress := ig.Context().IngressService.Detail(ns, name) // 原生对象
	SuccessResp(c, ingress)
}

func (ig *IngressCtl) DetailIngressYaml(c *gin.Context) {
	ns := c.DefaultQuery("ns", "default")
	name := c.DefaultQuery("name", "")

	ingress, err := ig.NetworkingV1().Ingresses(ns).Get(c, name, v1.GetOptions{})
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	b, err := yaml.Marshal(ingress)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}
	SuccessResp(c, string(b))
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

// ------

func (ig *IngressCtl) BasicAuthFileOrMap(c *gin.Context) {

	auth := &types.GenAuth{}
	err := c.ShouldBindJSON(auth)
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	//t := c.DefaultQuery("t", "1") // 1是auth-file  2是 auth-map
	//为了简单，目前只支持单个用户名和密码生成 ，如果要扩展 自行修改secret内容
	secret := corev1.Secret{}
	secret.Name = auth.SName
	secret.Namespace = auth.Namespace

	if auth.Type == 1 {
		secret.Data = map[string][]byte{
			"auth": []byte(auth.UName + ":" + types.HashApr1(auth.UPwd)),
		}
	} else {
		secret.Data = map[string][]byte{
			auth.UName: []byte(types.HashApr1(auth.UPwd)),
		}
	}

	//调用k8s client 创建出 密文
	_, err = ig.CoreV1().Secrets(auth.Namespace).Create(c, &secret, v1.CreateOptions{})
	if err != nil {
		InternalResp(c, RespField("reason", err))
		return
	}

	SuccessResp(c, "ok")
}
