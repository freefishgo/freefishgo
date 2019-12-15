package routers

import (
	"freeFishGo/examples/fishgo"
	"freeFishGo/router"
)

func init() {
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	fishgo.Mvc.AddMainRouter(&router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}"})
}
