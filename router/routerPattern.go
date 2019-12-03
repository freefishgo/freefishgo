package router

type Pattern struct {
	PatternString  string
	ControllerName string
	ActionName     string
	//正则匹配出来的变量地址映射变量映射
	PatternMap map[string]int
}
type freeFishUrl struct {
	ControllerName   string
	ControllerAction string
	OtherKeyMap      map[string]interface{}
}
