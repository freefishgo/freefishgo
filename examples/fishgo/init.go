package fishgo

import (
	"github.com/freeFishGo/middlewares"
)

var Mvc *middlewares.MvcApp

// 实例化一个mvc服务
func init() {
	Mvc = middlewares.NewFreeFishMvcApp()
}
