package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type ErrorController struct {
	beego.Controller
}

func (this *ErrorController) Error401() {
	this.Data["content"] = "限制访问，请求被拒绝"
	this.TplName = "deny.tpl"
}

func (this *ErrorController) Error403() {
	this.Data["content"] = "当前时间段不允许访问"
	this.TplName = "deny.tpl"
}
