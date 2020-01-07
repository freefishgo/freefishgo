package controllers

import (
	"github.com/freefishgo/freefishgo/middlewares/mvc"
)

func init() {
	mvc.SetStateCodeHandlers(&myStateCodeController{})
}

type myStateCodeController struct {
	mvc.StatusCodeController
}

// 500 错误处理函数
func (my *myStateCodeController) Error500() {
	my.StatusCodeController.Error500()
}

// 403 处理函数
func (my *myStateCodeController) Forbidden403() {
	my.StatusCodeController.Forbidden403()
}

// 404 处理函数
func (my *myStateCodeController) NotFind404() {
	my.StatusCodeController.NotFind404()
}
