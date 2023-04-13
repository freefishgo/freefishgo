// Copyright 2019 freefishgo Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mvc

import (
	"regexp"
	"strings"
)

type freeFishUrl struct {
	controllerName   string
	controllerAction string
	OtherKeyMap      map[string]interface{}
	ControllerInfo   *controllerInfo
}

// GetControllerName 获取控制器名称
func (f *freeFishUrl) GetControllerName(c *ActionRouter) string {
	if v, ok := f.OtherKeyMap["Controller"]; ok {
		return v.(string)
	} else {
		return c.controllerName
	}
}

// GetControllerAction 获取动作名称
func (f *freeFishUrl) GetControllerAction(c *ActionRouter) string {
	if v, ok := f.OtherKeyMap["Action"]; ok {
		return v.(string)
	} else {
		return c.actionName
	}
}

// / <summary>
// / 转义字符串中所有正则特殊字符
// / </summary>
// / <param name="input">传入字符串</param>
// / <returns></returns>
func filterRegexpString(input string) string {
	input = strings.Replace(input, "\\", "\\\\", -1) //先替换“\”，不然后面会因为替换出现其他的“\”
	//r := regexp.MustCompile("[\\*\\.\\?\\+\\$\\^\\[\\]\\(\\)\\{\\}\\|\\/]")
	r := regexp.MustCompile("[\\*\\.\\?\\+\\$\\^\\[\\]\\(\\)\\|]")
	arrList := r.FindAllStringSubmatch(input, -1)
	list := map[string]bool{}
	for _, v := range arrList {
		if _, ok := list[v[0]]; ok {
			continue
		}
		input = strings.Replace(input, v[0], "\\"+v[0], -1)
		list[v[0]] = true
	}
	return input
}
