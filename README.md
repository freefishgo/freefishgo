# freeFishGo
golang 通过结构体反射实现的典型的mvc架构 尝试可以看代码文件中examples/main.go

# 详细文档地址访问 http://freefishgo.com

## Installation

To install `freefishgo` use the `go get` command:

```bash
go get github.com/freefishgo/freefishgo
```

> If you already have `freefishgo` installed, updating `freefishgo` is simple:

```bash
go get -u github.com/freefishgo/freefishgo
```
> 如果你想快速构建项目 请查看 https://github.com/freefishgo/freefish

## 使用案例

```go
package main

import (
	"fmt"
	"github.com/freefishgo/freeFishGo"
	"github.com/freefishgo/freeFishGo/middlewares/mvc"
	"log"
	"time"
)
// 实现mvc控制器的处理ctrTest为控制器 {Controller}的值
type MainController struct {
	mvc.Controller
}
// OverwriteRouter()特殊定制指定action的路由
func (c *MainController) OverwriteRouter() []*mvc.ControllerActionInfo {
	log.Println("不是默认GetControllerInfo")
	tmp := make([]*mvc.ControllerActionInfo, 0)
	tmp = append(tmp, &mvc.ControllerActionInfo{RouterPattern: "{string}/{ Controller}/{Action}/{tstst1:string}er", ControllerActionFuncName: "MyControllerActionStrutPost"})
	return tmp
}
// 作为 Action的请求参数的映射值
type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}
func (c *MainController) MyControllerActionStrutPost(t *Test) {
	c.Data["Website"] = t.Id
	c.Data["Email"] = t.T1
	// 调用模板引擎   默认模板地址为{ Controller}/{Action}.fish    不含请求方式
	c.UseTplPath()
}
func (c *MainController) MyControllerActionStrutGet(t *Test) {
	c.Data["Website"] = t.Id
	c.Data["Email"] = t.T1
	//c.HttpContext.Response.Write([]byte("hahaha"))
	c.UseTplPath()
}
func (c *MainController) MyGET(t *Test) {
	c.Response.Write([]byte(fmt.Sprintf("数据为：%+v", t)))
}
func (c *MainController) My1(t *Test) {
	c.Response.Write([]byte(fmt.Sprintf("数据为：%+v", t)))
}
// 例子： 组装一个Middleware服务，实现打印mvc框架处理请求的时间
type mid struct {
	
}
// 中间件打印mvc框架处理请求的时间
func (*mid) Middleware(ctx *freeFishGo.HttpContext, next freeFishGo.Next) *freeFishGo.HttpContext {
		dt := time.Now()
    	log.Println(ctx.Request.URL)
    	ctxtmp := next(ctx)
    	log.Println("处理时间为:" + (time.Now().Sub(dt)).String())
    	return ctxtmp
}
// 中间件注册是调用函数进行该中间件最后的设置
func (*mid) LastInit(*freeFishGo.Config) {
	panic("implement me")
}
func main() {
	// 实例化一个mvc服务
	app := mvc.NewFreeFishMvcApp()
	// 注册控制器
	app.AddHandlers(&MainController{})
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	app.AddMainRouter(&mvc.MainRouter{RouterPattern: "/{ Controller}/{Action}", HomeController: "Main", IndexAction: "My"})
	build:= freeFishGo.NewFreeFishApplicationBuilder()
	// 通过注册中间件来实现注册服务
	build.UseMiddleware(&mid{})
	// 把mvc实例注册到管道中
	build.UseMiddleware(app)
	build.Run()
}

```
