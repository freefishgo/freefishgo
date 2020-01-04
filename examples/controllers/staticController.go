package controllers

import (
	"github.com/freefishgo/freefishgo/middlewares/mvc"
	"io"
	"os"
	"path/filepath"
)

type staticController struct {
	mvc.Controller
}

// 控制器注册
func init() {
	static := staticController{}
	//static.ActionRouterList = append(static.ActionRouterList,
	//	&mvc.ActionRouter{RouterPattern: "static/{path:allString}",
	//		ControllerActionFuncName: "StaticFile"})
	//static.ControllerRouter=&mvc.ControllerRouter{
	//	RouterPattern: "{Controller}/{Action}/{path:allString}",
	//}
	mvc.AddHandlers(&static)
}

type data struct {
	Path string `json:"path"`
}

// 提供静态资源服务
func (static *staticController) StaticFile(d *data) {
	if f, err := os.Open(filepath.Join("static", d.Path)); err == nil {
		//static.Response.Header().Set("Cache-Control","max-age=3600")
		switch filepath.Ext(d.Path) {
		case ".css":
			static.Response.Header().Set("Content-Type", "text/css")
			break
		case ".js":
			static.Response.Header().Set("Content-Type", "application/javascript")
			break
		}
		io.Copy(static.Response, f)
	} else {
		static.Response.WriteHeader(404)
		static.Response.Write([]byte(err.Error()))
	}
}
