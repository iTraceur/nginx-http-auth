package controllers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"nginx-http-auth/form"
	"nginx-http-auth/models"
	"nginx-http-auth/utils"
)

type UserStatusController struct {
	beego.Controller
	User *models.User
}

func (this *UserStatusController) ClearLoginInfo() {
	_ = this.DestroySession()
	this.User = nil
}

func (this *UserStatusController) Prepare() {
	username, ok := this.GetSession("uname").(string)
	if !ok {
		this.ClearLoginInfo()
		this.Redirect("/passport/login", 302)
		return
	}

	userAuthHash, ok := this.GetSession("userAuthHash").(string)
	if !ok {
		this.ClearLoginInfo()
		this.Redirect("/passport/login", 302)
		return
	}

	user, err := models.GetUserByUsername(username)
	if err != nil {
		logs.Error(fmt.Sprintf("get user by username(%s) failed: %s", username, err.Error()))
		this.ClearLoginInfo()
		this.Redirect("/passport/login", 302)
		return
	}

	hash := utils.Md5sum([]byte(user.Password))
	if hash != userAuthHash {
		logs.Warn(fmt.Sprintf("%s user auth hash dismatch, prev: %s, new: %s", username, userAuthHash, hash))
		this.ClearLoginInfo()
		this.Redirect("/passport/login", 302)
		return
	}

	this.User = user
	this.Data["user"] = user
}

type UserController struct {
	UserStatusController
}

// 用户列表
func (this *UserController) List() {
	this.Layout = "main.html"
	this.TplName = "user/list.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["HtmlHead"] = "user/head.html"
	this.LayoutSections["Nav"] = "nav.html"
	this.LayoutSections["Footer"] = "footer.html"
	this.LayoutSections["Script"] = "user/list-script.html"

	// 获取页码及分页大小
	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	}
	pageSize, err := this.GetInt("page_size")
	if err != nil {
		pageSize = 10
	}

	// 获取用户总数
	count := models.GetUserCount()
	// 获取分页器
	paginator := utils.Paginate(count, page, pageSize)
	if paginator.TotalPage < page {
		this.Redirect(fmt.Sprintf("/users?page=%d", paginator.TotalPage), 302)
		return
	}

	// 过滤分页用户列表
	userList := models.LimitList(page, pageSize)

	this.Data["paginator"] = paginator
	this.Data["users"] = userList
}

// 添加用户
func (this *UserController) Create() {
	this.Layout = "user/user-modal.html"
	this.TplName = "user/create.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Footer"] = "footer.html"
	this.LayoutSections["Script"] = "user/create-script.html"

	// GET 请求获取用户添加表单
	if this.Ctx.Input.Method() == "GET" {
		// 生成 XSRF Token
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	} else {
		// 初始化响应
		res := make(map[string]interface{})
		this.Data["json"] = res

		// 解析表单
		userForm := form.UserForm{}
		if err := this.ParseForm(&userForm); err != nil {
			res["errMsg"] = "提交的数据不合法"
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
			return
		}

		// 构造用户创建表单校验器，并进行校验
		uValidator := form.UserCreateValidator{
			Username: userForm.Username,
			Password: userForm.Password,
		}
		errMap := uValidator.Valid()
		if len(errMap) != 0 {
			res["errMap"] = errMap
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
			return
		}

		// 从表单中构造用户数据
		user := models.User{}
		utils.StructAssign(&userForm, &user)

		// 密码加密处理
		encodePwd, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
		if err != nil {
			errMap["password"] = "密码加密错误"
			res["errMap"] = errMap
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
			return
		}
		user.Password = string(encodePwd)

		// 创建用户
		if err := models.SaveUser(&user); err != nil {
			logs.Error(fmt.Sprintf("save user failed: %#v, %s", user, err.Error()))
			res["errMsg"] = err.Error()
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
			return
		}

		this.Ctx.Output.SetStatus(201)
		_ = this.ServeJSON()
	}
}

// 编辑用户
func (this *UserController) Update() {
	this.Layout = "user/user-modal.html"
	this.TplName = "user/update.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Script"] = "user/update-script.html"

	// 获取用户ID，并获取该ID对应的用户
	id, _ := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	user, err := models.GetUserById(id)
	if err != nil {
		logs.Error(fmt.Sprintf("get user by id(%d) failed: %s", id, err.Error()))
		if this.Ctx.Input.Method() == "GET" {
			this.Data["errMsg"] = err.Error()
		} else {
			this.Data["json"] = map[string]string{"errMsg": err.Error()}
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
		}
		return
	}

	// GET请求获取用户编辑表单
	if this.Ctx.Input.Method() == "GET" {
		user.Password = "" // 去除密码
		// 以用户数据填充用户编辑表单
		this.Data["userToEdit"] = user
		// 获取XSRF Token
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	} else if this.Ctx.Input.Method() == "POST" {
		// 初始化响应
		res := make(map[string]interface{})
		this.Data["json"] = res

		// 解析用户编辑表单
		userForm := form.UserForm{}
		if err := this.ParseForm(&userForm); err != nil {
			res["errMsg"] = "提交的数据不合法"
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
			return
		}

		// 构造表单校验器，并进行检验
		var uValidator form.UserUpdatedValidator
		uValidator = &form.UserUpdateWithOutPasswordValidator{
			Username: userForm.Username,
		}
		if userForm.Password != "" {
			uValidator = &form.UserUpdateWithPasswordValidator{
				Username: userForm.Username,
				Password: userForm.Password,
			}
		}
		errMap := uValidator.Valid()
		if len(errMap) != 0 {
			res["errMap"] = errMap
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
			return
		}

		// 从表单中构造用户数据
		userInfo := models.User{}
		utils.StructAssign(&userForm, &userInfo)

		// 如更改了用户名，则需要确认更改后的用户名在数据库中不存在
		if user.Username != userInfo.Username {
			_, err := models.GetUserByUsername(userInfo.Username)
			if err == nil {
				res["errMsg"] = fmt.Sprintf("已存在相同用户名(%s)的用户", userInfo.Username)
				this.Ctx.Output.SetStatus(400)
				_ = this.ServeJSON()
				return
			}
		}

		userInfo.Id = id
		// 密码更新
		if userInfo.Password != "" {
			// 加密更新的密码
			encodePwd, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
			if err != nil {
				errMap["password"] = "密码加密错误"
				res["errMap"] = errMap
				this.Ctx.Output.SetStatus(400)
				_ = this.ServeJSON()
				return
			}
			userInfo.Password = string(encodePwd)
		}

		// 保存用户
		if err := models.SaveUser(&userInfo); err != nil {
			logs.Error(fmt.Sprintf("save user failed: %#v, %s", user, err.Error()))
			res["errMsg"] = err.Error()
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
			return
		}

		_ = this.ServeJSON()
	}
}

// 删除用户
func (this *UserController) Delete() {
	this.Layout = "user/user-modal.html"
	this.TplName = "user/delete.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Script"] = "user/delete-script.html"

	// 获取用户ID
	id, _ := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)

	// GET请求获取用户删除确认提示
	if this.Ctx.Input.Method() == "GET" {
		// 获取用户ID对应的用户
		user, err := models.GetUserById(id)
		if err != nil {
			logs.Error(fmt.Sprintf("get user by id(%d) failed: %s", id, err.Error()))
			this.Data["errMsg"] = err.Error()
			return
		}
		// 生成XSRF Token
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.Data["user"] = user
	} else {
		res := make(map[string]interface{})
		this.Data["json"] = res
		// 删除用户ID对应的用户
		if err := models.DeleteUser(id); err != nil {
			logs.Error(fmt.Sprintf("delete user(%d) failed: %s", id, err.Error()))
			res["errMsg"] = err.Error()
			this.Ctx.Output.SetStatus(400)
			_ = this.ServeJSON()
		}

		this.Ctx.Output.SetStatus(204)
		_ = this.ServeJSON()
	}
}
