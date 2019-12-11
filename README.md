# freeFishGo
golang 通过结构体反射实现的典型的mvc架构
```go
package main

import (
	"fmt"
	free "freeFishGo"
	"freeFishGo/httpContext"
	"freeFishGo/router"
	"time"
)
// 实现mvc控制器的处理ctrTest为控制器 {Controller}的值
type ctrTestController struct {
	router.Controller
}
// GetControllerActionInfo()特殊定制指定action的路由和请求method  默认为 httpContext.GET
func (c *ctrTestController) GetControllerActionInfo() []*router.ControllerActionInfo {
	println("不是默认GetControllerInfo")
	tmp := make([]*router.ControllerActionInfo, 0)
	// 修改MyControllerActionStrut() 方法的路由规则为/{ Controller}/{Action}/{allString}即为/ctrTest/MyControllerActionStrut/{allString}
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}/{allString}", ControllerActionFuncName: "MyControllerActionStrut"})
	return tmp
}
// 作为 Action的请求参数的映射值
type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}
// MyControllerActionStrut为{Action}的值 该方法的默认路由为/ctrTest/MyControllerActionStrut
func (c *ctrTestController) MyControllerActionStrut(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}
// 新添加一个Action 该方法的路由为/ctrTest/My
func (c *ctrTestController) My(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}
// 例子： 组装一个Middleware服务，实现打印mvc框架处理请求的时间
type mid struct {
	
}

func (m *mid) Middleware(ctx *httpContext.HttpContext, next *free.MiddlewareLink) *httpContext.HttpContext {
	dt :=time.Now()
	ctxtmp:= next.Next(ctx)
	fmt.Println("处理时间为:"+(time.Now().Sub(dt)).String())
	return ctxtmp
}

func (m *mid) LastInit() {
	//panic("implement me")
}


func main() {
	// 实例化一个mvc服务
	app := free.NewFreeFishMvcApp()
	// 注册控制器
	app.AddHanlers(&ctrTestController{})
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	app.AddMainRouter(&router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}"})
	build:= free.NewFreeFishApplicationBuilder()
	// 通过注册中间件来实现注册服务
	build.UseMiddleware(&mid{})
	// 把mvc实例注册到管道中
	build.UseMiddleware(app)
	build.Run()
}

```
