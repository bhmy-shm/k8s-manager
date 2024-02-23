package service

type IngressCode uint64

var IngressTag map[IngressCode]string

const (
	IngressCodeCROS IngressCode = iota //跨域
	IngressCodeRewrite
	IngressCodeBasicAuth
	IngressCodeAuthSecret
	IngressCodeRateLimit
	IngressCodeRateLimitBurst
	IngressCodeServerSnippet
)

const (
	IngressOptionsCROSTAG        = "nginx.ingress.kubernetes.io/enable-cors"            //跨域功能
	IngressOptionsREWRITETAG     = "nginx.ingress.kubernetes.io/rewrite-target"         //路径重写
	IngressOptionsBasicAuth      = "nginx.ingress.kubernetes.io/auth-type"              //basic auth 身份验证类型
	IngressOptionsAuthSecret     = "nginx.ingress.kubernetes.io/auth-secret"            //auth secret 身份验证密文，名称
	IngressOptionsRateLimit      = "nginx.ingress.kubernetes.io/limit-rps"              //基于ip进行每秒限流，每秒请求数量
	IngressOptionsRateLimitBurst = "nginx.ingress.kubernetes.io/limit-burst-multiplier" //限流突发请求

	IngressOptionsServerSnippet        = "nginx.ingress.kubernetes.io/server-snippet"
	IngressOptionsConfigurationSnippet = "nginx.ingress.kubernetes.io/configuration-snippet"
)

func init() {
	IngressTag = map[IngressCode]string{
		IngressCodeCROS:           IngressOptionsCROSTAG,
		IngressCodeRewrite:        IngressOptionsREWRITETAG,
		IngressCodeBasicAuth:      IngressOptionsBasicAuth,
		IngressCodeAuthSecret:     IngressOptionsAuthSecret,
		IngressCodeRateLimit:      IngressOptionsRateLimit,
		IngressCodeRateLimitBurst: IngressOptionsRateLimitBurst,
	}
}

func findTage(code IngressCode) string {
	return IngressTag[code]
}
