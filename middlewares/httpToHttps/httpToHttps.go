package httpToHttps

import (
	"github.com/freeFishGo/freeFish"
	"strconv"
	"strings"
)

type HttpToHttps struct {
	HTTPPort  string
	HTTPSPort string
}

func (h *HttpToHttps) Middleware(ctx *freeFish.HttpContext, next freeFish.Next) *freeFish.HttpContext {
	_host := strings.Split(ctx.Request.Host, ":")
	if _host[1] == h.HTTPPort {
		_host[1] = h.HTTPSPort
		target := "https://" + strings.Join(_host, ":") + ctx.Request.URL.Path
		ctx.Response.Redirect(target)
		return ctx
	}
	return next(ctx)
}

func (h *HttpToHttps) LastInit(c *freeFish.Config) {
	h.HTTPPort = strconv.Itoa(c.Listen.HTTPPort)
	h.HTTPSPort = strconv.Itoa(c.Listen.HTTPSPort)
}
