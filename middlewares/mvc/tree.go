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
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//type tree struct {
//	StaticTree       map[string]*tree //静态路径
//	RegularTree      map[string]*tree //正则路径    静态路径大于正则路径
//	PrevTree         *tree            //上一个路由
//	Path             string           //当前路径路由匹配规则
//	AllPath          string           //当前路径的完整路径
//	IsRoot           bool             //是否根路由
//	Controller       Controller       //当前路径的处理程序
//	ControllerFunc   reflect.Type    //请求事件的处理函数
//	controllerName   string           //控制器名称
//	controllerAction string           //控制器处理方法
//}
// 路由映射模型
type controllerModelList map[string]*ActionRouter

type tree struct {
	ControllerList map[string]map[string]*controllerInfo //静态路径
	//主要路由节点
	MainRouterList controllerModelList
	// 动作路由映射模型
	ActionRouterList controllerModelList
	// 控制器路由映射模型
	ControllerRouterList  controllerModelList
	CloseMainRouter       map[string]map[string]bool
	CloseControllerRouter map[string]bool
}

func (c controllerModelList) AddControllerModelList(list ...*ActionRouter) controllerModelList {
	if c == nil {
		c = controllerModelList{}
	}
	for _, v := range list {
		v.makePattern()
		if strings.Contains(v.patternRe.String(), "{") {
			panic("添加的路由存在冲突，该路由为" + v.RouterPattern + ". 错误的变量")
		}
		if _, ok := c[v.patternRe.String()]; ok {
			panic("添加的路由存在冲突，该路由为" + v.RouterPattern)
		} else {
			c[v.patternRe.String()] = v
		}
	}
	return c
}

// 计算路由信息
func (c *ActionRouter) makePattern() {
	pathPattern := c.RouterPattern
	if len(pathPattern) == 0 {
		panic("设置的路由匹配模式不能为空")
	}
	if pathPattern[0] != '/' {
		pathPattern = "/" + pathPattern
	}
	waitSortMap := map[string]string{}
	waitSortArr := make([]int, 0)
	sortMap := map[string]int{}
	f := regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
	t := f.FindAllStringIndex(pathPattern, -1)
	if len(t) > 1 {
		panic("路由注册时发现{Controller}使用超过1次")
	}
	if len(t) == 1 {
		waitSortMap[strconv.Itoa(t[0][0])] = "Controller"
		waitSortArr = append(waitSortArr, t[0][0])
	}
	f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
	t = f.FindAllStringIndex(pathPattern, -1)
	if len(t) > 1 {
		panic("路由注册时发现{Action}使用超过1次")
	}
	if len(t) == 1 {
		waitSortMap[strconv.Itoa(t[0][0])] = "Action"
		waitSortArr = append(waitSortArr, t[0][0])
	}
	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*(int|string|allString)+[\ ]*}`)
	t = f.FindAllStringIndex(pathPattern, -1)
	for _, v := range t {
		sl := strings.Trim(strings.Split(pathPattern[v[0]+1:v[1]], ":")[0], " ")
		for _, v := range waitSortMap {
			if v == sl {
				panic("路由注册时" + c.RouterPattern + "时发现{" + sl + "}使用超过1次")
			}
		}
		waitSortArr = append(waitSortArr, v[0])
		waitSortMap[strconv.Itoa(v[0])] = sl
	}
	sort.Ints(waitSortArr)
	for k, v := range waitSortArr {
		sortMap[waitSortMap[strconv.Itoa(v)]] = k + 1
	}
	f = regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `([\w+$]+)`)
	f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `([\w+$]+)`)
	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*int[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `(-?[1-9]\d+)`)
	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*string[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `([\w+$]+)`)

	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*allString[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `(.*?)`)

	f = regexp.MustCompile(`{[\ ]*int[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `-?[1-9]\d+`)
	f = regexp.MustCompile(`{[\ ]*string[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `[\w+$]+`)

	f = regexp.MustCompile(`{[\ ]*allString[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `.*?`)

	c.patternRe = regexp.MustCompile("^" + pathPattern + "$")
	c.patternMap = sortMap
}

func newTree() *tree {
	tree := new(tree)
	tree.ControllerList = map[string]map[string]*controllerInfo{}
	return tree
}

// 控制器注册
func (t *tree) addPathTree(controllerName string, controllerAction string, controllerFunc reflect.Type, ControllerActionParameterStruct reflect.Type) {
	controllerInfo1 := new(controllerInfo)
	controllerInfo1.ControllerAction = controllerAction
	controllerInfo1.ControllerName = controllerName
	controllerInfo1.ControllerFunc = controllerFunc
	controllerInfo1.ControllerActionParameterStruct = ControllerActionParameterStruct
	//re := regexp.MustCompile(`([\w+$]+)Controller$`)
	//tmpControllerNameList := re.FindStringSubmatch(controllerName)
	//if len(tmpControllerNameList) == 2 {
	//	controllerName = strings.ToLower(tmpControllerNameList[1])
	//	if _, ok := t.ControllerList[controllerName]; !ok {
	//		t.ControllerList[controllerName] = map[string]*controllerInfo{}
	//	}
	//} else {
	//	controllerName = strings.ToLower(controllerName)
	//	if _, ok := t.ControllerList[controllerName]; !ok {
	//		t.ControllerList[controllerName] = map[string]*controllerInfo{}
	//	}
	//}
	controllerName = strings.ToLower(controllerName)
	if _, ok := t.ControllerList[controllerName]; !ok {
		t.ControllerList[controllerName] = map[string]*controllerInfo{}
	}
	if !isHaveHttpMethod(controllerAction) {
		controllerAction += "get"
	}
	if _, ok := t.ControllerList[controllerName][strings.ToLower(controllerAction)]; ok {
		panic("进行控制器注入时发现{Controller}:" + controllerName + ",{Action}:" + strings.ToLower(controllerAction) + "重复注入。请注意{Controller}{Action}是不区分大小写的，且actionGet等效于action")
	}
	t.ControllerList[controllerName][strings.ToLower(controllerAction)] = controllerInfo1
}

//  根据控制器名字和动作名字获取控制器
func (t *tree) getControllerInfoByControllerNameControllerAction(controllerName string, controllerAction string) (*controllerInfo, bool) {
	if v, ok := t.ControllerList[controllerName]; ok {
		if v, ok := v[controllerAction]; ok {
			return v, true
		}
	}
	return nil, false
}
