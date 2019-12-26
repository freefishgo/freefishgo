package freeFishGo

import (
	"time"
)

type Config struct {
	AppName    string
	RunMode    string
	ServerName string
	// 是否开启Gzip压缩
	EnableGzip bool
	// 一次写入长度超过多少需要Gzip压缩
	NeedGzipLen int
	// 是否启用Session
	EnableSession bool
	// Session的存在时间
	SessionAliveTime time.Duration
	// Session的客户端Cookie名字
	SessionCookieName string
	// 是否在管道最末端捕获Panic，以取代框架的处理Panic函数
	RecoverPanic bool
	// 捕获Panic的处理函数
	RecoverFunc func(ctx *HttpContext, e error, Stack []byte)
	Listen      Listen
}

type Listen struct {
	// 服务超时时间
	ServerTimeOut time.Duration
	//写超时时间
	WriteTimeout time.Duration

	MaxHeaderBytes int
	// 是否开启http服务
	EnableHTTP bool
	// http服务运行ip地址
	HTTPAddr string
	// http服务运行端口
	HTTPPort int
	// 是否开启Https服务
	EnableHTTPS bool
	// https服务运行ip地址
	HTTPSAddr string
	// https服务运行端口
	HTTPSPort int
	// httpsCertFile文件地址
	HTTPSCertFile string
	// HTTPSKeyFile文件地址
	HTTPSKeyFile string
}

const (
	VERSION = "1.00"
)

func NewConfig() *Config {
	return &Config{
		AppName:           "freeFishGo",
		ServerName:        "freeFishGoServer:" + VERSION,
		EnableGzip:        true,
		NeedGzipLen:       1 << 11,
		EnableSession:     true,
		SessionAliveTime:  time.Minute * 20,
		SessionCookieName: "fishCookie",
		Listen: Listen{
			ServerTimeOut: 0,
			EnableHTTP:    true,
			HTTPAddr:      "",
			HTTPPort:      8080,
			EnableHTTPS:   true,
			HTTPSAddr:     "",
			HTTPSPort:     8081,
			HTTPSCertFile: "server.pem",
			HTTPSKeyFile:  "server.key",
		},
	}
}
