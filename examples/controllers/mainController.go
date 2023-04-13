package controllers

import (
	"log"

	"github.com/freefishgo/freefishgo/middlewares/mvc"
)

// HomeController 实现mvc控制器的处理Main为控制器 {Controller}的值
type HomeController struct {
	mvc.Controller
}

// 注册控制器
func init() {
	mvc.AddHandlers(&HomeController{})
}

func (c *HomeController) Prepare() {
	log.Println("子类的Prepare")
}

// Finish 控制器结束时调用
func (c *HomeController) Finish() {
	log.Println("子类的Finish")
}

// 作为 Action的请求参数的映射值
type Test struct {
	T  []*string `json:"tt"`
	T1 string    `json:"tstst1"`
	Id bool      `json:"id"`
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Home/MyControllerActionStrut 最后的单词为请求方式  该例子为Post请求
func (c *HomeController) MyControllerActionStrut(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	// 调用模板引擎   默认模板地址为{ Controller}/{Action}.fish    即为Main/MyControllerActionStrut， c.UseTplPath()等效于c.UseTplPath("Home/MyControllerActionStrut")
	c.UseTplPath()
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Home/LayoutTestGet 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看

// 由于重写Controller的SetInfo方法
//
//&router.ControllerActionRouter{RouterPattern: "{string}/{ Controller}/{Action}/{tstst1:string}er", ControllerActionFuncName: "LayoutTestGet"}

// 所以实际路由为:/任意字符串/main/layoutTest/任意字符串er
func (c *HomeController) LayoutTestGet(Test *Test) {
	log.Println("控制器执行成功")
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	c.LayoutSections = map[string]string{}
	c.LayoutSections["Scripts"] = "Other/Script.fish"
	c.LayoutSections["HtmlHead"] = "Other/HtmlHead.fish"
	c.LayoutPath = "layout.fish"
	c.UseTplPath("Other/layoutSon.fish")
}

// My{Action}的值 该方法的默认路由为/Home/My 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看 重定向使用方法
func (c *HomeController) MyGET(Test *Test) {
	c.Response.Redirect("/haha/main/LayoutTestGet/fafafd4646er?id=我喜")
}

// My1为{Action}的值 该方法的默认路由为/Home/My1 get请求可以省略get后缀  查询具体字符串值可到httpContext包中查看
func (c *HomeController) My1() {
	c.UseTplPath()
}

func (c *HomeController) My2() {
	if conn, err := c.Response.WebSocket(); err == nil {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
			//conn.Close()
			break
		}
	} else {
		log.Println(err.Error())
	}
}
