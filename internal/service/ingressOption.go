package service

type IngressCode uint64

var IngressTag map[IngressCode]string

const (
	IngressCodeCROS IngressCode = iota //跨域
	IngressCodeRewrite
	IngressCodeLimit //
)

const (
	IngressOptionsCROSTAG    = "nginx.ingress.kubernetes.io/enable-cors"
	IngressOptionsREWRITETAG = "nginx.ingress.kubernetes.io/rewrite-target"
)

func init() {
	IngressTag = map[IngressCode]string{
		IngressCodeCROS:    IngressOptionsCROSTAG,
		IngressCodeRewrite: IngressOptionsREWRITETAG,
	}
}

func findTage(code IngressCode) string {
	return IngressTag[code]
}
