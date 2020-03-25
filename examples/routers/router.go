package routers

import (
	_ "github.com/freefishgo/freefishgo/examples/controllers"
	"github.com/freefishgo/freefishgo/middlewares/mvc"
)

func init() {
	// 设置主路由格式
	mvc.AddMainRouter(&mvc.MainRouter{RouterPattern: "/{ Controller}/{Action}", HomeController: "Main", IndexAction: "LayoutTest"})
}
