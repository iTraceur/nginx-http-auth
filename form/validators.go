package form

import (
	"strings"

	"github.com/beego/beego/v2/core/validation"
)

type UserAuthValidator struct {
	Username string `alias:"用户名" valid:"Required"`
	Password string `alias:"密码" valid:"Required"`
}

func (u *UserAuthValidator) Valid() map[string]string {
	valid := validation.Validation{}
	status, _ := valid.Valid(u)
	errMap := map[string]string{}
	if !status {
		for _, err := range valid.Errors {
			var alias = GetAlias(UserAuthValidator{}, err.Field)
			errMap[strings.ToLower(err.Field)] = alias + err.Message
		}
	}
	return errMap
}

type UserCreateValidator struct {
	Username string `alias:"用户名" valid:"Required;MinSize(5);MaxSize(25)"`
	Password string `alias:"密码" valid:"Required;MinSize(8);MaxSize(255)"`
}

func (u *UserCreateValidator) Valid() map[string]string {
	valid := validation.Validation{}
	status, _ := valid.Valid(u)
	errMap := map[string]string{}
	if !status {
		for _, err := range valid.Errors {
			var alias = GetAlias(UserCreateValidator{}, err.Field)
			errMap[strings.ToLower(err.Field)] = alias + err.Message
		}
	}
	return errMap
}

type UserUpdatedValidator interface {
	Valid() map[string]string
}

type UserUpdateWithPasswordValidator struct {
	Username string `alias:"用户名" valid:"Required;MinSize(5);MaxSize(25)"`
	Password string `alias:"密码" valid:"MinSize(8);MaxSize(255)"`
}

func (u *UserUpdateWithPasswordValidator) Valid() map[string]string {
	valid := validation.Validation{}
	status, _ := valid.Valid(u)
	errMap := map[string]string{}
	if !status {
		for _, err := range valid.Errors {
			var alias = GetAlias(UserUpdateWithPasswordValidator{}, err.Field)
			errMap[strings.ToLower(err.Field)] = alias + err.Message
		}
	}
	return errMap
}

type UserUpdateWithOutPasswordValidator struct {
	Username string `alias:"用户名" valid:"Required;MinSize(5);MaxSize(25)"`
}

func (u *UserUpdateWithOutPasswordValidator) Valid() map[string]string {
	valid := validation.Validation{}
	status, _ := valid.Valid(u)
	errMap := map[string]string{}
	if !status {
		for _, err := range valid.Errors {
			var alias = GetAlias(UserUpdateWithOutPasswordValidator{}, err.Field)
			errMap[strings.ToLower(err.Field)] = alias + err.Message
		}
	}
	return errMap
}
