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
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"

	freeFishGo "github.com/freefishgo/freefishgo"
)

type controllerRegister struct {
	tree                    *tree
	WebConfig               *MvcWebConfig
	staticViewsFile         map[string]template.HTML
	staticFileHandler       http.Handler
	stateCodeControllerInfo *statusCodeControllerInfo
}

type statusCodeControllerInfo struct {
	statusCodeController reflect.Type
	name                 string
}

func (cr *controllerRegister) doStatusCode(ctx *freeFishGo.HttpContext) (ctxTmp *freeFishGo.HttpContext) {
	ctxTmp = ctx
	switch ctx.Response.ReadStatusCode() {
	case 404:
		if !ctx.Response.GetStarted() {
			ic := cr.initStateCodeFunc(ctx)
			if ic.getController().isStopController {
				return ctx
			}
			ic.NotFind404()
			cr.templateHtml(ic, "NotFind404")
			ic.Finish()
		}
		break
	case 403:
		if !ctx.Response.GetStarted() {
			ic := cr.initStateCodeFunc(ctx)
			if ic.getController().isStopController {
				return ctx
			}
			ic.Forbidden403()
			cr.templateHtml(ic, "Forbidden403")
			ic.Finish()
		}
		break
	case 500:
		if !ctx.Response.GetStarted() {
			ic := cr.initStateCodeFunc(ctx)
			if ic.getController().isStopController {
				return ctx
			}
			ic.Error500()
			cr.templateHtml(ic, "Error500")
			ic.Finish()
		}
		break
	}
	return ctx
}
func (cr *controllerRegister) initStateCodeFunc(ctx *freeFishGo.HttpContext) IStatusCodeController {

	if !ctx.Response.GetStarted() {
		ctx.Response.ClearWriteCache()
		stateCodeC := reflect.New(cr.stateCodeControllerInfo.statusCodeController)
		Is := stateCodeC.Interface().(IStatusCodeController)
		Is.initController(ctx)
		Is.Prepare()
		return Is
	}
	return nil
}

func (cr *controllerRegister) templateHtml(ic IStatusCodeController, MethodByName string) {
	con := ic.getController()
	con.controllerName = cr.stateCodeControllerInfo.name
	con.actionName = MethodByName
	err := cr.tmpHtml(con)
	if err != nil {
		panic(err)
	}
}

// 实例化一个mvc注册器

func newControllerRegister() *controllerRegister {
	controllerRegister := new(controllerRegister)
	controllerRegister.tree = newTree()
	controllerRegister.staticViewsFile = map[string]template.HTML{}
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
			cr.tree.ActionRouterList = cr.tree.ActionRouterList.AddControllerModelList(&ActionRouter{RouterPattern: "/", controllerName: v.HomeController, actionName: v.IndexAction})
		}
		cr.tree.MainRouterList = cr.tree.MainRouterList.AddControllerModelList(&ActionRouter{RouterPattern: v.RouterPattern})
	}

}

// 如果主路由为空注册一个默认主路由
func (cr *controllerRegister) MainRouterNil() {
	if cr.tree.MainRouterList == nil || len(cr.tree.MainRouterList) == 0 {
		cr.AddMainRouter(&MainRouter{RouterPattern: "/{ Controller}/{Action}"})
	}
}

// 如果StateCode处理为空，则采用默认的错误处理
func (cr *controllerRegister) StateCodeNil() {
	if cr.stateCodeControllerInfo == nil {
		cr.SetStatusCodeHandlers(&StatusCodeController{})
	}
}

// AddStateHandlers 将Controller控制器注册到Mvc框架的定制状态处理程序中 如：404状态自定义  不传使用默认的
func (handlers *controllerRegister) SetStatusCodeHandlers(s IStatusCodeController) {
	if handlers.stateCodeControllerInfo == nil {
		sInfo := new(statusCodeControllerInfo)
		sInfo.statusCodeController = reflect.TypeOf(s).Elem()
		controllerNameList := strings.Split(sInfo.statusCodeController.String(), ".")
		controllerName := controllerNameList[len(controllerNameList)-1]
		f := regexp.MustCompile(`Controller$`)
		controllerName = f.ReplaceAllString(controllerName, "")
		sInfo.name = controllerName
		handlers.stateCodeControllerInfo = sInfo
	} else {
		panic("StateCode注册处理重复，请检查,重复为：" + reflect.TypeOf(s).Name())
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
			ctx.Response.SetError(err)
			if ctx != nil {
				ctx.Response.WriteHeader(500)
			}
			ctx.Response.SetStack(string(debug.Stack()))
		}
		c.doStatusCode(ctx)
	}()
	u, _ := url.Parse(ctx.Request.RequestURI)
	f := c.analysisUrlToGetAction(u, freeFishGo.HttpMethod(ctx.Request.Method))
	if f == nil {
		if ctx.Response.GetIsWriteInCache() {
			c.staticFileHandler.ServeHTTP(ctx.Response, ctx.Request.Request)
		} else {
			ctx.Response.SetIsWriteInCache(true)
			ctx.Response.SetMaxResponseCacheLen(1 << 9)
			c.staticFileHandler.ServeHTTP(ctx.Response, ctx.Request.Request)
		}
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
			MapStringToStruct(param, data)
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

				if t, err := template.New(c.tplPath).Delims(ctr.WebConfig.LayoutTemplateLeft, ctr.WebConfig.LayoutTemplateRight).Parse(string(b)); err == nil {
					if err := t.Execute(&buf, section); err == nil {
						if c.Data == nil {
							c.Response.Write(buf.Bytes())
							return nil
						}
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
				return errors.New("Controller:" + c.controllerName + "Action:" + c.actionName + "读取模板地址:" + c.LayoutPath + "时出错:" + err.Error())
			}
		}
		if c.tplPath == "" {
			c.tplPath = filepath.Join(c.controllerName, replaceActionName(c.actionName)+".fish")
		}
		if b, err := ctr.htmlTpl(c.tplPath); err == nil {
			if c.Data == nil {
				c.Response.Write([]byte(b))
				return nil
			}
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
	if v, ok := ctr.staticViewsFile[path]; ok {
		return v, nil
	} else {
		temPath := filepath.Join(ctr.WebConfig.ViewsPath, path)
		if b, err := ioutil.ReadFile(temPath); err == nil {
			html := template.HTML(b)
			if ctr.WebConfig.IsDevelopment {
				return html, nil
			}
			ctr.staticViewsFile[path] = html
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
		if k == "Action" || k == "Controller" {
			continue
		}
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
				if c.tree.CloseControllerRouter[ff.controllerName] {
					continue
				}
				if c.tree.CloseMainRouter[ff.controllerName][ff.controllerAction] {
					continue
				}
				ff.ControllerInfo = v
				return ff
			}
		}
	}

	for _, v := range c.tree.ControllerRouterList {
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
					continue
				}
				ff.ControllerInfo = v
				return ff
			}
		}
	}

	for _, v := range c.tree.ActionRouterList {
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
