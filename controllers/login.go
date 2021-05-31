package controllers

import (
	"fmt"
	"html/template"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/captcha"

	"nginx-http-auth/utils"
)

type LoginController struct {
	beego.Controller
}

func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
	cpt.ChallengeNums = 4
	cpt.StdWidth = 110
	cpt.StdHeight = 48
	cpt.Expiration = 120 * time.Second
}

var cpt *captcha.Captcha

func (this *LoginController) Get() {
	logtime := time.Now().Format("02/Jan/2006 03:04:05")
	clientIP := this.Ctx.Input.IP()
	directIPS, err := beego.AppConfig.Strings("ipControl::direct")
	denyIPS, err := beego.AppConfig.Strings("ipControl::deny")
	timeDirect, err := beego.AppConfig.Strings("timeControl::direct")
	timeDeny, err := beego.AppConfig.Strings("timeControl::deny")
	if err != nil {
		logs.Error(err.Error())
		this.Ctx.Output.SetStatus(500)
		this.Ctx.WriteString("Internal Server Error")
		return
	}

	if utils.IpCheck(clientIP, denyIPS) {
		logs.Notice(fmt.Sprintf("%s - - [%s] Login failed: IP %s is not allowed", clientIP, logtime, clientIP))
		this.Abort("401")
	}

	if utils.IpCheck(clientIP, directIPS) {
		this.SetSession("uname", clientIP)
		logs.Notice(fmt.Sprintf("%s - %s [%s] Login successed: direct IP", clientIP, clientIP, logtime))
		this.Ctx.Redirect(302, "/")
		return
	}

	if utils.TimeCheck(timeDeny) {
		logs.Notice(fmt.Sprintf("%s - - [%s] Login failed: access is not allowed for the current time period", clientIP, logtime))
		this.Abort("403")
	}
	if utils.TimeCheck(timeDirect) {
		this.SetSession("uname", "timeDirect")
		logs.Notice(fmt.Sprintf("%s - %s [%s] Login successed: direct access period", clientIP, "timeDirect", logtime))
		this.Ctx.Redirect(302, "/")
		return
	}

	uname := this.GetSession("uname")
	if uname != nil {
		this.Ctx.Redirect(302, "/")
		return
	}

	target := this.Ctx.Input.Header("X-Target")
	getTarget := this.GetString("target")
	if target == "" && getTarget == "" {
		target = "/"
	}
	if getTarget != "" {
		target = getTarget
	}
	this.Data["target"] = target
	loginFailed := this.GetSession("loginFailed")
	if loginFailed != nil {
		this.Data["captcha"] = true
	}
	var msg string
	switch loginFailed {
	case "1":
		msg = "用户名或密码错误！"
	case "2":
		msg = "用户不存在或已禁用！"
	case "3":
		msg = "验证码错误！"
	}
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.Data["msg"] = msg
	this.TplName = "login.tpl"
}

func (this *LoginController) Post() {
	logtime := time.Now().Format("02/Jan/2006 03:04:05")
	clientIP := this.Ctx.Input.IP()
	this.Ctx.Request.ParseForm()
	username := this.Ctx.Request.Form.Get("username")
	password := this.Ctx.Request.Form.Get("password")
	target := this.Ctx.Request.Form.Get("target")
	loginFailed := this.GetSession("loginFailed")
	if loginFailed != nil {
		if !cpt.VerifyReq(this.Ctx.Request) {
			this.SetSession("loginFailed", "3")
			logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: Captcha Wrong", clientIP, logtime))
			this.Ctx.Redirect(302, fmt.Sprintf("/passport/login?target=%s", target))
			return
		}
	}

	allowUsers, err := beego.AppConfig.Strings("userControl::allow")
	denyUsers, err := beego.AppConfig.Strings("userControl::deny")
	if err != nil {
		logs.Error(err.Error())
		this.Ctx.Output.SetStatus(500)
		this.Ctx.WriteString("Internal Server Error")
		return
	}

	if len(allowUsers) > 0 && !utils.InSlice(username, allowUsers) || len(denyUsers) > 0 && utils.InSlice(username, denyUsers) {
		this.SetSession("loginFailed", "2")
		logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed: user %s is not allowed", clientIP, logtime, username))
		this.Ctx.Redirect(302, fmt.Sprintf("/passport/login?target=%s", target))
		return
	}

	if utils.HttpAuth(username, password) {
		//登录成功设置session
		if target == "" || target == "/passport/login" {
			logs.Warning(fmt.Sprintf("%s - - [%s] Login Successed: Missing X-Target", clientIP, logtime))
			this.Ctx.Redirect(302, "/")
		}
		this.SetSession("uname", username)
		logs.Notice(fmt.Sprintf("%s - %s [%s] Login Successed", clientIP, username, logtime))
		this.Ctx.Redirect(302, target)
	} else {
		this.SetSession("loginFailed", "1")
		logs.Notice(fmt.Sprintf("%s - - [%s] Login Failed", clientIP, logtime))
		this.Ctx.Redirect(302, fmt.Sprintf("/passport/login?target=%s", target))
	}
}
