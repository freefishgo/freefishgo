package main

import (
	"freeFishGo/httpContext"
	"freeFishGo/router"
	"time"
)

type ctrTest struct {
	router.Controller
}

func (c *ctrTest) GetControllerInfo() *router.ControllerInfo {
	println("进入自定义的了")
	return nil
}
func (c *ctrTest) MyControllerActionStrut(Test httpContext.HttpContext) {
	c.HttpContext.Response.Write([]byte("MyControllerAction"))
}
func main() {
	app := NewFreeFish()
	app.AddHanlers(&ctrTest{})
	app.Run()
	time.Sleep(time.Hour)
}
