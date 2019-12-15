package fishgo

import (
	"freeFishGo"
	"freeFishGo/router"
)

var Mvc *freeFishGo.MvcApp

// 实例化一个mvc服务
func init() {
	Mvc = freeFishGo.NewFreeFishMvcApp()
}

func AddHandlers(ic ...router.IController) {
	Mvc.AddHandlers(ic...)
}

func AddMainRouter(list ...*router.ControllerActionInfo) {
	Mvc.AddMainRouter(list...)
}
