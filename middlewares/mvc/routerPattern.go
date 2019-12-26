package mvc

type freeFishUrl struct {
	controllerName   string
	controllerAction string
	OtherKeyMap      map[string]interface{}
	ControllerInfo   *controllerInfo
}

// 获取控制器名称
func (f *freeFishUrl) GetControllerName(c *ControllerActionInfo) string {
	if v, ok := f.OtherKeyMap["Controller"]; ok {
		return v.(string)
	} else {
		return c.controllerName
	}
}

// 获取动作名称
func (f *freeFishUrl) GetControllerAction(c *ControllerActionInfo) string {
	if v, ok := f.OtherKeyMap["Action"]; ok {
		return v.(string)
	} else {
		return c.actionName
	}
}
