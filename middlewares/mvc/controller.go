// Copyright 2019 freefishgo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mvc

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/freefishgo/freefishgo"
)

type controllerInfo struct {
	ControllerFunc                  reflect.Type //请求事件的处理函数
	ControllerName                  string       //控制器名称
	ControllerAction                string       //控制器处理方法
	ControllerActionParameterStruct reflect.Type
}

// Controller http请求逻辑控制器
type Controller struct {
	// 响应前端的处理 不建议使用
	Response freefishgo.IResponse
	// 重置控制器路由  必须包含{Action}变量
	ControllerRouter *ControllerRouter
	// 控制器单一方法设置路由
	ActionRouterList []*ActionRouter
	// 和前端一切的数据  都可以通过他获取
	Request        *freefishgo.Request
	controllerInfo *controllerInfo
	sonController  IController
	// 前端传来的数据都可以获取  包括路由格式化的数据  如 /{id:string}  可通过 Query["id"]获取值
	Query map[string]interface{}
	// 模板引擎中变量数据
	Data map[interface{}]interface{}
	// 如果母版页存在 则该内容会被填充到模板页的 .LayoutContent 变量中
	isStopController bool
	tplPath          string
	isUseTplPath     bool
	controllerName   string
	actionName       string
	// 母版页地址
	LayoutPath string
	//母版页子页面地址
	LayoutSections map[string]string
}

func (c *Controller) setQuery(m map[string]interface{}) {
	c.Query = m
}

// UseTplPath 使用模板的路径并启动调用模板
//
// 管道中任意地方调用Controller.HttpContext.Response.Write() 方法会失效
//
// 不使用路径会用v.ViewsPath/{Controller}/{Action}.fish
//
// 多路径只用最后一个
func (c *Controller) UseTplPath(tplPath ...string) {
	c.isUseTplPath = true
	len := len(tplPath)
	if len != 0 {
		c.tplPath = tplPath[len-1]
	}
}

func (c *Controller) getController() *Controller {
	return c
}

// IController 控制器的基本数据结构
type IController interface {
	getControllerInfo(*tree) *tree
	setSonController(IController)
	initController(ctx *freefishgo.HttpContext)
	getController() *Controller
	setQuery(map[string]interface{})
	// UseTplPath 使用模板
	UseTplPath(tplPath ...string)
	// SkipController 跳过控制器
	SkipController()
	// Prepare http控制器处理前调用方法
	Prepare()
	// Finish http控制器执行后调用方法
	Finish()
}

// IStatusCodeController 响应状态处理接口
type IStatusCodeController interface {
	IController
	// Error500 http500处理方法
	Error500()
	// NotFind404 http404处理方法
	NotFind404()
	// Forbidden403 http403处理方法
	Forbidden403()
}

// StatusCodeController http响应状态处理
type StatusCodeController struct {
	Controller
}

// Error500 http 500错误处理
func (s *StatusCodeController) Error500() {
	s.Response.WriteHeader(500)
	fmt.Fprintf(s.Response, `<html><body><div style="color: red;color: red;margin: 150px auto;width: 800px;"><div>500 Internal Server Error:  %s </div><pre>%s</pre></div></body></html>`, s.Response.Error(), s.Response.Stack())
}

// NotFind404 http 404处理
func (s *StatusCodeController) NotFind404() {
	s.Response.Write([]byte("404 page not found"))
}

// Forbidden403 http 403处理
func (s *StatusCodeController) Forbidden403() {
	s.Response.Write([]byte("403 Forbidden"))
}

// getControllerInfo 进行路由注册的基类 如果结构体含有Controller 则Controller去掉 如GetController 变位Get  忽略大小写
func (c *Controller) getControllerInfo(tree *tree) *tree {
	getType := reflect.TypeOf(c.sonController)
	controllerNameList := strings.Split(getType.String(), ".")
	controllerName := controllerNameList[len(controllerNameList)-1]
	for i := 0; i < getType.NumMethod(); i++ {
		me := getType.Method(i)
		actionName := me.Name
		if !isNotSkip(actionName) {
			continue
		}
		var controllerActionParameterStruct reflect.Type = nil
		if me.Type.NumIn() <= 2 {
			if me.Type.NumIn() == 2 {
				tmp := me.Type.In(1)
				if tmp.Kind() == reflect.Ptr {
					if tmp.Elem().Kind() == reflect.Struct {
						controllerActionParameterStruct = tmp.Elem()
					} else {
						panic("方法" + getType.String() + "." + me.Name + "错误:只能传结构体指针或者无参,且只能设置一个结构体指针")
					}
				} else {
					panic("方法" + getType.String() + "." + me.Name + "错误:只能传结构体指针或者无参,且只能设置一个结构体指针")
				}
			}
		}
		f := regexp.MustCompile(`Controller$`)
		controllerName = f.ReplaceAllString(controllerName, "")
		tree.addPathTree(controllerName, actionName, getType.Elem(), controllerActionParameterStruct)
	}
	if tree.CloseMainRouter == nil {
		tree.CloseMainRouter = map[string]map[string]bool{}
	}
	if tree.CloseControllerRouter == nil {
		tree.CloseControllerRouter = map[string]bool{}
	}
	if tree.StaticRouterList == nil {
		tree.StaticRouterList = map[string]*ActionRouter{}

	}
	controllerRouter := c.ControllerRouter
	if controllerRouter != nil {
		v := &ActionRouter{}
		v.RouterPattern = controllerRouter.RouterPattern
		v.controllerName = strings.ToLower(controllerName)
		tree.CloseControllerRouter[v.controllerName] = true
		//tree.CloseControllerRouter[actionRouter.controllerName]=actionRouter
		f := regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
		v.RouterPattern = f.ReplaceAllString(v.RouterPattern, strings.ToLower(controllerName))
		f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
		if !f.MatchString(v.RouterPattern) {
			panic("控制器路由注册时发现：控制器 " + getType.String() + "错误:错误原因为Controller注册时路由规则中必须含有 {Action}变量")
		}
		tree.ControllerRouterList = tree.ControllerRouterList.AddControllerModelList(v)

	}
	controllerActionInfoList := c.ActionRouterList
	if controllerActionInfoList == nil {
		return tree
	}
	for _, v := range controllerActionInfoList {
		_, ok := getType.MethodByName(v.ControllerActionFuncName)
		if !ok {
			panic(getType.String() + "方法" + v.ControllerActionFuncName + "不存在")
		}
		v.actionName = replaceActionName(v.ControllerActionFuncName)
		v.controllerName = strings.ToLower(controllerName)
		if tree.CloseMainRouter[v.controllerName] == nil {
			tree.CloseMainRouter[v.controllerName] = map[string]bool{}
			tree.CloseMainRouter[v.controllerName][v.actionName] = true
		} else {
			tree.CloseMainRouter[v.controllerName][v.actionName] = true
		}
		f := regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
		v.RouterPattern = f.ReplaceAllString(v.RouterPattern, strings.ToLower(controllerName))
		f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
		v.RouterPattern = f.ReplaceAllString(v.RouterPattern, v.actionName)
		if !regexp.MustCompile("{.*?}").MatchString(v.RouterPattern) {
			tree.StaticRouterList = addStaticAction(tree.StaticRouterList, v)
			continue
		}
		tree.ActionRouterList = tree.ActionRouterList.AddControllerModelList(v)
	}
	return tree
}

func replaceActionName(actionName string) string {
	actionName = strings.ToUpper(actionName)
	httpMethodList := []freefishgo.HttpMethod{freefishgo.MethodPost,
		freefishgo.MethodConnect, freefishgo.MethodDelete,
		freefishgo.MethodGet, freefishgo.MethodHead, freefishgo.MethodOptions,
		freefishgo.MethodPatch, freefishgo.MethodPut, freefishgo.MethodTrace}
	for _, v := range httpMethodList {
		f := regexp.MustCompile(string(v) + "$")
		if f.MatchString(actionName) {
			return strings.ToLower(f.ReplaceAllString(actionName, ""))
		}
	}
	return strings.ToLower(actionName)

}

func replaceActionNameIgnore(actionName string) string {
	httpMethodList := []freefishgo.HttpMethod{freefishgo.MethodPost,
		freefishgo.MethodConnect, freefishgo.MethodDelete,
		freefishgo.MethodGet, freefishgo.MethodHead, freefishgo.MethodOptions,
		freefishgo.MethodPatch, freefishgo.MethodPut, freefishgo.MethodTrace}
	for _, v := range httpMethodList {
		f, _ := regexp.Compile("(?i:" + strings.ToLower(string(v)) + ")$")
		if f.MatchString(actionName) {
			return f.ReplaceAllString(actionName, "")
		}
	}
	return actionName

}

func isHaveHttpMethod(actionName string) bool {
	actionName = strings.ToUpper(actionName)
	httpMethodList := []freefishgo.HttpMethod{freefishgo.MethodPost,
		freefishgo.MethodConnect, freefishgo.MethodDelete,
		freefishgo.MethodGet, freefishgo.MethodHead, freefishgo.MethodOptions,
		freefishgo.MethodPatch, freefishgo.MethodPut, freefishgo.MethodTrace}
	for _, v := range httpMethodList {
		f := regexp.MustCompile(string(v) + "$")
		if f.MatchString(actionName) {
			return true
		}
	}
	return false

}

// ActionRouter 单一动作器路由设置结构体
type ActionRouter struct {
	// 传设置控制器的方法
	ControllerActionFuncName string
	//路由设置  如：/{Controller}/{Action}/{id:int}
	// /home/index/123可以匹配成功
	RouterPattern  string
	controllerName string
	actionName     string
	// 正则表达式匹配
	PatternRe *regexp.Regexp
	//正则匹配出来的变量地址映射变量映射
	patternMap map[string]int
}

// ControllerRouter 单一控制器路由设置结构体，路由规则中必须包含`{Action}`变量
type ControllerRouter struct {
	//路由设置  如：/{Controller}/{Action}/{id:int}
	// /home/index/123可以匹配成功
	RouterPattern  string
	controllerName string
	patternRe      *regexp.Regexp
	//正则匹配出来的变量地址映射变量映射
	patternMap map[string]int
}

// Prepare 控制器执行前调用
func (c *Controller) Prepare() {
	//log.Println("父类的Prepare")
}

// Finish 控制器结束时调用
func (c *Controller) Finish() {
	//log.Println("父类的Finish")
}

// SkipController 停止执行控制器
func (c *Controller) SkipController() {
	c.isStopController = true
}

// setSonController 控制器注册
func (c *Controller) setSonController(son IController) {
	c.sonController = son
}

// initController http请求上下文注册
func (c *Controller) initController(ctx *freefishgo.HttpContext) {
	c.Response = ctx.Response
	c.Request = ctx.Request
	c.Data = map[interface{}]interface{}{}
}

// isNotSkip 过滤掉本地方法
func isNotSkip(methodName string) bool {
	skinList := map[string]bool{"SetHttpContext": true,
		"OverwriteActionRouter": true, "SetTplPath": true,
		"UseTplPath": true, "Prepare": true,
		"SkipController": true, "Finish": true}
	if _, ok := skinList[methodName]; ok {
		return false
	}
	return true
}
