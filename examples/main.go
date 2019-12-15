package main

import (
	"freeFishGo"
	_ "freeFishGo/examples/controllers"
	"freeFishGo/examples/fishgo"
	"freeFishGo/examples/middlewares"
	_ "freeFishGo/examples/routers"
)

var build *freeFishGo.ApplicationBuilder

func init() {
	build = freeFishGo.NewFreeFishApplicationBuilder()
}

func main() {
	// 通过注册中间件来实现注册服务
	build.UseMiddleware(&middlewares.Mid{})
	// 把mvc实例注册到管道中
	build.UseMiddleware(fishgo.Mvc)
	build.Run()
}
