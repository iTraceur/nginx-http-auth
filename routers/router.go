package routers

import (
	"nginx-http-auth/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/passport/login", &controllers.LoginController{})
	beego.Router("/passport/logout", &controllers.LogoutController{})
	beego.Router("/auth-proxy", &controllers.AuthProxyController{})
	beego.Router("/api/v1/:control", &controllers.ControlController{})
	beego.ErrorController(&controllers.ErrorController{})
}
