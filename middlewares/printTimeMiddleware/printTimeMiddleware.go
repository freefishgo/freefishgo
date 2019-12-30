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
package printTimeMiddleware

import (
	freeFishGo "github.com/freefishgo/freefishgo"
	"log"
	"strconv"
	"time"
)

// 例子： 组装一个Middleware服务，实现打印mvc框架处理请求的时间
type PrintTimeMiddleware struct {
}

// 中间件打印mvc框架处理请求的时间
func (m *PrintTimeMiddleware) Middleware(ctx *freeFishGo.HttpContext, next freeFishGo.Next) *freeFishGo.HttpContext {
	dt := time.Now()
	ctx.Response.IsWriteInCache = true
	ctxtmp := next(ctx)
	log.Println("路径:" + ctx.Request.URL.Path + "  处理时间为:" + (time.Now().Sub(dt)).String() + "  响应状态：" + strconv.Itoa(ctx.Response.ReadStatusCode()))
	return ctxtmp
}

// 中间件注册是调用函数进行该中间件最后的设置
func (*PrintTimeMiddleware) LastInit(*freeFishGo.Config) {
	//panic("implement me")
}
