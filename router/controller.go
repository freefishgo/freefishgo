package router

import (
	"freeFishGo/httpContext"
	"reflect"
	"regexp"
	"strings"
)

type ControllerInfo struct {
	ControllerFunc                  reflect.Type //请求事件的处理函数
	ControllerName                  string       //控制器名称
	ControllerAction                string       //控制器处理方法
	ControllerActionParameterStruct reflect.Type
}

// http请求逻辑控制器
type Controller struct {
	HttpContext    *httpContext.HttpContext
	controllerInfo *ControllerInfo
	sonController  IController
}

// 控制器的基本数据结构
type IController interface {
	getControllerInfo(*tree) *tree
	setSonController(IController)
	GetControllerActionInfo() []*ControllerActionInfo
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
		if !isNotSkin(actionName) {
			continue
		}
		var controllerActionParameterStruct reflect.Type = nil
		if me.Type.NumIn() == 2 {
			tmp := me.Type.In(1)
			if tmp.Kind() == reflect.Ptr {
				if tmp.Elem().Kind() == reflect.Struct {
					controllerActionParameterStruct = tmp.Elem()
				} else {
					panic("方法" + getType.String() + "." + me.Name + "错误:只能传结构体指针,且只能设置一个结构体指针")
				}
			} else {
				panic("方法" + getType.String() + "." + me.Name + "错误:只能传结构体指针,且只能设置一个结构体指针")
			}
		}
		tree.addPathTree(controllerName, actionName, getType.Elem(), controllerActionParameterStruct)
	}
	controllerActionInfoList := (c.sonController).GetControllerActionInfo()
	for _, v := range controllerActionInfoList {
		_, ok := getType.MethodByName(v.ControllerActionFuncName)
		if !ok {
			panic(getType.String() + "方法" + v.ControllerActionFuncName + "不存在")
		}
		v.actionName = strings.ToLower(v.ControllerActionFuncName)
		v.controllerName = strings.ToLower(controllerName)
		f := regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
		v.RouterPattern = f.ReplaceAllString(v.RouterPattern, v.controllerName)
		f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
		v.RouterPattern = f.ReplaceAllString(v.RouterPattern, v.actionName)
		tree.ControllerModelList = tree.ControllerModelList.AddControllerModelList(v)
	}
	return tree
}

type ControllerActionInfo struct {
	// 传设置控制器的方法
	ControllerActionFuncName string
	//路由设置  如：/{Controller}/{Action}/{id:int}
	// /home/index/123可以匹配成功
	RouterPattern string
	//允许的请求方法
	AllowMethod    []httpContext.HttpMethod
	allowMethod    map[httpContext.HttpMethod]bool
	controllerName string
	actionName     string
	patternRe      *regexp.Regexp
	//正则匹配出来的变量地址映射变量映射
	patternMap map[string]int
}

// 控制器属性设置 路由变量路由中只能出现一次
func (c *Controller) GetControllerActionInfo() []*ControllerActionInfo {
	println("默认GetControllerInfo")
	return make([]*ControllerActionInfo, 0)
}

// 控制器注册
func (c *Controller) setSonController(son IController) {
	c.sonController = son
}

// http请求上下文注册
func (c *Controller) SetHttpContext(ctx *httpContext.HttpContext) {
	c.HttpContext = ctx
}

// 过滤掉本地方法
func isNotSkin(methodName string) bool {
	skinList := map[string]bool{"SetHttpContext": true,
		"GetControllerInfo": true}
	if _, ok := skinList[methodName]; ok {
		return false
	}
	return true
}
