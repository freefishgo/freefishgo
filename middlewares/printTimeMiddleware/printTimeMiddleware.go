package printTimeMiddleware

import (
	"github.com/freefishgo/freeFish"
	"log"
	"strconv"
	"time"
)

// 例子： 组装一个Middleware服务，实现打印mvc框架处理请求的时间
type PrintTimeMiddleware struct {
}

// 中间件打印mvc框架处理请求的时间
func (m *PrintTimeMiddleware) Middleware(ctx *freeFish.HttpContext, next freeFish.Next) *freeFish.HttpContext {
	dt := time.Now()
	ctxtmp := next(ctx)
	log.Println("路径:" + ctx.Request.URL.Path + "  处理时间为:" + (time.Now().Sub(dt)).String() + "  响应状态：" + strconv.Itoa(ctx.Response.ReadStatusCode()))
	return ctxtmp
}

// 中间件注册是调用函数进行该中间件最后的设置
func (*PrintTimeMiddleware) LastInit(*freeFish.Config) {
	//panic("implement me")
}
