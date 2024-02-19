package controllers

import (
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"strings"
)

func RespField(k string, v interface{}) gin.H {
	return gin.H{
		k: v,
	}
}

func InternalResp(c *gin.Context, h gin.H) {
	c.JSON(http.StatusInsufficientStorage, h)
}

func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}

/*
IntoGvr 拆分资源 gvr

	@request = apps.v1.deployment  根据代码逻辑，将得到以下结果：
	@result：
		Version：“apps.v1”
		Resource：“deployment”
		Group：""

	@request = api.jtthink.com.v1beta1.tasks
	@result：
		Version: “api.jtthink.com.v1beta1”
		Resource: “tasks”
		Group: “api.jtthink.com”
*/
func IntoGvr(gvr string) schema.GroupVersionResource {
	list := strings.Split(gvr, ".")
	ret := schema.GroupVersionResource{}

	if len(list) == 2 {
		ret.Version, ret.Resource = list[0], list[1]
	} else if len(list) > 2 {
		lastIndex := len(list) - 1
		ret.Version, ret.Resource = list[lastIndex-1], list[lastIndex] // 最后一个
		ret.Group = strings.Join(list[0:lastIndex-1], ".")
	}
	return ret
}
