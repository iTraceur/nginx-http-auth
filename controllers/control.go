package controllers

import (
	"fmt"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/toolkits/file"

	"nginx-http-auth/g"
)

type ControlController struct {
	beego.Controller
}

func (this *ControlController) Get() {
	// 获取客户端IP
	clientIP := this.Ctx.Input.IP()

	// 获取用户Session
	username, ok := this.GetSession("uname").(string)
	if !ok {
		this.Redirect("/passport/login", 302)
		return
	}

	// 获取管理用户配置
	manageUsers, err := beego.AppConfig.Strings("manageUsers")
	if err != nil {
		logs.Warn(fmt.Sprintf("%s - get manage users failed: %s", clientIP, err.Error()))
		manageUsers = []string{"admin"}
	}

	// 管理用户校验
	if !utils.InSlice(username, manageUsers) {
		logs.Warn(fmt.Sprintf("%s - %s - access to control API was denied", clientIP, username))
		this.Abort("453")
	}

	// 获取管理类型
	control := this.Ctx.Input.Param(":control")
	switch control {
	case "version": // 获取版本
		_ = this.Ctx.Output.Body([]byte(g.VERSION))
	case "health": // 探活
		_ = this.Ctx.Output.Body([]byte("ok"))
	case "config": // 获取配置信息
		var json map[string]interface{}
		err = beego.AppConfig.Unmarshaler("", &json)
		if err != nil {
			logs.Error(fmt.Sprintf("%s - get config info failed: %s", clientIP, err.Error()))
			this.Abort("550")
		}
		this.Data["json"] = json
		_ = this.ServeJSON()
	case "reload": // 重新加载配置
		err := beego.LoadAppConfig("ini", filepath.Join(file.SelfDir(), "conf/app.conf"))
		if err != nil {
			logs.Error(fmt.Sprintf("%s - config reload failed: %s", clientIP, err.Error()))
			_ = this.Ctx.Output.Body([]byte("config reload failed"))
		} else {
			logs.Info(fmt.Sprintf("%s - config Reloaded", clientIP))
			_ = this.Ctx.Output.Body([]byte("config reloaded"))
		}
	}
}
