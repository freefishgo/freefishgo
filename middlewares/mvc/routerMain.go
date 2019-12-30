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
	"bytes"
	"encoding/json"
	"errors"
	freeFishGo "github.com/freefishgo/freefishgo"
	"html/template"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
)

type controllerRegister struct {
	tree       *tree
	WebConfig  *MvcWebConfig
	staticFile map[string]template.HTML
}

// 实例化一个mvc注册器

func newControllerRegister() *controllerRegister {
	controllerRegister := new(controllerRegister)
	controllerRegister.tree = newTree()
	controllerRegister.staticFile = map[string]template.HTML{}
	return controllerRegister
}

func (cr *controllerRegister) AddHandlers(ctl IController) {
	ctl.setSonController(ctl)
	cr.tree = ctl.getControllerInfo(cr.tree)
}

// 主路由节点注册，必须含有{Controller}和{Action}变量
func (cr *controllerRegister) AddMainRouter(ctlList ...*MainRouter) {

	for _, v := range ctlList {
		if v.HomeController != "" && v.IndexAction != "" {
			v.IndexAction = replaceActionName(v.IndexAction)
			v.HomeController = strings.ToLower(v.HomeController)
			cr.tree.ControllerModelList.AddControllerModelList(&ControllerActionRouter{RouterPattern: "/", controllerName: v.HomeController, actionName: v.IndexAction})
		}
		cr.tree.MainRouterList = cr.tree.MainRouterList.AddControllerModelList(&ControllerActionRouter{RouterPattern: v.RouterPattern})
	}

}

// 如果主路由为空注册一个默认主路由
func (cr *controllerRegister) MainRouterNil() {
	if cr.tree.MainRouterList == nil || len(cr.tree.MainRouterList) == 0 {
		cr.tree.MainRouterList = cr.tree.MainRouterList.AddControllerModelList(&ControllerActionRouter{RouterPattern: "/{ Controller}/{Action}"})
	}
}

// http服务逻辑处理程序
//func (c *controllerRegister) Middleware(ctx *freeFishGo.HttpContext) {
//	c.AnalysisRequest(ctx)
//}
func (c *controllerRegister) AnalysisRequest(ctx *freeFishGo.HttpContext) (cont *freeFishGo.HttpContext) {
	cont = ctx
	defer func() {
		if err := recover(); err != nil {
			err, _ := err.(error)
			if c.WebConfig.RecoverPanic {
				c.WebConfig.RecoverFunc(ctx, err, debug.Stack())
			} else {
				if ctx != nil {
					ctx.Response.WriteHeader(500)
					ctx.Response.Write([]byte(`<html><body><div style="color: red;color: red;margin: 150px auto;width: 800px;"><div>` + "服务器内部错误 500:" + err.Error() + "\r\n\r\n\r\n</div><pre>" + string(debug.Stack()) + `</pre></div></body></html>`))
				}
			}
		}
	}()
	u, _ := url.Parse(ctx.Request.RequestURI)
	f := c.analysisUrlToGetAction(u, freeFishGo.HttpMethod(ctx.Request.Method))
	if f == nil {
		ctx.Response.WriteHeader(404)
		return ctx
	}
	ctl := f.ControllerInfo
	action := reflect.New(ctl.ControllerFunc)
	var ic IController = action.Interface().(IController)
	ic.initController(ctx)
	ic.Prepare()
	con := ic.getController()
	if con.isStopController {
		return ctx
	} else {
		ctx.Request.ParseForm()
		data := fromToSimpleMap(ctx.Request.Form, f.OtherKeyMap)
		ic.setQuery(data)
		if ctl.ControllerActionParameterStruct != nil {
			var param interface{}
			param = reflect.New(ctl.ControllerActionParameterStruct).Interface()
			dataString, err := json.Marshal(data)
			if err != nil {
				panic(err.Error())
			}
			json.Unmarshal(dataString, param)
			action.MethodByName(ctl.ControllerAction).Call(getValues(param))
		} else {
			action.MethodByName(ctl.ControllerAction).Call(nil)
		}
		if !ctx.Response.GetStarted() {
			con.controllerName = ctl.ControllerName
			con.actionName = ctl.ControllerAction
			err := c.tmpHtml(con)
			if err != nil {
				panic(err)
			}
		}
		ic.Finish()
	}
	return ctx
}
func (ctr *controllerRegister) tmpHtml(c *Controller) error {
	if c.isUseTplPath || c.LayoutPath != "" {
		if c.LayoutPath != "" {
			if b, err := ctr.htmlTpl(c.LayoutPath); err == nil {
				section := map[string]interface{}{}
				if c.tplPath == "" {
					c.tplPath = filepath.Join(c.controllerName, replaceActionName(c.actionName)+".fish")
				}
				if b, err := ctr.htmlTpl(c.tplPath); err == nil {
					section["LayoutContent"] = b
				} else {
					return errors.New("Controller:" + c.controllerName + "Action:" + c.actionName + "读取模板地址:" + c.tplPath + "时出错:" + err.Error())
				}
				for k, v := range c.LayoutSections {
					if b1, err := ctr.htmlTpl(v); err == nil {
						section[k] = b1
					} else {
						return errors.New("操作母模板:" + c.LayoutPath + " 读取子模板: " + k + " 子模板地址:" + v + "时出错:" + err.Error())
					}
				}
				var buf bytes.Buffer

				if t, err := template.New(c.tplPath).Delims(ctr.WebConfig.TemplateLeft, ctr.WebConfig.TemplateRight).Parse(string(b)); err == nil {
					if err := t.Execute(&buf, section); err == nil {
						if t, err := template.New(c.tplPath).Delims(ctr.WebConfig.TemplateLeft, ctr.WebConfig.TemplateRight).Parse(buf.String()); err == nil {
							return t.Execute(c.Response, c.Data)
						} else {
							return errors.New("Controller:" + c.controllerName + "Action:" + c.actionName + "读取模板地址:" + c.tplPath + "时出错:" + err.Error())
						}
					} else {
						return errors.New("格式化模板页时 母版页地址:" + c.LayoutPath + "时出错:" + err.Error())
					}
				} else {
					return errors.New("Controller:" + c.controllerName + "Action:" + c.actionName + "读取模板地址:" + c.tplPath + "时出错:" + err.Error())
				}
			} else {
				return errors.New("Controller:" + c.controllerName + "Action:" + c.actionName + "读取模板地址:" + c.tplPath + "时出错:" + err.Error())
			}
		}
		if c.tplPath == "" {
			c.tplPath = filepath.Join(c.controllerName, replaceActionName(c.actionName)+".fish")
		}
		if b, err := ctr.htmlTpl(c.tplPath); err == nil {
			// 创建一个新的模板，并且载入内容
			if t, err := template.New(c.tplPath).Delims(ctr.WebConfig.TemplateLeft, ctr.WebConfig.TemplateRight).Parse(string(b)); err == nil {
				return t.Execute(c.Response, c.Data)
			} else {
				return errors.New("Controller:" + c.controllerName + "Action:" + c.actionName + "读取模板地址:" + c.tplPath + "时出错:" + err.Error())
			}
		} else {
			return errors.New("Controller:" + c.controllerName + "Action:" + c.actionName + "读取模板地址:" + c.tplPath + "时出错:" + err.Error())
		}
	}
	return nil
}

func (ctr *controllerRegister) htmlTpl(path string) (template.HTML, error) {
	if v, ok := ctr.staticFile[path]; ok {
		return v, nil
	} else {
		temPath := filepath.Join(ctr.WebConfig.ViewsPath, path)
		if b, err := ioutil.ReadFile(temPath); err == nil {
			html := template.HTML(b)
			if ctr.WebConfig.IsDevelopment {
				return html, nil
			}
			ctr.staticFile[path] = html
			return html, nil
		} else {
			return "", err
		}
	}
}

//根据参数获取对应的Values
func getValues(param ...interface{}) []reflect.Value {
	vals := make([]reflect.Value, 0, len(param))
	for i := range param {
		vals = append(vals, reflect.ValueOf(param[i]))
	}
	return vals
}

func fromToSimpleMap(v url.Values, addKeyVal map[string]interface{}) map[string]interface{} {
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
	return dic
}

// 根据url对象分析出控制处理器名称，并把其他规则数据提取出来
func (c *controllerRegister) analysisUrlToGetAction(u *url.URL, method freeFishGo.HttpMethod) *freeFishUrl {
	path := strings.ToLower(u.Path)
	for _, v := range c.tree.MainRouterList {
		sl := v.patternRe.FindStringSubmatch(path)
		if len(sl) != 0 {
			ff := new(freeFishUrl)
			ff.OtherKeyMap = map[string]interface{}{}
			for k, m := range v.patternMap {
				ff.OtherKeyMap[k] = sl[m]
			}
			ff.controllerAction = ff.GetControllerAction(v) + strings.ToLower(string(method))
			ff.controllerName = ff.GetControllerName(v)
			if v, ok := c.tree.getControllerInfoByControllerNameControllerAction(ff.controllerName, ff.controllerAction); ok {
				if c.tree.CloseMainRouter[ff.controllerName][ff.controllerAction] {
					break
				}
				ff.ControllerInfo = v
				return ff
			}
		}
	}

	for _, v := range c.tree.ControllerModelList {
		sl := v.patternRe.FindStringSubmatch(path)
		if len(sl) != 0 {
			ff := new(freeFishUrl)
			ff.OtherKeyMap = map[string]interface{}{}
			for k, m := range v.patternMap {
				ff.OtherKeyMap[k] = sl[m]
			}
			ff.controllerAction = ff.GetControllerAction(v) + strings.ToLower(string(method))
			ff.controllerName = ff.GetControllerName(v)
			if v, ok := c.tree.getControllerInfoByControllerNameControllerAction(ff.controllerName, ff.controllerAction); ok {
				ff.ControllerInfo = v
				return ff
			}
		}
	}

	return nil
}
