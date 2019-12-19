package config

import (
	"github.com/freeFishGo/httpContext"
)

type Config struct {
	AppName             string
	RunMode             string
	RouterCaseSensitive bool
	ServerName          string
	CopyRequestBody     bool
	EnableGzip          bool
	NeedGzipLen         int

	RecoverPanic bool
	RecoverFunc  func(ctx *httpContext.HttpContext, e error, Stack []byte)

	MaxMemory          int64
	EnableErrorsShow   bool
	EnableErrorsRender bool
	Listen             Listen
	Log                LogConfig

	IsOpenGzip bool
}

type Listen struct {
	AutoTLS       bool
	ServerTimeOut int64
	ListenTCP4    bool
	EnableHTTP    bool
	HTTPAddr      string
	HTTPPort      int
	EnableHTTPS   bool
	HTTPSAddr     string
	HTTPSPort     int
	HTTPSCertFile string
	HTTPSKeyFile  string
}
type WebConfig struct {
	TemplateLeft  string
	TemplateRight string
	ViewsPath     string
	RecoverPanic  bool
	RecoverFunc   func(ctx *httpContext.HttpContext, e error, Stack []byte)
	Session       SessionConfig
}

func NewWebConfig() *WebConfig {
	return &WebConfig{ViewsPath: "views", TemplateLeft: "{{", TemplateRight: "}}"}
}

type SessionConfig struct {
	SessionOn                    bool
	SessionProvider              string
	SessionName                  string
	SessionGCMaxLifetime         int64
	SessionProviderConfig        string
	SessionCookieLifeTime        int
	SessionAutoSetCookie         bool
	SessionDomain                string
	SessionDisableHTTPOnly       bool
	SessionEnableSidInHTTPHeader bool
	SessionNameInHTTPHeader      string
	SessionEnableSidInURLQuery   bool
}

type LogConfig struct {
}

const (
	VERSION = "1.00"
)

func NewConfig() *Config {
	return &Config{
		AppName:             "freeFishGo",
		RouterCaseSensitive: true,
		ServerName:          "freeFishGoServer:" + VERSION,
		CopyRequestBody:     false,
		EnableGzip:          false,
		MaxMemory:           1 << 26, //64MB
		EnableErrorsShow:    true,
		IsOpenGzip:          true,
		NeedGzipLen:         1 << 11,
		EnableErrorsRender:  true,
		Listen: Listen{
			ServerTimeOut: 0,
			ListenTCP4:    false,
			EnableHTTP:    true,
			AutoTLS:       false,
			HTTPAddr:      "",
			HTTPPort:      8080,
			EnableHTTPS:   false,
			HTTPSAddr:     "",
			HTTPSPort:     10443,
			HTTPSCertFile: "",
			HTTPSKeyFile:  "",
		},
		Log: LogConfig{},
	}
}
