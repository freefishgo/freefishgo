package mvc

import (
	"testing"
)

func TestDoStruct(t *testing.T) {
	type roles struct {
		RoleId   int
		RoleName string
	}
	type User struct {
		Name     string `json:"Name"`
		Age      int
		Email    string
		NickName string
		Telphone int
		Test     []int
		Roles    *roles
	}
	ul := map[string]interface{}{}
	ul["Name"] = "name"
	ul["Age"] = "40"
	ul["Test"] = []string{"464", "4464"}
	ul["Roles"] = `{"RoleId":1001,"RoleName":"administrator"}`
	u := new(User)
	MapStringToStruct(u, ul)
	if u.Name != "name" || u.Age != 40 || u.Roles.RoleId != 1001 || u.Test[1] != 4464 {
		t.Error("格式化失败")
	}
}
