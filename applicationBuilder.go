package freeFishGo

import (
	"github.com/freeFishGo/config"
	"github.com/freeFishGo/httpContext"
	"net/http"
	"runtime/debug"
	"strconv"
)

type ApplicationBuilder struct {
	server  *http.Server
	Config  *config.Config
	handler *ApplicationHandler
}

// 创建一个ApplicationBuilder管道
func NewFreeFishApplicationBuilder() *ApplicationBuilder {
	freeFish := new(ApplicationBuilder)
	freeFish.handler = newApplicationHandler()
	freeFish.Config = config.NewConfig()
	return freeFish
}
func (app *ApplicationBuilder) Run() {
	app.middlewareSorting()
	if app.Config.Listen.EnableHTTP {
		addr := app.Config.Listen.HTTPAddr + ":" + strconv.Itoa(app.Config.Listen.HTTPPort)
		app.handler.config = app.Config
		app.server = &http.Server{
			Addr: addr,
			//ReadTimeout:    MvcApp.Server.ReadTimeout,
			//WriteTimeout:   MvcApp.Server.WriteTimeout,
			//MaxHeaderBytes: MvcApp.Server.MaxHeaderBytes,
			Handler: app.handler,
		}
		if err := app.server.ListenAndServe(); err != nil {
			panic(err.Error())
		}

	}
}

func newApplicationHandler() *ApplicationHandler {
	return new(ApplicationHandler)
}

type ApplicationHandler struct {
	middlewareList []IMiddleware
	middlewareLink *MiddlewareLink
	config         *config.Config
}

// http服务逻辑处理程序
func (app *ApplicationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := new(httpContext.HttpContext)
	ctx.SetContext(rw, r)
	ctx.Response.SessionAliveTime = app.config.SessionAliveTime
	defer func() {
		if ctx != nil && ctx.Response.Gzip != nil {
			ctx.Response.Gzip.Close()
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			err, _ := err.(error)
			if app.config.RecoverPanic {
				app.config.RecoverFunc(ctx, err, debug.Stack())
			} else {
				if ctx != nil {
					ctx.Response.WriteHeader(500)
					ctx.Response.Write([]byte(`<html><body><div style="color: red;color: red;margin: 150px auto;width: 800px;"><div>` + "服务器内部错误 500:" + err.Error() + "\r\n\r\n\r\n</div><pre>" + string(debug.Stack()) + `</pre></div></body></html>`))
				}
			}
		}
	}()
	ctx.Response.IsOpenGzip = app.config.IsOpenGzip
	ctx.Response.NeedGzipLen = app.config.NeedGzipLen
	ctx = app.middlewareLink.val.Middleware(ctx, app.middlewareLink.next.innerNext)
	if !ctx.Response.Started {
		ctx.Response.ResponseWriter.WriteHeader(ctx.Response.ReadStatusCode())
	}
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
func (link *MiddlewareLink) innerNext(ctx *httpContext.HttpContext) *httpContext.HttpContext {
	return link.val.Middleware(ctx, link.next.innerNext)
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
		tmpMid.val = &lastFrameMiddleware{}
		tmpMid.val.LastInit()
	}
	return app
}

// 框架最后一个中间件
type lastFrameMiddleware struct {
}

func (last *lastFrameMiddleware) Middleware(ctx *httpContext.HttpContext, next Next) *httpContext.HttpContext {
	return ctx
}
func (last *lastFrameMiddleware) LastInit() {

}
