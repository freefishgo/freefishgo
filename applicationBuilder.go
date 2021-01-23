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
package freefishgo

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

// DefaultApplicationBuilder is the default ApplicationBuilder used by Serve.
var defaultApplicationBuilder *ApplicationBuilder

func GetDefaultApplicationBuilder() *ApplicationBuilder {
	checkDefaultApplicationBuilderNil()
	return defaultApplicationBuilder
}

func SetDefaultApplicationBuilderConfig(config *Config) {
	checkDefaultApplicationBuilderNil()
	defaultApplicationBuilder.Config = config
}

func checkDefaultApplicationBuilderNil() {
	if defaultApplicationBuilder == nil {
		defaultApplicationBuilder = NewFreeFishApplicationBuilder()
	}
	if defaultApplicationBuilder.Config == nil {
		defaultApplicationBuilder.Config = NewConfig()
	}
}

// ApplicationBuilder管道构造器
type ApplicationBuilder struct {
	Config  *Config
	handler *applicationHandler
}

// 向管道注入session去数据的接口
func (app *ApplicationBuilder) InjectionSession(session ISession) {
	app.handler.session = session
}

// 向默认管道注入session去数据的接口
func InjectionSession(session ISession) {
	checkDefaultApplicationBuilderNil()
	defaultApplicationBuilder.InjectionSession(session)
}

// 创建一个ApplicationBuilder管道
func NewFreeFishApplicationBuilder() *ApplicationBuilder {
	freeFish := new(ApplicationBuilder)
	freeFish.handler = newApplicationHandler()
	freeFish.Config = NewConfig()
	return freeFish
}

// 启动默认中间件web服务
func Run() {
	checkDefaultApplicationBuilderNil()
	defaultApplicationBuilder.Run()
}

// 启动web服务
func (app *ApplicationBuilder) Run() {
	app.middlewareSorting()
	app.handler.config = app.Config
	errChan := make(chan error)
	if app.Config.EnableSession {
		if app.handler.session == nil {
			app.handler.session = NewSessionMgr(app.handler.config.SessionAliveTime)
		}
		app.handler.session.Init(app.handler.config.SessionAliveTime)
	}
	if app.Config.Listen.EnableHTTP {
		addr := app.Config.Listen.HTTPAddr + ":" + strconv.Itoa(app.Config.Listen.HTTPPort)
		go func() {
			log.Println("http on " + addr)
			errChan <- (&http.Server{
				Addr:           addr,
				ReadTimeout:    app.Config.Listen.ServerTimeOut,
				WriteTimeout:   app.Config.Listen.WriteTimeout,
				MaxHeaderBytes: app.Config.Listen.MaxHeaderBytes,
				Handler:        app.handler,
			}).ListenAndServe()
		}()
	}
	if app.Config.Listen.EnableHTTPS {
		addr := app.Config.Listen.HTTPSAddr + ":" + strconv.Itoa(app.Config.Listen.HTTPSPort)
		go func() {
			log.Println("https on " + addr)
			errChan <- (&http.Server{
				Addr:           addr,
				ReadTimeout:    app.Config.Listen.ServerTimeOut,
				WriteTimeout:   app.Config.Listen.WriteTimeout,
				MaxHeaderBytes: app.Config.Listen.MaxHeaderBytes,
				Handler:        app.handler,
			}).ListenAndServeTLS(app.Config.Listen.HTTPSCertFile, app.Config.Listen.HTTPSKeyFile)
		}()
	}
	for {
		select {
		case e := <-errChan:
			panic(e)

		}
	}
}

func newApplicationHandler() *applicationHandler {
	return new(applicationHandler)
}

type applicationHandler struct {
	middlewareList []IMiddleware
	middlewareLink *middlewareLink
	config         *Config
	session        ISession
}

// http服务逻辑处理程序
func (app *applicationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := new(HttpContext)
	ctx.setContext(rw, r)
	ctx.Response.getYourself().maxResponseCacheLen = app.config.MaxResponseCacheLen
	if app.config.EnableSession {
		ctx.Response.setISession(app.session)
		ctx.Response.getYourself().SessionCookieName = app.config.SessionCookieName
		ctx.Response.getYourself().SessionAliveTime = app.config.SessionAliveTime
		cookie, err := ctx.Request.Cookie(app.config.SessionCookieName)
		if err == nil {
			ctx.Response.getYourself().SessionId = cookie.Value
		}
	}
	defer func() {
		if ctx != nil && ctx.Response.getYourself().Gzip != nil {
			ctx.Response.getYourself().Gzip.Close()
		}
	}()
	defer func() {
		if app.config.EnableSession {
			ctx.Response.UpdateSession()
		}
		if err := recover(); err != nil {
			ctx.Response.getYourself().SetIsWriteInCache(false)
			ctx.Response.SetError(err)
			ctx.Response.SetStack(string(debug.Stack()))
			if app.config.RecoverPanic {
				app.config.RecoverFunc(ctx)
			} else {
				if ctx != nil {
					ctx.Response.WriteHeader(500)
					fmt.Fprintf(ctx.Response, `<html><body><div style="color: red;color: red;margin: 150px auto;width: 800px;"><div>500 Internal Server Error:  %s </div><pre>%s</pre></div></body></html>`, ctx.Response.Error(), ctx.Response.Stack())
				}
			}
		}
	}()
	ctx.Response.getYourself().IsOpenGzip = app.config.EnableGzip
	if ctx.Response.getYourself().IsOpenGzip {
		ctx.Response.Header().Set("Content-Encoding", "gzip")
	}
	ctx = app.middlewareLink.val.Middleware(ctx, app.middlewareLink.next.innerNext)
	ctx.Response.SetIsWriteInCache(false)
	ctx.Response.Write(nil)
}

// 下一个中间件
type Next func(*HttpContext) *HttpContext

// 中间件类型接口
type IMiddleware interface {
	// 中间件的逻辑处理函数 框架会调用
	Middleware(ctx *HttpContext, next Next) *HttpContext
	// 注册框架后 框架会自动调用这个函数
	LastInit(*Config)
}

type middlewareLink struct {
	val  IMiddleware
	next *middlewareLink
}

// 执行下一个中间件
func (link *middlewareLink) innerNext(ctx *HttpContext) (cont *HttpContext) {
	cont = ctx
	defer func() {
		if err := recover(); err != nil {
			ctx.Response.SetError(err)
			if ctx != nil {
				ctx.Response.WriteHeader(500)
			}
			ctx.Response.SetStack(string(debug.Stack()))
		}
	}()
	return link.val.Middleware(ctx, link.next.innerNext)
}

// 中间件注册接口
func (app *ApplicationBuilder) UseMiddleware(middleware ...IMiddleware) {
	if app.handler.middlewareList == nil {
		app.handler.middlewareList = []IMiddleware{}
	}
	app.handler.middlewareList = append(app.handler.middlewareList, middleware...)
}

// 中间件func注册接口
func (app *ApplicationBuilder) UseMiddlewareFunc(middlewareFunc ...func(ctx *HttpContext, next Next) *HttpContext) {
	if app.handler.middlewareList == nil {
		app.handler.middlewareList = []IMiddleware{}
	}
	for _, v := range middlewareFunc {
		mid := &innerMiddlewareFunc{
			f: v,
		}
		app.handler.middlewareList = append(app.handler.middlewareList, mid)
	}
}

type innerMiddlewareFunc struct {
	f func(ctx *HttpContext, next Next) *HttpContext
}

func (m *innerMiddlewareFunc) Middleware(ctx *HttpContext, next Next) *HttpContext {
	return m.f(ctx, next)
}

func (m *innerMiddlewareFunc) LastInit(config *Config) {
	//panic("implement me")
}

// 向默认中间件注册接口
func UseMiddleware(middleware ...IMiddleware) {
	checkDefaultApplicationBuilderNil()
	defaultApplicationBuilder.UseMiddleware(middleware...)
}

// 中间件func注册接口
func UseMiddlewareFunc(middlewareFunc ...func(ctx *HttpContext, next Next) *HttpContext) {
	checkDefaultApplicationBuilderNil()
	defaultApplicationBuilder.UseMiddlewareFunc(middlewareFunc...)
}

// 中间件排序
func (app *ApplicationBuilder) middlewareSorting() *ApplicationBuilder {
	app.handler.middlewareLink = new(middlewareLink)
	tmpMid := app.handler.middlewareLink
	for i := 0; i < len(app.handler.middlewareList); i++ {
		tmpMid.val = app.handler.middlewareList[i]
		tmpMid.val.LastInit(app.Config)
		tmpMid.next = new(middlewareLink)
		tmpMid = tmpMid.next
	}
	if tmpMid.val == nil {
		tmpMid.val = &lastFrameMiddleware{}
		tmpMid.val.LastInit(app.Config)
	}
	return app
}

// 框架最后一个中间件
type lastFrameMiddleware struct {
}

func (last *lastFrameMiddleware) Middleware(ctx *HttpContext, next Next) *HttpContext {
	return ctx
}
func (last *lastFrameMiddleware) LastInit(config *Config) {

}
