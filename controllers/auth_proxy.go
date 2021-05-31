package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type AuthProxyController struct {
	beego.Controller
}

func (this *AuthProxyController) Get() {
	this.Ctx.Output.Header("Cache-Control", "no-cache")
	uname := this.GetSession("uname")
	if uname == nil {
		this.Ctx.Abort(401, "401")
		return
	}
	this.Ctx.Output.Body([]byte("ok"))
}
