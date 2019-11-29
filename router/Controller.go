package router

import (
	"freeFishGo/httpContext"
	"reflect"
	"strings"
)

type ControllerInfo struct {
	Path string
}

// http请求逻辑控制器
type Controller struct {
	HttpContext    *httpContext.HttpContext
	controllerInfo *ControllerInfo
	sonController  *IController
}

type IController interface {
	getControllerInfo(*tree) *tree
	setSonController(*IController)
	GetControllerInfo() *ControllerInfo
}

// 进行路由注册的基类 如果结构体含有Controller 则Controller去掉 如GetController 变位Get  忽略大小写
func (c *Controller) getControllerInfo(tree *tree) *tree {
	getType := reflect.TypeOf(*c.sonController)
	controllerNameList := strings.Split(getType.String(), ".")
	controllerName := controllerNameList[len(controllerNameList)-1]
	println(controllerName)
	getValue := reflect.ValueOf(*c.sonController)
	for i := 0; i < getType.NumMethod(); i++ {
		me := getType.Method(i)
		actionName := me.Name
		println(actionName)
		println(getValue.Method(i).String())
	}
	(*c.sonController).GetControllerInfo()
	c.controllerInfo = new(ControllerInfo)
	return tree
}

func (c *Controller) GetControllerInfo() *ControllerInfo {
	println("默认GetControllerInfo")
	return new(ControllerInfo)
}

func (c *Controller) setSonController(son *IController) {
	c.sonController = son
}
