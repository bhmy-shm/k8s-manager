package types

import "github.com/GehirnInc/crypt/apr1_crypt"

// IngressModel 列表返回参数
type (
	IngressOptions struct {
		IsCros      bool //是否开启跨域
		IsRewrite   bool //是否开启路径重写
		IsAuth      bool
		IsRateLimit bool
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
	IsUpdate    bool   //是否更新ingress
	Name        string //ingress名称
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

type GenAuth struct {
	Type      int    `json:"type"`
	Namespace string `json:"namespace"`
	SName     string `json:"sname"` // secretname
	UName     string `json:"uname"` // 用户名
	UPwd      string `json:"upwd"`  //用户密码
}

func HashApr1(password string) string {
	s, err := apr1_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		panic(err)
	}
	return s
}
