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
	freeFishGo "github.com/freefishgo/freefishgo"
)

type MvcWebConfig struct {
	// 是否启用开发模式
	IsDevelopment bool
	// html模板引擎变量左标记符号
	TemplateLeft string
	// html模板引擎变量右标记符号
	TemplateRight string
	// html模板的父目录
	ViewsPath string

	homeController string
	indexAction    string
	// 是否在Mvc框架最末端捕获Panic，以取代Mvc框架的处理Panic函数
	RecoverPanic bool
	// 捕获Panic的处理函数
	RecoverFunc func(ctx *freeFishGo.HttpContext, e error, Stack []byte)
}

// 实例化一个 MvcWebConfig
func NewWebConfig() *MvcWebConfig {
	return &MvcWebConfig{ViewsPath: "", TemplateLeft: "{{", TemplateRight: "}}", IsDevelopment: false}
}
