package main

import (
	"freeFishGo/router"
	"time"
)

type ctrTest struct {
	router.Controller
}

//func (c *ctrTest) GetControllerInfo() *router.ControllerInfo {
//	println("进入自定义的了")
//	return nil
//}

func main() {
	app := NewFreeFish()
	app.AddHanlers(&ctrTest{})
	app.Run()
	time.Sleep(time.Hour)
}
