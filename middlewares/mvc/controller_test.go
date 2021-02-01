package mvc

import (
	"testing"
)

func TestReplaceActionNameIgone(T *testing.T) {
	t := []struct {
		f string
		r string
	}{
		{f: "nihao", r: "nihao"},
		{f: "nihao", r: "nihao"},
		{f: "niHao", r: "niHao"},
		{f: "nihaoget", r: "nihao"},
		{f: "nihaoGet", r: "nihao"},
		{f: "nIhaoGeT", r: "nIhao"},
	}
	for _, v := range t {
		if replaceActionNameIgnore(v.f) != v.r {
			T.Errorf("测试失败%+v,因为%s不等于%s", v, replaceActionNameIgnore(v.f), v.r)
		}
	}
}
