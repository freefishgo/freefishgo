package main

import (
	"freeFishGo/router"
)

type ctrTest struct {
	router.Controller
}

func (c *ctrTest) GetControllerActionInfo() []*router.ControllerActionInfo {
	println("不是默认GetControllerInfo")
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}/{id:int}{string}/{int}int", ControllerActionFuncName: "MyControllerActionStrut"})
	return tmp
}

type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}

func (c *ctrTest) MyControllerActionStrut(Test *Test) {
	c.HttpContext.Response.Write([]byte(Test.Id))
}
func main() {
	app := NewFreeFish()
	// 注册控制器
	app.AddHanlers(&ctrTest{})
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	app.AddMainRouter(&router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}"})
	app.Run()
}
