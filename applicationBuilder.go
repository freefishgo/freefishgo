package freeFishGo

import (
	"freeFishGo/config"
	"freeFishGo/httpContext"
	"log"
	"net/http"
	"strconv"
)

type ApplicationBuilder struct {
	Server  *http.Server
	Config  *config.Config
	handler *ApplicationHandler
}

func NewFreeFishApplicationBuilder() *ApplicationBuilder {
	freeFish := new(ApplicationBuilder)
	freeFish.handler = NewApplicationHandler()
	freeFish.Config = config.NewConfig()
	return freeFish
}
func (app *ApplicationBuilder) Run() {
	app.middlewareSorting()
	if app.Config.Listen.EnableHTTP {
		addr := app.Config.Listen.HTTPAddr + ":" + strconv.Itoa(app.Config.Listen.HTTPPort)
		app.Server = &http.Server{
			Addr: addr,
			//ReadTimeout:    MvcApp.Server.ReadTimeout,
			//WriteTimeout:   MvcApp.Server.WriteTimeout,
			//MaxHeaderBytes: MvcApp.Server.MaxHeaderBytes,
			Handler: app.handler,
		}
		if err := app.Server.ListenAndServe(); err != nil {
			panic(err.Error())
		}

	}
}

func NewApplicationHandler() *ApplicationHandler {
	return new(ApplicationHandler)
}

type ApplicationHandler struct {
	middlewareList []IMiddleware
	middlewareLink *MiddlewareLink
}

// http服务逻辑处理程序
func (app *ApplicationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := new(httpContext.HttpContext)
	ctx.SetContext(rw, r)
	app.middlewareLink.val.Middleware(ctx, app.middlewareLink.next.Next)

}

// 创建一个ApplicationBuilder管道
func NewApplicationBuilder() *ApplicationBuilder {
	return new(ApplicationBuilder)
}

type Next func(*httpContext.HttpContext) *httpContext.HttpContext

// 中间件类型接口
type IMiddleware interface {
	Middleware(ctx *httpContext.HttpContext, next Next) *httpContext.HttpContext
	//注册框架后 框架会自动调用这个函数
	LastInit()
}
type MiddlewareLink struct {
	val  IMiddleware
	next *MiddlewareLink
}

// 执行下一个中间件
func (link *MiddlewareLink) Next(ctx *httpContext.HttpContext) *httpContext.HttpContext {
	return link.val.Middleware(ctx, link.next.Next)
}

// 中间件注册接口
func (app *ApplicationBuilder) UseMiddleware(middleware IMiddleware) *ApplicationBuilder {
	if app.handler.middlewareList == nil {
		app.handler.middlewareList = []IMiddleware{}
	}
	app.handler.middlewareList = append(app.handler.middlewareList, middleware)
	return app
}

// 中间件排序
func (app *ApplicationBuilder) middlewareSorting() *ApplicationBuilder {
	app.handler.middlewareLink = new(MiddlewareLink)
	tmpMid := app.handler.middlewareLink
	for i := 0; i < len(app.handler.middlewareList); i++ {
		tmpMid.val = app.handler.middlewareList[i]
		tmpMid.val.LastInit()
		tmpMid.next = new(MiddlewareLink)
		tmpMid = tmpMid.next
	}
	if tmpMid.val == nil {
		tmpMid.val = &LastFrameMiddleware{}
		tmpMid.val.LastInit()
	}
	return app
}

// 框架最后一个中间件
type LastFrameMiddleware struct {
}

func (last *LastFrameMiddleware) Middleware(ctx *httpContext.HttpContext, next Next) *httpContext.HttpContext {
	return ctx
}
func (last *LastFrameMiddleware) LastInit() {
	log.Println("管道最后一层设置完成")
}
