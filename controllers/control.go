package controllers

import (
	"fmt"
	"path/filepath"
	"reflect"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/toolkits/file"

	"nginx-http-auth/g"
	"nginx-http-auth/utils"
)

type ControlController struct {
	beego.Controller
}

func (this *ControlController) Get() {
	clientIP := this.Ctx.Input.IP()
	logtime := time.Now().Format("02/Jan/2006 03:04:05")

	uname := this.GetSession("uname")
	if uname == nil {
		this.Ctx.Redirect(302, "/passport/login")
		return
	}

	controlUsers, err := beego.AppConfig.Strings("controlUsers")
	if err != nil {
		logs.Error(err.Error())
		this.Ctx.Output.SetStatus(500)
		this.Ctx.WriteString("Internal Server Error")
		return
	}
	if !utils.InSlice(uname.(string), controlUsers) {
		logs.Debug(uname.(string), controlUsers, controlUsers[0], reflect.TypeOf(controlUsers[0]))
		this.Ctx.Output.SetStatus(401)
		this.Ctx.Output.Body([]byte("Not Allowed"))
		return
	}

	control := this.Ctx.Input.Param(":control")
	switch control {
	case "version":
		this.Ctx.Output.Body([]byte(g.VERSION))
	case "health":
		this.Ctx.Output.Body([]byte("ok"))
	case "config":
		var json map[string]interface{}
		err = beego.AppConfig.Unmarshaler("", &json)
		if err != nil {
			logs.Error(err.Error())
			this.Ctx.Output.SetStatus(500)
			this.Ctx.WriteString("Internal Server Error")
			return
		}
		this.Data["json"] = json
		this.ServeJSON()
	case "reload":
		err := beego.LoadAppConfig("ini", filepath.Join(file.SelfDir(), "conf/app.conf"))
		if err != nil {
			logs.Error(fmt.Sprintf("%s - - [%s] Config reload failed: %s", clientIP, logtime, err.Error()))
			this.Ctx.Output.Body([]byte("config reload failed"))
		} else {
			logs.Notice(fmt.Sprintf("%s - - [%s] Config Reloaded", clientIP, logtime))
			this.Ctx.Output.Body([]byte("config reloaded"))
		}
	}
}
