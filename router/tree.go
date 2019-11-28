package router

import (
	"freeFishGo/httpContext"
	"net/http"
)

type tree struct {
	Tree   *tree
	Path   string
	Handle http.Handler
}
type Controller struct {
	HttpContext *httpContext.HttpContext
}

func (t *tree) AddHttpPath(ctl Controller) {

}
