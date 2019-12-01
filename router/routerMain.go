package router

import (
	"encoding/json"
	"fmt"
	"freeFishGo/httpContext"
	"net/http"
	"net/url"
	"reflect"
)

type ControllerRegister struct {
	tree *tree
}

func NewControllerRegister() *ControllerRegister {
	controllerRegister := new(ControllerRegister)
	controllerRegister.tree = newTree()
	return controllerRegister
}

func (cr *ControllerRegister) AddHandlers(ctl IController) {
	ctl.setSonController(ctl)
	cr.tree = ctl.getControllerInfo(cr.tree)
}

// http服务逻辑处理程序
func (c *ControllerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	println(r.RequestURI)
	c.analysisRequest(rw, r)
}
func (c *ControllerRegister) analysisRequest(rw http.ResponseWriter, r *http.Request) (ctx *httpContext.HttpContext) {
	ctx = new(httpContext.HttpContext)
	ctx.Request = r
	ctx.Response = rw
	tmp := c.tree.ControllerList["ctrtest"]["mycontrolleractionstrut"].ControllerFunc
	t := reflect.New(tmp)
	tt := t.Interface()
	var ic IController = tt.(IController)
	ic.SetHttpContext(ctx)
	r.ParseForm()
	test := reflect.New(c.tree.ControllerList["ctrtest"]["mycontrolleractionstrut"].ControllerActionParameterStruct).Interface()
	fmt.Printf("数据：%+v", test)
	json.Unmarshal(fromToSimpleMap(r.Form), test)
	fmt.Printf("数据：%+v", test)
	t.MethodByName(c.tree.ControllerList["ctrtest"]["mycontrolleractionstrut"].ControllerAction).Call(getValues(test))
	return
}

//根据参数获取对应的Values
func getValues(param ...interface{}) []reflect.Value {
	vals := make([]reflect.Value, 0, len(param))
	for i := range param {
		vals = append(vals, reflect.ValueOf(param[i]))
	}
	return vals
}

func fromToSimpleMap(v url.Values) []byte {
	dic := map[string]interface{}{}
	for k, val := range v {
		if len(val) == 1 {
			dic[k] = val[0]
		} else {
			dic[k] = val
		}
	}
	data, err := json.Marshal(dic)
	if err != nil {
		panic(err.Error())
	}
	return data
}
