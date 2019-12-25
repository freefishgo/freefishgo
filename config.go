package freeFishGo

import (
	"time"
)

type Config struct {
	AppName           string
	RunMode           string
	ServerName        string
	EnableGzip        bool
	NeedGzipLen       int
	SessionAliveTime  time.Duration
	RecoverPanic      bool
	RecoverFunc       func(ctx *HttpContext, e error, Stack []byte)
	SessionCookieName string
	EnableSession     bool
	Listen            Listen
}

type Listen struct {
	ServerTimeOut  time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	EnableHTTP     bool
	HTTPAddr       string
	HTTPPort       int
	EnableHTTPS    bool
	HTTPSAddr      string
	HTTPSPort      int
	HTTPSCertFile  string
	HTTPSKeyFile   string
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
