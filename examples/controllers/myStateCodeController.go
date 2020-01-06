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

func (my *myStateCodeController) Error500() {
	my.UseTplPath()
}
func (my *myStateCodeController) Forbidden403() {
	my.StateCodeController.Forbidden403()
}
func (my *myStateCodeController) NotFind404() {
	my.StateCodeController.NotFind404()
}
