package router

import (
	"fmt"
	"freeFishGo/httpContext"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//type tree struct {
//	StaticTree       map[string]*tree //静态路径
//	RegularTree      map[string]*tree //正则路径    静态路径大于正则路径
//	PrevTree         *tree            //上一个路由
//	Path             string           //当前路径路由匹配规则
//	AllPath          string           //当前路径的完整路径
//	IsRoot           bool             //是否根路由
//	Controller       Controller       //当前路径的处理程序
//	ControllerFunc   reflect.Type    //请求事件的处理函数
//	controllerName   string           //控制器名称
//	controllerAction string           //控制器处理方法
//}
// 路由映射模型
type ControllerModelList map[string]*ControllerActionInfo

type tree struct {
	ControllerList map[string]map[string]*ControllerInfo //静态路径
	//主要路由节点
	MainRouterList ControllerModelList
	// 路由映射模型
	ControllerModelList ControllerModelList
}

func (c ControllerModelList) AddControllerModelList(list ...*ControllerActionInfo) ControllerModelList {
	if c == nil {
		c = ControllerModelList{}
	}
	for _, v := range list {
		v.makePattern()
		if strings.Contains(v.patternRe.String(), "{") {
			panic("添加的路由存在冲突，该路由为" + v.RouterPattern + ". 错误的变量")
		}
		if len(v.AllowMethod) == 0 {
			v.AllowMethod = append(v.AllowMethod, httpContext.MethodGet)
		}
		v.allowMethod = map[httpContext.HttpMethod]bool{}
		for _, cc := range v.AllowMethod {
			v.allowMethod[cc] = true
		}
		if _, ok := c[v.patternRe.String()]; ok {
			panic("添加的路由存在冲突，该路由为" + v.RouterPattern)
		} else {
			c[v.patternRe.String()] = v
		}
	}
	return c
}

// 计算路由信息
func (c *ControllerActionInfo) makePattern() {
	pathPattern := c.RouterPattern
	if len(pathPattern) == 0 {
		panic("设置的路由匹配模式不能为空")
	}
	if pathPattern[0] != '/' {
		pathPattern = "/" + pathPattern
	}
	waitSortMap := map[string]string{}
	waitSortArr := make([]int, 0)
	sortMap := map[string]int{}
	f := regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
	t := f.FindAllStringIndex(pathPattern, -1)
	if len(t) > 1 {
		panic("路由注册时发现{Controller}使用超过1次")
	}
	if len(t) == 1 {
		waitSortMap[strconv.Itoa(t[0][0])] = "Controller"
		waitSortArr = append(waitSortArr, t[0][0])
	}
	f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
	t = f.FindAllStringIndex(pathPattern, -1)
	if len(t) > 1 {
		panic("路由注册时发现{Action}使用超过1次")
	}
	if len(t) == 1 {
		waitSortMap[strconv.Itoa(t[0][0])] = "Action"
		waitSortArr = append(waitSortArr, t[0][0])
	}
	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*(int|string|allString)+[\ ]*}`)
	t = f.FindAllStringIndex(pathPattern, -1)
	for _, v := range t {
		sl := strings.Trim(strings.Split(pathPattern[v[0]+1:v[1]], ":")[0], " ")
		for _, v := range waitSortMap {
			if v == sl {
				panic("路由注册时" + c.RouterPattern + "时发现{" + sl + "}使用超过1次")
			}
		}
		println(sl)
		waitSortArr = append(waitSortArr, v[0])
		waitSortMap[strconv.Itoa(v[0])] = sl
	}
	sort.Ints(waitSortArr)
	for k, v := range waitSortArr {
		sortMap[waitSortMap[strconv.Itoa(v)]] = k + 1
	}
	fmt.Println(fmt.Sprintf("%+v", sortMap))
	f = regexp.MustCompile(`{[\ ]*Controller[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `([\w+$]+)`)
	f = regexp.MustCompile(`{[\ ]*Action[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `([\w+$]+)`)
	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*int[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `(-?[1-9]\d+)`)
	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*string[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `([\w+$]+)`)

	f = regexp.MustCompile(`{[\ ]*[a-zA-Z][\w+$]+[\ ]*:[\ ]*allString[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `(.*?)`)

	f = regexp.MustCompile(`{[\ ]*int[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `-?[1-9]\d+`)
	f = regexp.MustCompile(`{[\ ]*string[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `[\w+$]+`)

	f = regexp.MustCompile(`{[\ ]*allString[\ ]*}`)
	pathPattern = f.ReplaceAllString(pathPattern, `.*?`)

	c.patternRe = regexp.MustCompile("^" + pathPattern + "$")
	c.patternMap = sortMap
}

func newTree() *tree {
	tree := new(tree)
	tree.ControllerList = map[string]map[string]*ControllerInfo{}
	return tree
}

// 控制器注册
func (t *tree) addPathTree(controllerName string, controllerAction string, controllerFunc reflect.Type, ControllerActionParameterStruct reflect.Type) {
	controllerInfo := new(ControllerInfo)
	controllerInfo.ControllerAction = controllerAction
	controllerInfo.ControllerName = controllerName
	controllerInfo.ControllerFunc = controllerFunc
	controllerInfo.ControllerActionParameterStruct = ControllerActionParameterStruct
	if _, ok := t.ControllerList[controllerName]; !ok {
		t.ControllerList[strings.ToLower(controllerName)] = map[string]*ControllerInfo{}
	}

	t.ControllerList[strings.ToLower(controllerName)][strings.ToLower(controllerAction)] = controllerInfo
}

//  根据控制器名字和动作名字获取控制器
func (t *tree) getControllerInfoByControllerNameControllerAction(controllerName string, controllerAction string) (*ControllerInfo, bool) {
	if v, ok := t.ControllerList[controllerName]; ok {
		if v, ok := v[controllerAction]; ok {
			return v, true
		}
	}
	return nil, false
}
