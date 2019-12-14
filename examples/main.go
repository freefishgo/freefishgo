package main

import (
	"fmt"
	free "freeFishGo"
	"freeFishGo/httpContext"
	"freeFishGo/router"
	"log"
	"time"
)

// 实现mvc控制器的处理ctrTest为控制器 {Controller}的值
type ctrTestController struct {
	router.Controller
}

// GetControllerActionInfo()特殊定制指定action的路由
func (c *ctrTestController) GetControllerActionInfo() []*router.ControllerActionInfo {
	log.Println("不是默认GetControllerInfo")
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "{string}/{ Controller}/{Action}/{tstst1:string}er", ControllerActionFuncName: "MyControllerActionStrutPost"})
	return tmp
}

// 作为 Action的请求参数的映射值
type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/ctrTest/MyControllerActionStrut 最后的单词为请求方式  该例子为Post请求
func (c *ctrTestController) MyControllerActionStrutPost(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	// 调用模板引擎   默认模板地址为{ Controller}/{Action}.fish    不含请求方式
	c.UseTplPath()
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/ctrTest/MyControllerActionStrut 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看
func (c *ctrTestController) MyControllerActionStrutGet(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	//c.HttpContext.Response.Write([]byte("hahaha"))
	c.UseTplPath()
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/ctrTest/My 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看
func (c *ctrTestController) MyGET(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/ctrTest/My1 get请求可以省略get后缀  查询具体字符串值可到httpContext包中查看
func (c *ctrTestController) My1(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}

// 例子： 组装一个Middleware服务，实现打印mvc框架处理请求的时间
type mid struct {
}

// 中间件打印mvc框架处理请求的时间
func (m *mid) Middleware(ctx *httpContext.HttpContext, next free.Next) *httpContext.HttpContext {
	dt := time.Now()
	log.Println(ctx.Request.URL)
	ctxtmp := next(ctx)
	log.Println("处理时间为:" + (time.Now().Sub(dt)).String())
	return ctxtmp
}

// 中间件注册是调用函数进行该中间件最后的设置
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
	build := free.NewFreeFishApplicationBuilder()
	// 通过注册中间件来实现注册服务
	build.UseMiddleware(&mid{})
	// 把mvc实例注册到管道中
	build.UseMiddleware(app)
	build.Run()
}
