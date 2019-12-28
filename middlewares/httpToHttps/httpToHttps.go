package httpToHttps

import (
	"github.com/freefishgo/freefishgo"
	"strconv"
	"strings"
)

type HttpToHttps struct {
	HTTPPort  string
	HTTPSPort string
}

func (h *HttpToHttps) Middleware(ctx *freeFishGo.HttpContext, next freeFishGo.Next) *freeFishGo.HttpContext {
	_host := strings.Split(ctx.Request.Host, ":")
	if _host[1] == h.HTTPPort {
		_host[1] = h.HTTPSPort
		target := "https://" + strings.Join(_host, ":") + ctx.Request.URL.Path
		if len(ctx.Request.URL.RawQuery) > 0 {
			target += "?" + ctx.Request.URL.RawQuery
		}
		ctx.Response.Redirect(target)
		return ctx
	}
	return next(ctx)
}

func (h *HttpToHttps) LastInit(c *freeFishGo.Config) {
	h.HTTPPort = strconv.Itoa(c.Listen.HTTPPort)
	h.HTTPSPort = strconv.Itoa(c.Listen.HTTPSPort)
}
