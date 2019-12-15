package controllers

import (
	"fmt"
	"freeFishGo/examples/fishgo"
	"freeFishGo/router"
	"log"
)

// 实现mvc控制器的处理Main为控制器 {Controller}的值
type MainController struct {
	router.Controller
}

// 注册控制器
func init() {
	fishgo.AddHandlers(&MainController{})
}

// SetInfo()特殊定制指定action的路由
func (c *MainController) SetInfo() []*router.ControllerActionInfo {
	log.Println("不是默认GetControllerInfo")
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "{string}/{ Controller}/{Action}/{tstst1:string}er", ControllerActionFuncName: "MyControllerActionStrutPost"})
	return tmp
}

// 作为 Action的请求参数的映射值
type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Main/MyControllerActionStrut 最后的单词为请求方式  该例子为Post请求
func (c *MainController) MyControllerActionStrutPost(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	// 调用模板引擎   默认模板地址为{ Controller}/{Action}.fish    不含请求方式
	c.UseTplPath()
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Main/MyControllerActionStrut 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看
func (c *MainController) MyControllerActionStrutGet(Test *Test) {
	c.Data["Website"] = Test.Id
	c.Data["Email"] = Test.T1
	//c.HttpContext.Response.Write([]byte("hahaha"))
	c.UseTplPath()
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Main/My 最后的单词为请求方式该例子为Get请求  查询具体字符串值可到httpContext包中查看
func (c *MainController) MyGET(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}

// MyControllerActionStrut为{Action}的值 该方法的默认路由为/Main/My1 get请求可以省略get后缀  查询具体字符串值可到httpContext包中查看
func (c *MainController) My1(Test *Test) {
	c.HttpContext.Response.Write([]byte(fmt.Sprintf("数据为：%+v", Test)))
}
