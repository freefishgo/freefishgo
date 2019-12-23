package main

import (
	"github.com/freeFishGo"
	_ "github.com/freeFishGo/examples/controllers"
	"github.com/freeFishGo/examples/fishgo"
	_ "github.com/freeFishGo/examples/routers"
	"github.com/freeFishGo/middlewares/httpToHttps"
	"github.com/freeFishGo/middlewares/printTimeMiddleware"
)

var build *freeFishGo.ApplicationBuilder

func init() {
	build = freeFishGo.NewFreeFishApplicationBuilder()
}

func main() {
	// 通过注册中间件来打印任务处理时间服务
	build.UseMiddleware(&printTimeMiddleware.PrintTimeMiddleware{})
	// 利用中间件来实现http到https的转换
	build.UseMiddleware(&httpToHttps.HttpToHttps{})
	// 把mvc实例注册到管道中
	build.UseMiddleware(fishgo.Mvc)
	build.Config.Listen.HTTPPort = 8080
	build.Run()
}
