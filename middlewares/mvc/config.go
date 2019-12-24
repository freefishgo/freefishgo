package mvc

import "github.com/freeFishGo/httpContext"

type WebConfig struct {
	TemplateLeft  string
	TemplateRight string
	ViewsPath     string
	RecoverPanic  bool
	RecoverFunc   func(ctx *httpContext.HttpContext, e error, Stack []byte)
}

func NewWebConfig() *WebConfig {
	return &WebConfig{ViewsPath: "", TemplateLeft: "{{", TemplateRight: "}}"}
}
