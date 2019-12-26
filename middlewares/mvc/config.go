package mvc

import "github.com/freefishgo/freeFishGo"

type WebConfig struct {
	IsDevelopment bool
	TemplateLeft  string
	TemplateRight string
	ViewsPath     string

	homeController string
	indexAction    string

	RecoverPanic bool
	RecoverFunc  func(ctx *freeFishGo.HttpContext, e error, Stack []byte)
}

func NewWebConfig() *WebConfig {
	return &WebConfig{ViewsPath: "", TemplateLeft: "{{", TemplateRight: "}}", IsDevelopment: false}
}
