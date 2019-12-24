package fishgo

import (
	"github.com/freefishgo/freeFish/middlewares/mvc"
)

var Mvc *mvc.MvcApp

// 实例化一个mvc服务
func init() {
	Mvc = mvc.NewFreeFishMvcApp()
}
