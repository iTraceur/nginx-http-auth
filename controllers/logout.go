package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Get() {
	clientIP := this.Ctx.Input.IP()
	username, ok := this.GetSession("uname").(string)
	if !ok {
		this.Redirect("/passport/login", 302)
		return
	}

	_ = this.DestroySession()
	logs.Notice(fmt.Sprintf("%s - user(%s) logout", clientIP, username))
	this.Redirect("/passport/login", 302)
}
