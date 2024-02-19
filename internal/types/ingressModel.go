package types

// IngressModel 列表返回参数
type (
	IngressOptions struct {
		IsCros    bool //是否开启跨域
		IsRewrite bool //是否开启路径重写
	}
	IngressModel struct {
		Name       string
		NameSpace  string
		CreateTime string
		Host       string
		Options    IngressOptions
	}
)

// IngressPath 配置
type IngressPath struct {
	Path    string `json:"path"`
	SvcName string `json:"svcName"`
	Port    string `json:"port"`
}

// IngressRules  规则
type IngressRules struct {
	Host  string         `json:"host"`
	Paths []*IngressPath `json:"paths"`
}

// IngressPost 提交Ingress 对象
type IngressPost struct {
	Name        string
	Namespace   string
	Rules       []*IngressRules //路由规则
	Annotations string          //标签
}

/*
	路由规则
	rules:
	  - host: www.shm.bhmy
	    http:
	      paths:
	        - path: /
	          pathType: Prefix
	          backend:
	            service:
	              name: testrouter
	              port:
	                number: 80
*/

/*
	标签
	metadata:
	  name: testingress
	  annotations:
	    nginx.ingress.kubernetes.io/server-snippet:
*/
