package main

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
		app.handlers.AddHandlers(&ctrles[i])
	}
}

func (app *app) Run() {
	if app.Config.Listen.EnableHTTP {
		go func() {
			addr := app.Config.Listen.HTTPAddr + ":" + strconv.Itoa(app.Config.Listen.HTTPPort)
			app.Server = &http.Server{
				Addr: addr,
				//ReadTimeout:    app.Server.ReadTimeout,
				//WriteTimeout:   app.Server.WriteTimeout,
				//MaxHeaderBytes: app.Server.MaxHeaderBytes,
				Handler: app.handlers,
			}
			app.Server.ListenAndServe()
		}()
	}
}
