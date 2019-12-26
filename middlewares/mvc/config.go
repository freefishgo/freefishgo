package mvc

import "github.com/freefishgo/freeFishGo"

type MvcWebConfig struct {
	// 是否启用开发模式
	IsDevelopment bool
	// html模板引擎变量左标记符号
	TemplateLeft string
	// html模板引擎变量右标记符号
	TemplateRight string
	// html模板的父目录
	ViewsPath string

	homeController string
	indexAction    string
	// 是否在Mvc框架最末端捕获Panic，以取代Mvc框架的处理Panic函数
	RecoverPanic bool
	// 捕获Panic的处理函数
	RecoverFunc func(ctx *freeFishGo.HttpContext, e error, Stack []byte)
}

// 实例化一个 MvcWebConfig
func NewWebConfig() *MvcWebConfig {
	return &MvcWebConfig{ViewsPath: "", TemplateLeft: "{{", TemplateRight: "}}", IsDevelopment: false}
}
