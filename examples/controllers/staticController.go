package controllers

import (
	"github.com/freeFishGo/examples/fishgo"
	"github.com/freeFishGo/router"
	"io/ioutil"
	"path/filepath"
)

type staticController struct {
	router.Controller
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
		static.HttpContext.Response.Write(b)
	} else {
		static.HttpContext.Response.WriteHeader(404)
		static.HttpContext.Response.Write([]byte(err.Error()))
	}
}

func (static *staticController) SetInfo() []*router.ControllerActionInfo {
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "static/{path:allString}", ControllerActionFuncName: "StaticFile"})
	return tmp
}
