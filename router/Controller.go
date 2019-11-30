package router

import (
	"freeFishGo/httpContext"
	"reflect"
	"strings"
)

type ControllerInfo struct {
	ControllerFunc   reflect.Type //请求事件的处理函数
	ControllerName   string       //控制器名称
	ControllerAction string       //控制器处理方法
	ParameterStruct  reflect.Type
}

// http请求逻辑控制器
type Controller struct {
	HttpContext    *httpContext.HttpContext
	controllerInfo *ControllerInfo
	sonController  IController
}

type IController interface {
	getControllerInfo(*tree) *tree
	setSonController(IController)
	GetControllerInfo() *ControllerInfo
	SetHttpContext(ctx *httpContext.HttpContext)
}

// 进行路由注册的基类 如果结构体含有Controller 则Controller去掉 如GetController 变位Get  忽略大小写
func (c *Controller) getControllerInfo(tree *tree) *tree {
	getType := reflect.TypeOf(c.sonController)
	controllerNameList := strings.Split(getType.String(), ".")
	controllerName := controllerNameList[len(controllerNameList)-1]
	for i := 0; i < getType.NumMethod(); i++ {
		me := getType.Method(i)
		actionName := me.Name
		if isNotSkin(actionName) {
			tree.addPathTree(controllerName, actionName, getType.Elem())
		} else {
			continue
		}
		if me.Type.NumIn() == 2 {
			tmp := me.Type.In(1)
			if tmp.Kind() == reflect.Struct {
				for i := 0; i < tmp.NumField(); i++ {
					field := tmp.Field(i)
					println(field.Tag)
					println(field.Name)
				}
			} else {
				panic("方法" + getType.String() + "." + me.Name + "错误:不能设置参数为非结构体,且只能设置一个结构体")
			}
		}
	}

	(c.sonController).GetControllerInfo()
	c.controllerInfo = new(ControllerInfo)
	return tree
}

func (c *Controller) GetControllerInfo() *ControllerInfo {
	println("默认GetControllerInfo")
	return new(ControllerInfo)
}

func (c *Controller) setSonController(son IController) {
	c.sonController = son
}
func (c *Controller) SetHttpContext(ctx *httpContext.HttpContext) {
	c.HttpContext = ctx
}

func isNotSkin(methodName string) bool {
	skinList := map[string]bool{"SetHttpContext": true,
		"GetControllerInfo": true}
	if _, ok := skinList[methodName]; ok {
		return false
	}
	return true
}
