package middlewares

import (
	"github.com/freeFishGo"
	"github.com/freeFishGo/config"
	"github.com/freeFishGo/httpContext"
	"github.com/freeFishGo/router"
	"log"
	"os"
	"path/filepath"
)

type MvcApp struct {
	handlers *router.ControllerRegister
	//Server   *http.Server
	//Config   *config.Config
}

// http服务逻辑处理程序
func (mvc *MvcApp) Middleware(ctx *httpContext.HttpContext, next freeFishGo.Next) (c *httpContext.HttpContext) {
	c = ctx
	ctx = mvc.handlers.AnalysisRequest(ctx, mvc.handlers.WebConfig)
	return next(ctx)
}
func (mvc *MvcApp) LastInit(cnf *config.Config) {
	mvc.handlers.MainRouterNil()
	log.Println("MVC注册成功并完成LastInit初始化")
}

func NewFreeFishMvcApp() *MvcApp {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(dir)
	freeFish := new(MvcApp)
	freeFish.handlers = router.NewControllerRegister()
	return freeFish
}

func (app *MvcApp) AddHandlers(ic ...router.IController) {
	for i := 0; i < len(ic); i++ {
		app.handlers.AddHandlers(ic[i])
	}
}

// 主节点路由匹配原则注册     目前系统变量支持格式为 `/{ Controller}/{Action}/{id:int}/{who:string}/{allString}`
//
// 如果不进行路由注册  默认为/{ Controller}/{Action}   router.ControllerActionInfo中 ControllerActionFuncName不用设置  设置了也不会生效
func (app *MvcApp) AddMainRouter(list ...*router.ControllerActionInfo) {
	app.handlers.AddMainRouter(list...)
}
