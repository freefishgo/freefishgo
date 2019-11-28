package router

import (
	"net/http"
)

type ControllerRegister struct {
	Tree  *tree
	Count int
}

func NewControllerRegister() *ControllerRegister {
	return new(ControllerRegister)
}

// http服务逻辑处理程序
func (c *ControllerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	println(&c)
	c.Count++
	println(c.Count)
}
