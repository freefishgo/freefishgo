package main

import (
	"freeFishGo/router"
	"time"
)

type ctrTest struct {
	router.Controller
}

func (c *ctrTest) GetControllerActionInfo() []*router.ControllerActionInfo {
	println("不是默认GetControllerInfo")
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}/{id:string}", ControllerActionFuncName: "MyControllerActionStrut"})
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
	app.AddHanlers(&ctrTest{})
	app.Run()
	time.Sleep(time.Hour)
}
