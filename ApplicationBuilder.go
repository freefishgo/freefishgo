package freeFishGo

import "freeFishGo/httpContext"

type ApplicationBuilder struct {
	middlewareList []IMiddleware
	middlewareLink *MiddlewareLink
}

// 创建一个ApplicationBuilder管道
func NewApplicationBuilder() *ApplicationBuilder {
	return new(ApplicationBuilder)
}

// 中间件类型接口
type IMiddleware interface {
	Middleware(ctx *httpContext.HttpContext, next *MiddlewareLink) *httpContext.HttpContext
}
type MiddlewareLink struct {
	val  IMiddleware
	next *MiddlewareLink
}

// 执行下一个中间件
func (link *MiddlewareLink) Next(ctx *httpContext.HttpContext) {
	link.next.val.Middleware(ctx, link.next.next)
}

// 中间件注册接口
func (app *ApplicationBuilder) UseMiddleware(middleware IMiddleware) *ApplicationBuilder {
	if app.middlewareList == nil {
		app = NewApplicationBuilder()
	}
	app.middlewareList = append(app.middlewareList, middleware)
	return app
}

// 中间件排序
func (app *ApplicationBuilder) middlewareSorting() *ApplicationBuilder {
	app.middlewareLink = new(MiddlewareLink)
	tmpMid := app.middlewareLink
	for i := 0; i < len(app.middlewareList); i++ {
		tmpMid.val = app.middlewareList[i]
		tmpMid.next = new(MiddlewareLink)
		tmpMid = tmpMid.next
	}
	if tmpMid.next == nil {
		tmpMid.next = new(MiddlewareLink)
		tmpMid.next.val = &LastFrameMiddleware{}
	}
	return app
}

// 框架最后一个中间件
type LastFrameMiddleware struct {
}

func (last *LastFrameMiddleware) Middleware(ctx *httpContext.HttpContext, next *MiddlewareLink) *httpContext.HttpContext {
	return ctx
}
