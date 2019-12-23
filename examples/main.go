package main

import (
	"github.com/freeFishGo/examples/conf"
	_ "github.com/freeFishGo/examples/controllers"
	"github.com/freeFishGo/examples/fishgo"
	_ "github.com/freeFishGo/examples/routers"
	"github.com/freeFishGo/middlewares/printTimeMiddleware"
)

func main() {
	// 通过注册中间件来打印任务处理时间服务
	conf.Build.UseMiddleware(&printTimeMiddleware.PrintTimeMiddleware{})
	// 利用中间件来实现http到https的转换
	//conf.Build.UseMiddleware(&httpToHttps.HttpToHttps{})
	// 把mvc实例注册到管道中
	conf.Build.UseMiddleware(fishgo.Mvc)
	conf.Build.Config.Listen.HTTPPort = 8080
	conf.Build.Run()
}
