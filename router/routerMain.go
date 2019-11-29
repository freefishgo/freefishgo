package router

import (
	"net/http"
)

type ControllerRegister struct {
	tree *tree
}

func NewControllerRegister() *ControllerRegister {
	return new(ControllerRegister)
}

func (r *ControllerRegister) AddHandlers(ctl *IController) {
	(*ctl).setSonController(ctl)
	r.tree = (*ctl).getControllerInfo(r.tree)
}

// http服务逻辑处理程序
func (c *ControllerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	println(&c)
}
