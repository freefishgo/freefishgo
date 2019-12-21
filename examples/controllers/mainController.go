package controllers

import (
	"github.com/freeFishGo/examples/fishgo"
	"github.com/freeFishGo/router"
)

// 实现mvc控制器的处理Main为控制器 {Controller}的值
type MainController struct {
	router.Controller
}

// 注册控制器
func init() {
	fishgo.Mvc.AddHandlers(&MainController{})
}

// SetInfo()特殊定制指定action的路由
func (c *MainController) SetInfo() []*router.ControllerActionInfo {
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "{ Controller}/{Action}/{tstst1:string}er", ControllerActionFuncName: "LayoutTestGet"})
	return tmp
}

// 作为 Action的请求参数的映射值
type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Main/MyControllerActionStrut 最后的单词为请求方式  该例子为Post请求
func (c *MainController) MyControllerActionStrut(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	// 调用模板引擎   默认模板地址为{ Controller}/{Action}.fish    即为Main/MyControllerActionStrut， c.UseTplPath()等效于c.UseTplPath("Main/MyControllerActionStrut")
	c.UseTplPath()
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Main/LayoutTestGet 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看

// 由于重写Controller的SetInfo方法
//
//&router.ControllerActionInfo{RouterPattern: "{string}/{ Controller}/{Action}/{tstst1:string}er", ControllerActionFuncName: "LayoutTestGet"}

//所以实际路由为:/任意字符串/main/layoutTest/任意字符串er
func (c *MainController) LayoutTestGet(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	c.LayoutSections = map[string]string{}
	c.LayoutSections["Scripts"] = "Other/Script.fish"
	c.LayoutSections["HtmlHead"] = "Other/HtmlHead.fish"
	c.LayoutPath = "layout.fish"
	c.UseTplPath("Other/layoutSon.fish")
}

// My{Action}的值 该方法的默认路由为/Main/My 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看 重定向使用方法
func (c *MainController) MyGET(Test *Test) {
	c.Response.Redirect("/haha/main/LayoutTestGet/fafafd4646er?id=我喜")
}

// My1为{Action}的值 该方法的默认路由为/Main/My1 get请求可以省略get后缀  查询具体字符串值可到httpContext包中查看
func (c *MainController) My1() {
	c.Response.WriteJson("Test")
}
