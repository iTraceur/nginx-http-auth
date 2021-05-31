package main

import (
	"flag"
	"fmt"
	"os"

	_ "nginx-http-auth/routers"
	_ "github.com/beego/beego/v2/server/web/session/redis"
	beego "github.com/beego/beego/v2/server/web"

	"nginx-http-auth/g"
)

func main() {
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
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

