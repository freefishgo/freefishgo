package main

import (
	"encoding/json"
	"fmt"
	"github.com/freefishgo/freefishgo"
	_ "github.com/freefishgo/freefishgo/examples/routers"
	"github.com/freefishgo/freefishgo/middlewares/mvc"
	"github.com/freefishgo/freefishgo/middlewares/printTimeMiddleware"
	"os"
)

type config struct {
	*freefishgo.Config
	WebConfig *mvc.MvcWebConfig
}

func init() {
	conf := new(config)
	f, err := os.Open("conf/app.conf")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(f).Decode(conf)
	if err != nil {
		panic(err)
	}
	freefishgo.SetDefaultApplicationBuilderConfig(conf.Config)
	mvc.SetDefaultMvcWebConfig(conf.WebConfig)
}

func main() {
	// 通过注册中间件来打印任务处理时间服务
	freefishgo.UseMiddleware(&printTimeMiddleware.PrintTimeMiddleware{})
	freefishgo.UseMiddlewareFunc(func(ctx *freefishgo.HttpContext, next freefishgo.Next) *freefishgo.HttpContext {
		fmt.Println("MiddlewareFunc")
		return next(ctx)
	})
	// 把mvc实例注册到管道中
	freefishgo.UseMiddleware(mvc.GetDefaultMvcApp())
	freefishgo.Run()

}
