package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"nginx-http-auth/g"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	uname := this.GetSession("uname")
	if uname == nil {
		this.Ctx.Redirect(302, "/passport/login")
		return
	}
	this.Ctx.Output.Body([]byte("nginx-http-auth, version " + g.VERSION))
}
