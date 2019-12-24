package mvc

import "github.com/freefishgo/freeFish"

type WebConfig struct {
	TemplateLeft  string
	TemplateRight string
	ViewsPath     string
	RecoverPanic  bool
	RecoverFunc   func(ctx *freeFish.HttpContext, e error, Stack []byte)
}

func NewWebConfig() *WebConfig {
	return &WebConfig{ViewsPath: "", TemplateLeft: "{{", TemplateRight: "}}"}
}
