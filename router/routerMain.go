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

// 主路由节点注册
func (cr *ControllerRegister) AddMainRouter(ctlList ...*ControllerActionInfo) {
	cr.tree.MainRouterList = cr.tree.MainRouterList.AddControllerModelList(ctlList...)
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
	controllerName := "ctrtest"
	controllerAction := "mycontrolleractionstrut"
	u, _ := url.Parse(ctx.Request.RequestURI)
	f := c.analysisUrlToGetAction(u)
	if f == nil {
		controllerName = "ctrtest"
		controllerAction = "mycontrolleractionstrut"
	} else {
		controllerName = f.GetControllerName()
		controllerAction = f.GetControllerAction()
	}
	ctl, ok := c.tree.getControllerInfoByControllerNameControllerAction(controllerName, controllerAction)
	if ok {
		action := reflect.New(ctl.ControllerFunc)
		var ic IController = action.Interface().(IController)
		ic.SetHttpContext(ctx)
		r.ParseForm()
		var param interface{}
		if ctl.ControllerActionParameterStruct != nil {
			param = reflect.New(ctl.ControllerActionParameterStruct).Interface()
			json.Unmarshal(fromToSimpleMap(r.Form, f.OtherKeyMap), param)
		}
		fmt.Printf("数据：%+v", param)
		action.MethodByName(ctl.ControllerAction).Call(getValues(param))
	} else {
		ctx.Response.Write([]byte("404错误"))
	}
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

func fromToSimpleMap(v url.Values, addKeyVal map[string]interface{}) []byte {
	dic := map[string]interface{}{}
	for k, val := range v {
		if len(val) == 1 {
			dic[k] = val[0]
		} else {
			dic[k] = val
		}
	}
	for k, val := range addKeyVal {
		dic[k] = val
	}
	data, err := json.Marshal(dic)
	if err != nil {
		panic(err.Error())
	}
	return data
}

// 根据url对象分析出控制处理器名称，并把其他规则数据提取出来
func (c *ControllerRegister) analysisUrlToGetAction(u *url.URL) (f *freeFishUrl) {
	for _, v := range c.tree.MainRouterList {
		sl := v.patternRe.FindStringSubmatch(u.Path)
		if len(sl) != 0 {
			f := new(freeFishUrl)
			f.OtherKeyMap = map[string]interface{}{}
			f.controllerAction = v.actionName
			f.controllerName = v.controllerName
			for k, m := range v.patternMap {
				f.OtherKeyMap[k] = sl[m]
			}
			break
		}
	}
	return f
}
