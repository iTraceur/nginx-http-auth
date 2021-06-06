package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/redis"

	_ "nginx-http-auth/form"
	_ "nginx-http-auth/routers"
	_ "nginx-http-auth/utils"

	"nginx-http-auth/g"
	"nginx-http-auth/models"
)

func init() {
	_ = orm.RegisterDriver("sqlite3", orm.DRSqlite)
	_ = orm.RegisterDataBase("default", "sqlite3", "./data.db")
	_ = orm.RunSyncdb("default", false, true)
}

func main() {
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if beego.BConfig.RunMode == beego.DEV {
		orm.Debug = true
	}

	if authProvider, err := beego.AppConfig.String("authProvider"); err != nil || authProvider == "local" {
		models.CreateAdminUser()
	}

	// 启用Session和XSRF过滤
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "SessionID"
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600
	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = "272ae48b83413ca9982db969e22f1ece"
	beego.BConfig.WebConfig.XSRFExpire = 3600

	beego.Run()
}

