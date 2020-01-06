package controllers

import (
	"github.com/freefishgo/freefishgo/middlewares/mvc"
)

func init() {
	mvc.SetStateCodeHandlers(&myStateCodeController{})
}

type myStateCodeController struct {
	mvc.StateCodeController
}

// 500 错误处理函数
func (my *myStateCodeController) Error500() {
	my.StateCodeController.Error500()
}

// 403 处理函数
func (my *myStateCodeController) Forbidden403() {
	my.StateCodeController.Forbidden403()
}

// 404 处理函数
func (my *myStateCodeController) NotFind404() {
	my.StateCodeController.NotFind404()
}
