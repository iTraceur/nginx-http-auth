package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"nginx-http-auth/g"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.TplName = "index.html"
	data := map[string]string{
		"title": "nginx-http-auth",
		"p":     fmt.Sprintf("Version: %s", g.VERSION),
	}
	this.Data["data"] = data
}
