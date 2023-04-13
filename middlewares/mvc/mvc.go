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
	"net/http"

	"github.com/freefishgo/freefishgo"
)

// 默认的MvcApp
var defaultMvcApp *MvcApp

type MvcApp struct {
	handlers *controllerRegister
	//Server   *http.Server
	Config *MvcWebConfig
}

func SetDefaultMvcWebConfig(webConfig *MvcWebConfig) {
	checkDefaultMvcApp()
	defaultMvcApp.Config = webConfig
	defaultMvcApp.handlers.WebConfig = webConfig
}
func GetDefaultMvcApp() *MvcApp {
	checkDefaultMvcApp()
	return defaultMvcApp
}

// Middleware Web服务逻辑处理程序
func (mvc *MvcApp) Middleware(ctx *freefishgo.HttpContext, next freefishgo.Next) (c *freefishgo.HttpContext) {
	c = ctx
	ctx = mvc.handlers.AnalysisRequest(ctx)
	return next(ctx)
}

// LastInit 框架注册完成时  进行最后的配置
func (mvc *MvcApp) LastInit(cnf *freefishgo.Config) {
	mvc.handlers.WebConfig = mvc.Config
	handle := http.FileServer(http.Dir(mvc.handlers.WebConfig.StaticDir))
	mvc.handlers.staticFileHandler = handle
	mvc.handlers.MainRouterNil()
	mvc.handlers.StateCodeNil()
}

// NewFreeFishMvcApp 实例化生成一个Mvc对象
func NewFreeFishMvcApp() *MvcApp {
	freeFish := new(MvcApp)
	freeFish.handlers = newControllerRegister()
	freeFish.Config = freeFish.handlers.WebConfig
	return freeFish
}

func checkDefaultMvcApp() {
	if defaultMvcApp == nil {
		defaultMvcApp = NewFreeFishMvcApp()
	}
	if defaultMvcApp.Config == nil {
		defaultMvcApp.Config = NewWebConfig()
	}
	defaultMvcApp.handlers.WebConfig = defaultMvcApp.Config
}

// AddHandlers 将Controller控制器注册到Mvc框架对象中 即是添加路由动作
func (mvc *MvcApp) AddHandlers(ic ...IController) {
	for i := 0; i < len(ic); i++ {
		mvc.handlers.AddHandlers(ic[i])
	}
}

// SetStatusCodeHandlers 将Controller控制器注册到Mvc框架的定制状态处理程序中 如：404状态自定义  不传使用默认的
func (mvc *MvcApp) SetStatusCodeHandlers(s IStatusCodeController) {
	mvc.handlers.SetStatusCodeHandlers(s)
}

// SetStatusCodeHandlers 将Controller控制器注册到Mvc框架的定制状态处理程序中 如：404状态自定义  不传使用默认的
func SetStatusCodeHandlers(s IStatusCodeController) {
	checkDefaultMvcApp()
	defaultMvcApp.SetStatusCodeHandlers(s)
}

// AddHandlers 将Controller控制器注册到默认的Mvc框架对象中 即是添加路由动作
func AddHandlers(ic ...IController) {
	checkDefaultMvcApp()
	defaultMvcApp.AddHandlers(ic...)
}

// AddMainRouter 主节点路由匹配原则注册     目前系统变量支持格式为 `/{ Controller}/{Action}/{id:int}/{who:string}/{allString}`
//
// 如果不进行路由注册  默认为/{ Controller}/{Action}   router.ControllerActionInfo中 ControllerActionFuncName不用设置  设置了也不会生效
func (mvc *MvcApp) AddMainRouter(list ...*MainRouter) {
	for _, v := range list {
		if mvc.Config.homeController == "" || mvc.Config.indexAction == "" && (v.HomeController != "" && v.IndexAction != "") {
			mvc.Config.homeController = v.HomeController
			mvc.Config.indexAction = v.IndexAction
			mvc.handlers.AddMainRouter(v)
		} else {
			v.IndexAction = ""
			v.HomeController = ""
			mvc.handlers.AddMainRouter(v)
		}
	}
}

// AddMainRouter 默认mvc框架 主节点路由匹配原则注册     目前系统变量支持格式为 `/{ Controller}/{Action}/{id:int}/{who:string}/{allString}`
//
// 如果不进行路由注册  默认为/{ Controller}/{Action}   router.ControllerActionInfo中 ControllerActionFuncName不用设置  设置了也不会生效
func AddMainRouter(list ...*MainRouter) {
	checkDefaultMvcApp()
	defaultMvcApp.AddMainRouter(list...)
}

// MainRouter 主路由
type MainRouter struct {
	//路由设置  如：/{Controller}/{Action}/{id:int}
	// /home/index/123可以匹配成功
	RouterPattern string
	// Controller名称
	HomeController string
	// 动作名称
	IndexAction string
}
