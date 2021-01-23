# freeFishGo(1分钟学会golang写后端接口)
## freeFishGo是通过结构体反射实现的典型的mvc的web框架
## FreeFishGo is a typical MVC web framework implemented through struct reflection

## 使用案例
```go
package main

import (
	"fmt"
	"github.com/freefishgo/freefishgo"
	"github.com/freefishgo/freefishgo/middlewares/mvc"
	"log"
	"time"
)
// 实现mvc控制器的处理ctrTest为控制器 {Controller}的值
type MainController struct {
	mvc.Controller
}

// 作为 Action的请求参数的映射值
type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}
// 路由地址为 /main/MyControllerActionStrut 请求方式为post
func (c *MainController) MyControllerActionStrutPost(t *Test) {
	c.Data["Website"] = t.Id
	c.Data["Email"] = t.T1
	// 调用模板引擎   默认模板地址为{ Controller}/{Action}.fish    不含请求方式
	c.UseTplPath()
}
// 路由地址为 /main/MyControllerActionStrut 请求方式为get
func (c *MainController) MyControllerActionStrutGet(t *Test) {
	c.Data["Website"] = t.Id
	c.Data["Email"] = t.T1
	//c.HttpContext.Response.Write([]byte("hahaha"))
	c.UseTplPath()
}
// 路由地址为 /main/My 请求方式为get
func (c *MainController) MyGET(t *Test) {
	c.Response.Write([]byte(fmt.Sprintf("数据为：%+v", t)))
}
// 路由地址为 /main/My1 请求方式为get
func (c *MainController) My1(t *Test) {
	c.Response.Write([]byte(fmt.Sprintf("数据为：%+v", t)))
}
func main() {
	// 注册控制器
	mvc.AddHandlers(&MainController{})
	// 把mvc实例注册到管道中
	freefishgo.UseMiddleware(mvc.GetDefaultMvcApp())
	freefishgo.Run()
}
```
# 如果你想快速构建项目 请查看 https://github.com/freefishgo/freefish
# 使用freefishgo的优势(详细文档地址访问 http://freefishgo.com)
1.通过的中间件注入对http请求进行流处理 清楚明白处理流程<br/>
2.自定义中间件只需要实现IMiddleware接口 该接口只有两个方法,方便快捷,实现任意地点对请求的处理<br/>
3.提供Mvc服务也是通过中间件注入的 该中间件能通过动作器的参数进行参数注入,免去忘记请求参数的烦恼,再也不用再写反序列化获取参数了,框架帮你完成 <br/>
4.url中任意位置的字符串也能注入到请求参数的<br/>
5.你还在为每一个动作器都需要写路由地址而烦恼吗？ freefishgo解除你的烦恼,你可以不进行路由设置,实现所有控制器动作器路由注入 <br/>

1. Flow processing of HTTP requests through middleware injection to clearly understand the processing process<br/>
2. Custom middleware only needs to implement IMiddleware interface. This interface has only two methods, which are convenient and fast, and can handle requests at any place<br/>
3. The Mvc service is also injected through the middleware. This middleware can inject parameters through the parameters of the actor, so as to avoid the trouble of forgetting the request parameters, and no longer need to write deserialization to get parameters<br/>
4. Strings at any location in the url can also be injected into the request parameters<br/>
5. Are you still worried that every actuator needs to write a routing address? Freefishgo takes care of your worries, you can implement all controller actions without routing Settings<br/>
## Installation

To install `freefishgo` use the `go get` command:

```bash
go get github.com/freefishgo/freefishgo
```

> If you already have `freefishgo` installed, updating `freefishgo` is simple:

```bash
go get -u github.com/freefishgo/freefishgo
```

