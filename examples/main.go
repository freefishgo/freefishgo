package main

import (
	"fmt"
	free "freeFishGo"
	"freeFishGo/httpContext"
	"freeFishGo/router"
	"log"
	"time"
)

type ctrTestController struct {
	router.Controller
}

func (c *ctrTestController) GetControllerActionInfo() []*router.ControllerActionInfo {
	log.Println("不是默认GetControllerInfo")
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "{string}/{ Controller}/{Action}/{tstst1:string}er", ControllerActionFuncName: "MyControllerActionStrut"})
	return tmp
}

type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}

func (c *ctrTestController) MyControllerActionStrut(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	//c.HttpContext.Response.Write([]byte("hahaha"))
	c.UseTplPath()
}

func (c *ctrTestController) My(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}
func (c *ctrTestController) My1(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}

type mid struct {
}

func (m *mid) Middleware(ctx *httpContext.HttpContext, next free.Next) *httpContext.HttpContext {
	dt := time.Now()
	log.Println(ctx.Request.URL)
	ctxtmp := next(ctx)
	log.Println("处理时间为:" + (time.Now().Sub(dt)).String())
	return ctxtmp
}

func (m *mid) LastInit() {
	//panic("implement me")
}

func main() {
	app := free.NewFreeFishMvcApp()
	// 注册控制器
	app.AddHanlers(&ctrTestController{})
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	app.AddMainRouter(&router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}"})
	build := free.NewFreeFishApplicationBuilder()
	build.UseMiddleware(&mid{})
	build.UseMiddleware(app)
	build.Run()
}
