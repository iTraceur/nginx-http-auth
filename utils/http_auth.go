package utils

import (
	"fmt"
	"strings"
	"time"

	"crypto/tls"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// 远程用户认证
func HttpAuth(username string, password string) bool {
	// 远程用户认证接口
	api, err := beego.AppConfig.String("authAPI")
	if err != nil {
		return false
	}

	// 请求远程接口进行认证
	req := httplib.Post(api)
	// 跳过证书验证
	if strings.HasSuffix(api, "https:") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	// 设置请求头
	req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	// 添加用户名和密码请求数据
	req.Param("username", username)
	req.Param("password", password)

	// 获取响应，1秒连接超时，3秒响应超时
	resp, err := req.SetTimeout(1*time.Second, 3*time.Second).Response()
	if err != nil {
		logs.Error(fmt.Sprintf("Http auth failed: %s", err.Error()))
		return false
	} else if resp.StatusCode >= 400 {
		var body []byte
		b := resp.Body
		if _, err := b.Read(body); err != nil {
			logs.Error(fmt.Sprintf("Http auth failed: %d %s", resp.StatusCode, string(body)))
		}
		return false
	}
	return true
}
