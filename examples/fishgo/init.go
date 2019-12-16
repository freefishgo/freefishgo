package fishgo

import (
	"github.com/freeFishGo"
)

var Mvc *freeFishGo.MvcApp

// 实例化一个mvc服务
func init() {
	Mvc = freeFishGo.NewFreeFishMvcApp()
}
