package form

import "time"

type UserForm struct {
	Id        int64     `form:"-"`
	Name      string    `form:"name"`
	Username  string    `form:"username"`
	Password  string    `form:"password"`
	ClientIp  string    `form:"client_ip"`
	Active    bool      `form:"active"`
	Temporary bool      `form:"temporary"`
	Created   time.Time `form:"-"`
	Updated   time.Time `form:"-"`
}
