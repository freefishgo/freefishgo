package main

import (
	"freeFishGo/config"
	"freeFishGo/router"
	"net/http"
	"strconv"
)

type app struct {
	Handlers *router.ControllerRegister
	Server   *http.Server
	Config   *config.Config
}

func NewFreeFish() *app {
	freeFish := new(app)
	freeFish.Handlers = router.NewControllerRegister()
	freeFish.Config = config.NewConfig()
	return freeFish
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
				Handler: app.Handlers,
			}
			app.Server.ListenAndServe()
		}()
	}
}
