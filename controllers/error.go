package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type ErrorController struct {
	beego.Controller
}

func (this *ErrorController) Prepare() {
	this.Layout = "main-simple.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "deny-head.html"
}

func (this *ErrorController) Error401() {
	this.Data["title"] = "未授权页面"
	this.Data["content"] = "未登录或认证状态失效"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error403() {
	this.Data["title"] = "拒绝访问"
	this.Data["content"] = "限制访问，请求被拒绝"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error404() {
	this.Data["title"] = "Ooops!"
	this.Data["content"] = "页面未找到"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error405() {
	this.Data["title"] = "非法请求"
	this.Data["content"] = "不允许的请求方法"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error417() {
	this.Data["title"] = "非法请求"
	this.Data["content"] = "XSRF校验失败"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error422() {
	this.Data["title"] = "非法请求"
	this.Data["content"] = "请求未提交XSRF参数"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error450() {
	this.Data["title"] = "拒绝访问"
	this.Data["content"] = "IP不匹配，请求被拒绝"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error451() {
	this.Data["title"] = "拒绝访问"
	this.Data["content"] = "限制IP访问，请求被拒绝"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error452() {
	this.Data["title"] = "拒绝访问"
	this.Data["content"] = "当前时间段不允许访问"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error453() {
	this.Data["title"] = "拒绝访问"
	this.Data["content"] = "限制用户访问，请求被拒绝"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error500() {
	this.Data["title"] = "服务异常"
	this.Data["content"] = "服务器内部错误"
	this.TplName = "deny.html"
}

func (this *ErrorController) Error550() {
	this.Data["title"] = "服务异常"
	this.Data["content"] = "服务未正确配置"
	this.TplName = "deny.html"
}
