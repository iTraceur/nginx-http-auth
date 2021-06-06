package utils

import "github.com/beego/beego/v2/server/web"

func init() {
	_ = web.AddFuncMap("add", addition)
	_ = web.AddFuncMap("sub", subtract)
}

func addition(num, n int) int {
	return num + n
}

func subtract(num, n int) int {
	return num - n
}
