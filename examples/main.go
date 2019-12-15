package main

import (
	free "freeFishGo"
	"freeFishGo/examples/controllers"
	"freeFishGo/examples/middlewares"
	"freeFishGo/router"
)

func main() {
	// 实例化一个mvc服务
	app := free.NewFreeFishMvcApp()
	// 注册控制器
	app.AddHandlers(&controllers.MainController{})
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	app.AddMainRouter(&router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}"})
	build := free.NewFreeFishApplicationBuilder()
	// 通过注册中间件来实现注册服务
	build.UseMiddleware(&middlewares.Mid{})
	// 把mvc实例注册到管道中
	build.UseMiddleware(app)
	build.Run()
}
