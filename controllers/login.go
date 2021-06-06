package controllers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
	beeUtils "github.com/beego/beego/v2/core/utils"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/captcha"

	"nginx-http-auth/form"
	"nginx-http-auth/models"
	"nginx-http-auth/utils"
)

type LoginController struct {
	beego.Controller
}

// 初始化验证码插件
func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
	cpt.ChallengeNums = 4
	cpt.StdWidth = 110
	cpt.StdHeight = 48
	cpt.Expiration = 120 * time.Second
}

var cpt *captcha.Captcha

func (this *LoginController) ClearLoginInfo() {
	_ = this.DelSession("uname")
	_ = this.DelSession("loginFailed")
	_ = this.DelSession("userAuthHash")
}

func (this *LoginController) Get() {
	this.Layout = "main.html"
	this.TplName = "login.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "login-head.html"
	this.LayoutSections["Footer"] = "footer.html"

	// 获取客户端IP
	clientIP := this.Ctx.Input.IP()

	// 获取IP绑定、直通IP列表、限制IP列表、直通时间段、限制时间段配置
	var directIPs, denyIPs, timeDirect, timeDeny []string
	var err error
	if directIPs, err = beego.AppConfig.Strings("ipControl::direct"); err != nil {
		logs.Warn(fmt.Sprintf("%s - get direct ips failed: %s", clientIP, err.Error()))
		directIPs = []string{}
	}
	if denyIPs, err = beego.AppConfig.Strings("ipControl::deny"); err != nil {
		logs.Warn(fmt.Sprintf("%s - get deny ips failed: %s", clientIP, err.Error()))
		denyIPs = []string{}
	}
	if timeDirect, err = beego.AppConfig.Strings("timeControl::direct"); err != nil {
		logs.Warn(fmt.Sprintf("%s - get time direct failed: %s", clientIP, err.Error()))
		timeDirect = []string{}
	}
	if timeDeny, err = beego.AppConfig.Strings("timeControl::deny"); err != nil {
		logs.Warn(fmt.Sprintf("%s - get time deny failed: %s", clientIP, err.Error()))
		timeDeny = []string{}
	}

	// 是否忽略直通，强制登录
	forceLogin, err := this.GetBool("force_login")
	if err != nil {
		forceLogin = false
	}

	// 强制登录前需要清空直通登录后的登录状态
	if forceLogin {
		this.ClearLoginInfo()
	}

	// 获取目标跳转路径
	target := this.Ctx.Input.Header("X-Target")
	getTarget := this.GetString("target")
	if getTarget != "" {
		target = getTarget
	}
	this.Data["target"] = target

	// 这处的目标跳转路径处理需要放在"this.Data["target"] = target"赋值之后，
	// 目的是为了确保在开启了本地认证方式下管理员直接登录时能够正确跳转到用户管理页面
	if target == "" || target == "/passport/login" {
		target = "/"
	}

	// 拒绝在“限制IP列表”中的客户端IP访问
	if utils.IpCheck(clientIP, denyIPs) {
		logs.Warn(fmt.Sprintf("%s - login denied: IP deny", clientIP))
		this.Abort("451")
	}

	// 允许在“直通IP列表”中的客户端IP访问
	if !forceLogin && utils.IpCheck(clientIP, directIPs) {
		_ = this.SetSession("uname", clientIP)
		logs.Info(fmt.Sprintf("%s - login directed: IP direct", clientIP))
		this.Redirect(target, 302)
		return
	}

	// 拒绝“限制时间段”内的时间访问
	if utils.TimeCheck(timeDeny) {
		logs.Warn(fmt.Sprintf("%s - Login denied: time deny", clientIP))
		this.Abort("452")
	}

	// 允许“直通时间段”的时间访问
	if !forceLogin && utils.TimeCheck(timeDirect) {
		_ = this.SetSession("uname", "timeDirect")
		logs.Info(fmt.Sprintf("%s - login directed: time direct", clientIP))
		this.Redirect(target, 302)
		return
	}

	// 如果用户已登录，直接跳转目标页面
	uname := this.GetSession("uname")
	if uname != nil {
		this.Redirect(target, 302)
		return
	}

	// 如之前登录失败过，则本次登录需要输入验证
	loginFailed := this.GetSession("loginFailed")
	if loginFailed != nil {
		this.Data["captcha"] = true
	}
	// 登录失败提示信息
	var msg string
	switch loginFailed {
	case "1":
		msg = "用户名或密码错误！"
	case "2":
		msg = "用户不存在或已禁用！"
	case "3":
		msg = "验证码错误！"
	}
	this.Data["msg"] = msg

	// 生成 XSRF Token
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
}

func (this *LoginController) Post() {
	// 获取客户端IP
	clientIP := this.Ctx.Input.IP()

	// 解析登录表单
	_ = this.Ctx.Request.ParseForm()
	username := this.Ctx.Request.Form.Get("username")
	password := this.Ctx.Request.Form.Get("password")
	target := this.Ctx.Request.Form.Get("target")

	// 检验验证码，只在用户登录失败后再登录时需要校验验证码
	loginFailed := this.GetSession("loginFailed")
	if loginFailed != nil {
		if !cpt.VerifyReq(this.Ctx.Request) {
			_ = this.SetSession("loginFailed", "3")
			logs.Warn(fmt.Sprintf("%s - %s - login failed: captcha wrong", clientIP, username))
			this.Redirect(fmt.Sprintf("/passport/login?target=%s", target), 302)
			return
		}
	}

	// 获取管理用户配置, 用户控制配置
	var manageUsers, allowUsers, denyUsers []string
	var err error
	if manageUsers, err = beego.AppConfig.Strings("manageUsers"); err != nil {
		logs.Warn(fmt.Sprintf("%s - get manage users failed: %s", clientIP, err.Error()))
		manageUsers = []string{"admin"}
	}
	if allowUsers, err = beego.AppConfig.Strings("userControl::allow"); err != nil {
		logs.Warn(fmt.Sprintf("%s - get allow users failed: %s", clientIP, err.Error()))
		allowUsers = []string{}
	}
	if denyUsers, err = beego.AppConfig.Strings("userControl::deny"); err != nil {
		logs.Warn(fmt.Sprintf("%s - get demy users failed: %s", clientIP, err.Error()))
		denyUsers = []string{}
	}

	// 用户访问控制，当“允许列表”不为空时，则只允许“允许列表”中的用户登录，否则允许不在“拒绝列表”中的用户登录
	if len(allowUsers) > 0 && !beeUtils.InSlice(username, allowUsers) || len(denyUsers) > 0 && beeUtils.InSlice(username, denyUsers) {
		logs.Warn(fmt.Sprintf("%s - login failed: user %s is not allowed", clientIP, username))
		this.Abort("453")
		return
	}

	// 获取认证方式
	provider, err := beego.AppConfig.String("authProvider")
	if err != nil || provider == "" {
		provider = "local"
	}

	// 本地认证
	if provider == "local" {
		// 表单校验
		uValidator := form.UserAuthValidator{
			Username: username,
			Password: password,
		}
		errMap := uValidator.Valid()
		if len(errMap) == 0 {
			// 查询用户名对应的用户
			user, err := models.GetUserByUsername(username)
			if err != nil { // 查询用户名对应的用户失败
				_ = this.SetSession("loginFailed", "2")
				logs.Error(fmt.Sprintf("%s - %s - login failed: %s", clientIP, username, err.Error()))
				this.Redirect(fmt.Sprintf("/passport/login?target=%s", target), 302)
				return
			}
			if user == nil || !user.Active { // 用户名对应的用户不存在或已禁用
				_ = this.SetSession("loginFailed", "2")
				logs.Warn(fmt.Sprintf("%s - %s - login failed: user is not active", clientIP, username))
				this.Redirect(fmt.Sprintf("/passport/login?target=%s", target), 302)
				return
			}

			ipBinding, err := beego.AppConfig.Bool("ipBinding")
			if err != nil {
				ipBinding = false
			}
			if !beeUtils.InSlice(username, manageUsers) && ipBinding && clientIP != user.ClientIp {
				logs.Warn(fmt.Sprintf("%s - %s - login failed: client ip (%s) mismatch", clientIP, username, user.ClientIp))
				this.Abort("450")
			}

			// 校验密码
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				_ = this.SetSession("loginFailed", "1")
				logs.Warn(fmt.Sprintf("%s - %s - login failed: %s", clientIP, username, err.Error()))
				this.Redirect(fmt.Sprintf("/passport/login?target=%s", target), 302)
				return
			}

			// 临时用户只能登录一次
			if user.Temporary {
				user.Active = false
				_ = models.SaveUser(user)
			}

			// 设置用户认证哈希，用于用户密码更改后强制其重新登录
			authKey := utils.Md5sum([]byte(user.Password))
			_ = this.SetSession("userAuthHash", authKey)
		}
	} else { // 远程认证
		if !utils.HttpAuth(username, password) {
			_ = this.SetSession("loginFailed", "1")
			logs.Warn(fmt.Sprintf("%s - %s - login failed: auth api returnd failure", clientIP, username))
			this.Redirect(fmt.Sprintf("/passport/login?target=%s", target), 302)
			return
		}
	}

	// 登录成功设置session
	logs.Notice(fmt.Sprintf("%s - %s login successed", clientIP, username))
	_ = this.SetSession("uname", username)

	if target == "" || target == "/passport/login" {
		// 无目标跳转路径时，如用户为管理用户且认证方式为本地，则跳转到用户管理页面，否则跳转到首页
		if beeUtils.InSlice(username, manageUsers) && provider == "local" {
			this.Redirect("/users", 302)
		} else {
			logs.Warn(fmt.Sprintf("%s - login successed: missing X-Target", clientIP))
			this.Redirect("/", 302)
		}
	} else {
		// 有目标跳转路径时，则跳转到对应的路径页面
		this.Redirect(target, 302)
	}
}
