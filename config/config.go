package config

import (
	"github.com/freeFishGo/httpContext"
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
	RecoverFunc       func(ctx *httpContext.HttpContext, e error, Stack []byte)
	SessionCookieName string
	EnableSession     bool
	Listen            Listen
}

type Listen struct {
	AutoTLS        bool
	ServerTimeOut  time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	ListenTCP4     bool
	EnableHTTP     bool
	HTTPAddr       string
	HTTPPort       int
	EnableHTTPS    bool
	HTTPSAddr      string
	HTTPSPort      int
	HTTPSCertFile  string
	HTTPSKeyFile   string
}
type WebConfig struct {
	TemplateLeft  string
	TemplateRight string
	ViewsPath     string
	RecoverPanic  bool
	RecoverFunc   func(ctx *httpContext.HttpContext, e error, Stack []byte)
}

func NewWebConfig() *WebConfig {
	return &WebConfig{ViewsPath: "views", TemplateLeft: "{{", TemplateRight: "}}"}
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
			ListenTCP4:    false,
			EnableHTTP:    true,
			AutoTLS:       false,
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
