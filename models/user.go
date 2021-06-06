package models

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func init() {
	orm.RegisterModel(new(User))
}

type User struct {
	Id        int64
	Name      string `orm:"null"`
	Username  string `orm:"index;unique;size(50);description(用户名)"`
	Password  string
	ClientIp  string    `orm:"null;size(15)"`
	Active    bool      `default:"true"`
	Temporary bool      `default:"false"`
	Created   time.Time `orm:"auto_now_add;type(datetime)"`
	Updated   time.Time `orm:"auto_now;type(datetime)"`
}

// 创建管理员
func CreateAdminUser() {
	user, err := GetUserByUsername("admin")
	if err == nil {
		return
	}
	encodePwd, _ := bcrypt.GenerateFromPassword([]byte("nginx-http-auth"), bcrypt.DefaultCost)
	user = &User{Username: "admin", Password: string(encodePwd), Active: true}
	o := orm.NewOrm()
	if _, err = o.Insert(user); err != nil {
		logs.Error(fmt.Sprintf("insert manage user(admin) failed: %s", err.Error()))
	}
	logs.Info("create manage user (admin) successed")
}

//创建、更新用户
func SaveUser(u *User) (err error) {
	o := orm.NewOrm()

	if u.Id == 0 {
		_, err = GetUserByUsername(u.Username)
		if err == nil {
			return errors.New("该用户已存在")
		}
		if _, err = o.Insert(u); err != nil {
			logs.Error(fmt.Sprintf("insert user failed: %#v, %s", u, err.Error()))
			return errors.New("用户创建失败")
		}
	} else {
		var user *User
		user, err = GetUserById(u.Id)
		if err != nil {
			return errors.New("该用户不存在")
		}

		user.Username = u.Username
		user.Name = u.Name
		user.ClientIp = u.ClientIp
		user.Active = u.Active
		user.Temporary = u.Temporary

		if u.Password != "" {
			user.Password = u.Password
		}

		if _, err = o.Update(user); err != nil {
			logs.Error(fmt.Sprintf("update user failed: %#v, %s", user, err.Error()))
			return errors.New("用户保存失败")
		}
	}
	return
}

// 删除用户
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	if err = o.Read(&v, "Id"); err == nil {
		if _, err := o.Delete(&User{Id: id}); err != nil {
			logs.Error(fmt.Sprintf("delete user (uid: %d) failed: %s", id, err.Error()))
		}
	}
	return
}

// 通过用户ID获取用户
func GetUserById(id int64) (user *User, err error) {
	o := orm.NewOrm()
	user = &User{Id: id}

	if err = o.Read(user); err != nil {
		logs.Error(fmt.Sprintf("get user by id(%d) failed: %s", id, err.Error()))
		return nil, errors.New("该用户不存在")
	}
	return
}

// 通过用户名获取用户
func GetUserByUsername(uname string) (user *User, err error) {
	o := orm.NewOrm()
	user = &User{Username: uname}

	if err = o.QueryTable("user").Filter("Username", user.Username).One(user); err != nil {
		logs.Error(fmt.Sprintf("get user by username(%s) failed: %s", uname, err.Error()))
		return nil, errors.New("该用户不存在")
	}
	return
}

// 查询所有用户
func QueryAllUser() (userList []interface{}, err error) {
	var list []User
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	if _, err = qs.All(&list); err == nil {
		for _, v := range list {
			userList = append(userList, v)
		}
		return userList, nil
	}
	return nil, err
}

// 获取分页用户列表
func LimitList(pageNo int, pageSize int) (users []User) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	qs = qs.OrderBy("-updated", "-created").Limit(pageSize, (pageNo-1)*pageSize)
	_, _ = qs.All(&users)
	return
}

// 获取用户总数
func GetUserCount() int64 {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))

	num, err := qs.Count()
	if err == nil {
		return num
	} else {
		return 0
	}
}
