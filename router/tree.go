package router

import (
	"fmt"
	"log"
	"reflect"
)

type tree struct {
	StaticTree     map[string]*tree //静态路径
	RegularTree    map[string]*tree //正则路径    静态路径大于正则路径
	PrevTree       *tree            //上一个路由
	Path           string           //当前路径路由匹配规则
	AllPath        string           //当前路径的完整路径
	IsRoot         bool             //是否根路由
	Controller     Controller       //当前路径的处理程序
	ControllerFunc interface{}      //请求事件的处理函数
}

func newTree() *tree {
	return new(tree)
}

func reflectInterface(funcInter interface{}, paramsValue []reflect.Value) {
	v := reflect.ValueOf(funcInter)
	if v.Kind() != reflect.Func {
		log.Fatal("funcInter is not func")
	}
	values := v.Call(paramsValue) //方法调用并返回值
	for i := range values {
		fmt.Println(values[i])
	}
}
