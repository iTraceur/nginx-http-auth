package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type AuthProxyController struct {
	beego.Controller
}

func (this *AuthProxyController) Get() {
	this.TplName = "index.html"
	data := map[string]string{
		"title": "Are you ok?",
		"p":     "I'm ok.",
	}
	this.Ctx.Output.Header("Cache-Control", "no-cache")
	this.Data["data"] = data
}
