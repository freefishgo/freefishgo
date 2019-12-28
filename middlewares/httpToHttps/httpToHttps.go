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
package httpToHttps

import (
	freeFishGo "github.com/freefishgo/freefishgo"
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
