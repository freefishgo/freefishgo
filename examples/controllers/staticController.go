package controllers

import (
	"github.com/freefishgo/freeFishGo/examples/fishgo"
	"github.com/freefishgo/freeFishGo/middlewares/mvc"
	"io"
	"os"
	"path/filepath"
)

type staticController struct {
	mvc.Controller
}

// 控制器注册
func init() {
	fishgo.Mvc.AddHandlers(&staticController{})
}

type data struct {
	Path string `json:"path"`
}

// 提供静态资源服务
func (static *staticController) StaticFile(d *data) {
	if f, err := os.Open(filepath.Join("static", d.Path)); err == nil {
		io.Copy(static.Response, f)
	} else {
		static.Response.WriteHeader(404)
		static.Response.Write([]byte(err.Error()))
	}
}

func (static *staticController) SetInfo() []*mvc.ControllerActionRouter {
	tmp := make([]*mvc.ControllerActionRouter, 0)
	tmp = append(tmp, &mvc.ControllerActionRouter{RouterPattern: "static/{path:allString}", ControllerActionFuncName: "StaticFile"})
	return tmp
}
