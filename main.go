package main

import (
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

type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}

func (c *ctrTest) MyControllerActionStrut(Test *Test) {
	c.HttpContext.Response.Write([]byte(Test.T1))
}
func main() {
	app := NewFreeFish()
	app.AddHanlers(&ctrTest{})
	app.Run()
	time.Sleep(time.Hour)
}
