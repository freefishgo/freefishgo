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

package freefishgo

import (
	"time"
)

type Config struct {
	AppName    string
	RunMode    string
	ServerName string
	// EnableGzip 是否开启Gzip压缩
	EnableGzip bool
	// EnableSession 是否启用Session
	EnableSession bool
	// SessionAliveTime Session的存在时间
	SessionAliveTime time.Duration
	// MaxResponseCacheLen 响应数据最大缓存长度
	MaxResponseCacheLen int
	// SessionCookieName Session的客户端Cookie名字
	SessionCookieName string
	// RecoverPanic 是否在管道最末端捕获Panic，以取代框架的处理Panic函数
	RecoverPanic bool
	// RecoverFunc 捕获Panic的处理函数
	RecoverFunc func(ctx *HttpContext)
	Listen      Listen
}

type Listen struct {
	// ServerTimeOut 服务超时时间
	ServerTimeOut time.Duration
	// WriteTimeout 写超时时间
	WriteTimeout time.Duration

	MaxHeaderBytes int
	// EnableHTTP 是否开启http服务
	EnableHTTP bool
	// HTTPAddr http服务运行ip地址
	HTTPAddr string
	// HTTPPort http服务运行端口
	HTTPPort int
	// EnableHTTPS 是否开启Https服务
	EnableHTTPS bool
	// HTTPSAddr https服务运行ip地址
	HTTPSAddr string
	// HTTPSPort https服务运行端口
	HTTPSPort int
	// HTTPSCertFile httpsCertFile文件地址
	HTTPSCertFile string
	// HTTPSKeyFile HTTPSKeyFile文件地址
	HTTPSKeyFile string
}

const (
	VERSION = "1.00"
)

func NewConfig() *Config {
	return &Config{
		AppName:             "freeFishGo",
		ServerName:          "freeFishGoServer:" + VERSION,
		EnableGzip:          true,
		EnableSession:       true,
		SessionAliveTime:    time.Minute * 20,
		SessionCookieName:   "fishCookie",
		MaxResponseCacheLen: 2 << 12,
		Listen: Listen{
			ServerTimeOut: 0,
			EnableHTTP:    true,
			HTTPAddr:      "",
			HTTPPort:      8080,
			EnableHTTPS:   false,
			HTTPSAddr:     "",
			HTTPSPort:     8081,
			HTTPSCertFile: "server.pem",
			HTTPSKeyFile:  "server.key",
		},
	}
}
