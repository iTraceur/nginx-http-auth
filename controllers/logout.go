package controllers

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Get() {
	logtime := time.Now().Format("02/Jan/2006 03:04:05")
	clientIP := this.Ctx.Input.IP()
	uname := this.GetSession("uname")
	if uname != nil {
		this.DelSession("uname")
		this.DelSession("loginFailed")
		logs.Notice(fmt.Sprintf("%s - %s [%s] Logout successed", clientIP, uname, logtime))
	}
	this.Ctx.Redirect(302, "/")
}
