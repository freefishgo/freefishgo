package mvc

import "testing"

func TestFilterRegexpString(t *testing.T) {
	data := []struct {
		f string
		r string
	}{
		{
			f: ".*$.a*$fdfs:f??.*$?",
			r: "\\.\\*\\$\\.a\\*\\$fdfs:f\\?\\?\\.\\*\\$\\?",
		},
		{
			f: ".",
			r: "\\.",
		},
		{
			f: "?",
			r: "\\?",
		},
		{
			f: "(",
			r: "\\(",
		},
		{
			f: ")",
			r: "\\)",
		},
		{
			f: "[",
			r: "\\[",
		},
		{
			f: "]",
			r: "\\]",
		},
		{
			f: "{",
			r: "{",
		},
		{
			f: "}",
			r: "}",
		},
		{
			f: ":",
			r: ":",
		},
		{
			f: ".*$?",
			r: "\\.\\*\\$\\?",
		},
	}
	for _, v := range data {
		if filterRegexpString(v.f) != v.r {
			t.Errorf("测试失败%+v,因为%s不等于%s", v, filterRegexpString(v.f), v.r)
		}
	}
}
