package router

import (
	"reflect"
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

type tree struct {
	ControllerList map[string]map[string]*ControllerInfo //静态路径
}

func newTree() *tree {
	tree := new(tree)
	tree.ControllerList = map[string]map[string]*ControllerInfo{}
	return tree
}

func (t *tree) addPathTree(controllerName string, controllerAction string, controllerFunc reflect.Type) {
	controllerInfo := new(ControllerInfo)
	controllerInfo.ControllerAction = controllerAction
	controllerInfo.ControllerName = controllerName
	controllerInfo.ControllerFunc = controllerFunc
	if _, ok := t.ControllerList[controllerName]; !ok {
		t.ControllerList[strings.ToLower(controllerName)] = map[string]*ControllerInfo{}
	}

	t.ControllerList[strings.ToLower(controllerName)][strings.ToLower(controllerAction)] = controllerInfo
}
