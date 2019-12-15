package middlewares

import (
	free "freeFishGo"
	"freeFishGo/httpContext"
	"log"
	"time"
)

// 例子： 组装一个Middleware服务，实现打印mvc框架处理请求的时间
type Mid struct {
}

// 中间件打印mvc框架处理请求的时间
func (m *Mid) Middleware(ctx *httpContext.HttpContext, next free.Next) *httpContext.HttpContext {
	dt := time.Now()
	log.Println(ctx.Request.URL)
	ctxtmp := next(ctx)
	log.Println("处理时间为:" + (time.Now().Sub(dt)).String())
	return ctxtmp
}

// 中间件注册是调用函数进行该中间件最后的设置
func (m *Mid) LastInit() {
	//panic("implement me")
}
