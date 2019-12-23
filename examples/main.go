package main

import (
	"encoding/json"
	"github.com/freeFishGo"
	appConfig "github.com/freeFishGo/config"
	_ "github.com/freeFishGo/examples/controllers"
	"github.com/freeFishGo/examples/fishgo"
	_ "github.com/freeFishGo/examples/routers"
	"github.com/freeFishGo/middlewares/httpToHttps"
	"github.com/freeFishGo/middlewares/mvc/router"
	"github.com/freeFishGo/middlewares/printTimeMiddleware"
	"os"
)

var build *freeFishGo.ApplicationBuilder

type config struct {
	*appConfig.Config
	WebConfig *router.WebConfig
}

func init() {
	build = freeFishGo.NewFreeFishApplicationBuilder()
	conf := new(config)
	f, _ := os.Open("conf/app.conf")
	json.NewDecoder(f).Decode(conf)
	build.Config = conf.Config
	fishgo.Mvc.Config = conf.WebConfig

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
