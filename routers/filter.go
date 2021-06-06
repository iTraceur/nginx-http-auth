package routers

import (
	"strings"

	"github.com/beego/beego/v2/core/utils"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.InsertFilter("/*", beego.BeforeRouter, LoginFilter)
	beego.InsertFilter("/users/*", beego.BeforeRouter, UserFilter)
	beego.InsertFilter("/api/*", beego.BeforeRouter, UserFilter)
}

// 登录拦截器
var LoginFilter = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("uname").(string)
	if !ok {
		uri := strings.Split(ctx.Request.RequestURI, "?")[0]
		uri = strings.TrimRight(uri, "/")
		switch uri {
		case "/auth-proxy", "/users/add", "/users/edit", "/users/delete":
			ctx.Abort(401, "401")
		case "/passport/login":
		default:
			ctx.Redirect(302, "/passport/login")
		}
	}
}

// 访问控制拦截器
var UserFilter = func(ctx *context.Context) {
	uri := ctx.Request.RequestURI
	// 用户Session
	username, _ := ctx.Input.Session("uname").(string)

	// 获取管理用户配置
	manageUsers, err := beego.AppConfig.Strings("manageUsers")
	if err != nil {
		manageUsers = []string{"admin"}
	}

	// 只有认证方式设置为"local"时，才能访问"/users"及其子路由
	if strings.HasPrefix(uri, "/users") {
		if provider, _ := beego.AppConfig.String("authProvider"); provider != "local" {
			ctx.Output.SetStatus(403)
			_ = ctx.Output.Body([]byte("需开启本地认证！"))
		}
	}

	// 只有在管理员用户列表中的用户才可以访问以"/users"或"/api"开头的路由
	if strings.HasPrefix(uri, "/users") || strings.HasPrefix(uri, "/api") {
		if !utils.InSlice(username, manageUsers) {
			ctx.Output.SetStatus(403)
			_ = ctx.Output.Body([]byte("非管理员，请求被拒绝！"))
		}
	}
}
