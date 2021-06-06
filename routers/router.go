package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"nginx-http-auth/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/passport/login", &controllers.LoginController{})
	beego.Router("/passport/logout", &controllers.LogoutController{})
	beego.Router("/users", &controllers.UserController{}, "get:List")
	beego.Router("/users/add", &controllers.UserController{}, "get:Create;post:Create")
	beego.Router("/users/edit/:id", &controllers.UserController{}, "get:Update;post:Update")
	beego.Router("/users/delete/:id", &controllers.UserController{}, "get:Delete;delete:Delete")
	beego.Router("/auth-proxy", &controllers.AuthProxyController{})
	beego.Router("/api/v1/:control", &controllers.ControlController{})
	beego.ErrorController(&controllers.ErrorController{})
}
