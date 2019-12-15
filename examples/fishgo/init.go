package fishgo

import (
	"freeFishGo"
)

var Mvc *freeFishGo.MvcApp

// 实例化一个mvc服务
func init() {
	Mvc = freeFishGo.NewFreeFishMvcApp()
}
