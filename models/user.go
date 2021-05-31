package models

import (
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(User))
}

type User struct {
	Id          int
	Name        string
	Username	string
	Password	string
}
