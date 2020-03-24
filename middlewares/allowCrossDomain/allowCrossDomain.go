package allowCrossDomain

import (
	"github.com/freefishgo/freefishgo"
	"net/http"
)

// 组装一个Middleware服务，实现允许跨域请求
type AllowCrossDomain struct {
	// 允许跨域源
	AllowOrigin string
	// 允许跨域方法 格式 GET, POST, DELETE,PUT
	AllowMethods string
}

// 中间件实现允许跨域请求
func (allow *AllowCrossDomain) Middleware(ctx *freefishgo.HttpContext, next freefishgo.Next) *freefishgo.HttpContext {
	if http.MethodOptions == ctx.Request.Method {
		ctx.Response.Header().Set("Access-Control-Allow-Origin", allow.AllowOrigin)
		ctx.Response.Header().Set("Access-Control-Allow-Methods", allow.AllowMethods)
		return ctx
	}
	return next(ctx)
}

// 中间件注册时调用函数进行该中间件最后的设置
func (allow *AllowCrossDomain) LastInit(config *freefishgo.Config) {
	if allow.AllowOrigin == "" {
		allow.AllowOrigin = "*"
	}
	if allow.AllowMethods == "" {
		allow.AllowMethods = "GET, POST, DELETE,PUT"
	}
}
