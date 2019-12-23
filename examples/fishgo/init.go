package fishgo

import (
	"github.com/freeFishGo/middlewares/mvc"
)

var Mvc *mvc.MvcApp

// 实例化一个mvc服务
func init() {
	Mvc = mvc.NewFreeFishMvcApp()
}
