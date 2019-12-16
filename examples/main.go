package main

import (
	"github.com/freeFishGo"
	_ "github.com/freeFishGo/examples/controllers"
	"github.com/freeFishGo/examples/fishgo"
	"github.com/freeFishGo/examples/middlewares"
	_ "github.com/freeFishGo/examples/routers"
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
	build.Config.Listen.HTTPPort = 8080
	build.Run()
}
