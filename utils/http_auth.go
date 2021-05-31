package utils

import (
	"strings"
	"time"
	"crypto/tls"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
)

func HttpAuth(username string, password string) bool {
	api, err := beego.AppConfig.String("authAPI")
	if err != nil {
		return false
	}

	req := httplib.Post(api)
	if strings.HasSuffix(api, "https:") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	req.Header("User-Agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	req.Param("username", username)
	req.Param("password", password)

	resp, err := req.SetTimeout(1*time.Second, 3*time.Second).Response()
	if err != nil {
		logs.Error(err)
		return false
	} else if resp.StatusCode >= 400 {
		var body []byte
		resp.Body.Read(body)
		logs.Error("Http auth failed: %d %s", resp.StatusCode, string(body))
		return false
	}
	return true
}
