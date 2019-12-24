package controllers

import (
	"github.com/freeFishGo/examples/fishgo"
	"github.com/freeFishGo/middlewares/mvc"
	"io/ioutil"
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
	if b, err := ioutil.ReadFile(filepath.Join("static", d.Path)); err == nil {
		static.Response.Write(b)
	} else {
		static.Response.WriteHeader(404)
		static.Response.Write([]byte(err.Error()))
	}
}

func (static *staticController) SetInfo() []*mvc.ControllerActionInfo {
	tmp := make([]*mvc.ControllerActionInfo, 0)
	tmp = append(tmp, &mvc.ControllerActionInfo{RouterPattern: "static/{path:allString}", ControllerActionFuncName: "StaticFile"})
	return tmp
}
