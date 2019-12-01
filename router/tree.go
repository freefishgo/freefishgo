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

// 控制器注册
func (t *tree) addPathTree(controllerName string, controllerAction string, controllerFunc reflect.Type, ControllerActionParameterStruct reflect.Type) {
	controllerInfo := new(ControllerInfo)
	controllerInfo.ControllerAction = controllerAction
	controllerInfo.ControllerName = controllerName
	controllerInfo.ControllerFunc = controllerFunc
	controllerInfo.ControllerActionParameterStruct = ControllerActionParameterStruct
	if _, ok := t.ControllerList[controllerName]; !ok {
		t.ControllerList[strings.ToLower(controllerName)] = map[string]*ControllerInfo{}
	}

	t.ControllerList[strings.ToLower(controllerName)][strings.ToLower(controllerAction)] = controllerInfo
}

//  根据控制器名字和动作名字获取控制器
func (t *tree) getControllerInfoByControllerNameControllerAction(controllerName string, controllerAction string) (*ControllerInfo, bool) {
	if v, ok := t.ControllerList[controllerName]; ok {
		if v, ok := v[controllerAction]; ok {
			return v, true
		}
	}
	return nil, false
}
