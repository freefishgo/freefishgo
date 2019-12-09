package freeFishGo

import (
	"freeFishGo/config"
	"freeFishGo/router"
	"net/http"
	"strconv"
)

type app struct {
	handlers *router.ControllerRegister
	Server   *http.Server
	Config   *config.Config
}

func NewFreeFish() *app {
	freeFish := new(app)
	freeFish.handlers = router.NewControllerRegister()
	freeFish.Config = config.NewConfig()
	return freeFish
}

func (app *app) AddHanlers(ctrles ...router.IController) {
	for i := 0; i < len(ctrles); i++ {
		app.handlers.AddHandlers(ctrles[i])
	}
}

// 主节点路由匹配原则注册     目前系统变量支持格式为 `/{ Controller}/{Action}/{id:int}/{who:string}`
// 如果不进行路由注册  默认为/{ Controller}/{Action}   router.ControllerActionInfo中 ControllerActionFuncName不用设置  设置了也不会生效
func (app *app) AddMainRouter(list ...*router.ControllerActionInfo) {
	app.handlers.AddMainRouter(list...)
}

func (app *app) Run() {
	app.handlers.MainRouterNil()
	if app.Config.Listen.EnableHTTP {
		addr := app.Config.Listen.HTTPAddr + ":" + strconv.Itoa(app.Config.Listen.HTTPPort)
		app.Server = &http.Server{
			Addr: addr,
			//ReadTimeout:    app.Server.ReadTimeout,
			//WriteTimeout:   app.Server.WriteTimeout,
			//MaxHeaderBytes: app.Server.MaxHeaderBytes,
			Handler: app.handlers,
		}
		app.Server.ListenAndServe()
	}
}
