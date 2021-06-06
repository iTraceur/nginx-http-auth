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
	username := this.GetSession("uname")
	if username != nil {
		_ = this.DelSession("uname")
		_ = this.DelSession("loginFailed")
		_ = this.DelSession("userAuthHash")
		logs.Notice(fmt.Sprintf("%s - user(%s) logout", clientIP, username))
	}
	this.Redirect("/passport/login", 302)
}
