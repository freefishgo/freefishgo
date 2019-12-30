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
	freeFishGo "github.com/freefishgo/freefishgo"
	"reflect"
	"regexp"
	"strings"
)

type controllerInfo struct {
	ControllerFunc                  reflect.Type //请求事件的处理函数
	ControllerName                  string       //控制器名称
	ControllerAction                string       //控制器处理方法
	ControllerActionParameterStruct reflect.Type
}

// http请求逻辑控制器
type Controller struct {
	//响应前端的处理 不建议使用
	Response freeFishGo.IResponse
	// 和前端一切的数据  都可以通过他获取
	Request        *freeFishGo.Request
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

// 使用模板的路径并启动调用模板
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

// 控制器的基本数据结构
type IController interface {
	getControllerInfo(*tree) *tree
	setSonController(IController)
	OverwriteRouter() []*ControllerActionRouter
	initController(ctx *freeFishGo.HttpContext)
	getController() *Controller
	setQuery(map[string]interface{})
	Prepare()
	Finish()
}

// 进行路由注册的基类 如果结构体含有Controller 则Controller去掉 如GetController 变位Get  忽略大小写
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
	controllerActionInfoList := (c.sonController).OverwriteRouter()
	if controllerActionInfoList == nil {
		return tree
	}
	for _, v := range controllerActionInfoList {
		_, ok := getType.MethodByName(v.ControllerActionFuncName)
		if !ok {
			panic(getType.String() + "方法" + v.ControllerActionFuncName + "不存在")
		}
		if isHaveHttpMethod(v.ControllerActionFuncName) {
			v.actionName = replaceActionName(v.ControllerActionFuncName)
			v.controllerName = strings.ToLower(controllerName)
			if tree.CloseMainRouter[v.controllerName] == nil {
				tree.CloseMainRouter[v.controllerName] = map[string]bool{}
				tree.CloseMainRouter[v.controllerName][strings.ToLower(v.ControllerActionFuncName)] = true
			} else {
				tree.CloseMainRouter[v.controllerName][strings.ToLower(v.ControllerActionFuncName)] = true
			}
		} else {
			v.actionName = replaceActionName(v.ControllerActionFuncName)
			v.controllerName = strings.ToLower(controllerName)
			if tree.CloseMainRouter[v.controllerName] == nil {
				tree.CloseMainRouter[v.controllerName] = map[string]bool{}
				tree.CloseMainRouter[v.controllerName][strings.ToLower(v.ControllerActionFuncName)+"get"] = true
			} else {
				tree.CloseMainRouter[v.controllerName][strings.ToLower(v.ControllerActionFuncName)+"get"] = true
			}
		}
		//f := regexp.MustCompile(`Controller$`)
		//controllerName = f.ReplaceAllString(controllerName, "")
		f := regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
		v.RouterPattern = f.ReplaceAllString(v.RouterPattern, strings.ToLower(controllerName))
		f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
		v.RouterPattern = f.ReplaceAllString(v.RouterPattern, v.actionName)
		tree.ControllerModelList = tree.ControllerModelList.AddControllerModelList(v)
	}
	return tree
}

func replaceActionName(actionName string) string {
	actionName = strings.ToUpper(actionName)
	httpMethodList := []freeFishGo.HttpMethod{freeFishGo.MethodPost,
		freeFishGo.MethodConnect, freeFishGo.MethodDelete,
		freeFishGo.MethodGet, freeFishGo.MethodHead, freeFishGo.MethodOptions,
		freeFishGo.MethodPatch, freeFishGo.MethodPut, freeFishGo.MethodTrace}
	for _, v := range httpMethodList {
		f := regexp.MustCompile(string(v) + "$")
		if f.MatchString(actionName) {
			return strings.ToLower(f.ReplaceAllString(actionName, ""))
		}
	}
	return strings.ToLower(actionName)

}

func isHaveHttpMethod(actionName string) bool {
	actionName = strings.ToUpper(actionName)
	httpMethodList := []freeFishGo.HttpMethod{freeFishGo.MethodPost,
		freeFishGo.MethodConnect, freeFishGo.MethodDelete,
		freeFishGo.MethodGet, freeFishGo.MethodHead, freeFishGo.MethodOptions,
		freeFishGo.MethodPatch, freeFishGo.MethodPut, freeFishGo.MethodTrace}
	for _, v := range httpMethodList {
		f := regexp.MustCompile(string(v) + "$")
		if f.MatchString(actionName) {
			return true
		}
	}
	return false

}

// 单一路由设置结构体
type ControllerActionRouter struct {
	// 传设置控制器的方法
	ControllerActionFuncName string
	//路由设置  如：/{Controller}/{Action}/{id:int}
	// /home/index/123可以匹配成功
	RouterPattern  string
	controllerName string
	actionName     string
	patternRe      *regexp.Regexp
	//正则匹配出来的变量地址映射变量映射
	patternMap map[string]int
}

// 控制器属性设置 路由变量路由中只能出现一次
func (c *Controller) OverwriteRouter() []*ControllerActionRouter {
	return nil
}

// 控制器执行前调用
func (c *Controller) Prepare() {
	//log.Println("父类的Prepare")
}

// 控制器结束时调用
func (c *Controller) Finish() {
	//log.Println("父类的Finish")
}

// 停止执行控制器
func (c *Controller) SkipController() {
	c.isStopController = true
}

// 控制器注册
func (c *Controller) setSonController(son IController) {
	c.sonController = son
}

// http请求上下文注册
func (c *Controller) initController(ctx *freeFishGo.HttpContext) {
	c.Response = ctx.Response
	c.Request = ctx.Request
	c.Data = map[interface{}]interface{}{}
}

// 过滤掉本地方法
func isNotSkip(methodName string) bool {
	skinList := map[string]bool{"SetHttpContext": true,
		"OverwriteRouter": true, "SetTplPath": true, "UseTplPath": true, "Prepare": true, "SkipController": true, "Finish": true}
	if _, ok := skinList[methodName]; ok {
		return false
	}
	return true
}
