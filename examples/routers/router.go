package routers

import (
	"github.com/freefishgo/freefishgo/middlewares/mvc"
)

func init() {
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	mvc.AddMainRouter(&mvc.MainRouter{RouterPattern: "/{ Controller}/{Action}", HomeController: "Main", IndexAction: "LayoutTest"})
}
