# freeFishGo
golang 通过结构体反射实现的典型的mvc架构
```go
//继承router.Controller
type ctrTest struct {
	router.Controller
}
// 更改指定方法的路由规则，未更改的采用的是主路由规则     未设置主路由的采用默认的/{ Controller}/{Action}
func (c *ctrTest) GetControllerActionInfo() []*router.ControllerActionInfo {
	println("不是默认GetControllerInfo")
	tmp := make([]*router.ControllerActionInfo, 0)
	tmp = append(tmp, &router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}/{id:int}{string}/{int}", ControllerActionFuncName: "MyControllerActionStrut"})
	return tmp
}

type Test struct {
	T  []string `json:"tt"`
	T1 string   `json:"tstst1"`
	Id string   `json:"id"`
}
// 控制器动作的处理方法    参数Test是通过反射自动注入的, 其中路由规则中的{id:int}中的id也能映射到参数中
func (c *ctrTest) MyControllerActionStrut(Test *Test) {
	c.HttpContext.Response.Write([]byte(Test.Id))
}
func main() {
	app := NewFreeFish()
	// 注册控制器
	app.AddHanlers(&ctrTest{})
	// 注册主路由ControllerActionFuncName字段不用设置        设置了也不会生效
	app.AddMainRouter(&router.ControllerActionInfo{RouterPattern: "/{ Controller}/{Action}"})
	app.Run()
}
```
