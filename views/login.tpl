<!DOCTYPE html>
<html >
<head>
    <meta charset="UTF-8">
    <title>认证登录</title>
    <link rel="stylesheet" href="/static/css/normalize.css">
    <link rel="stylesheet" href="/static/css/login.css">
</head>
<body>
<div class="login">
    <h1>认证登录</h1>
    <form action="/passport/login" method="post">
        {{if .msg}}
        <div class="msg">{{.msg}}</div>
        {{end}}
        <input type="text" name="username" placeholder="用户名" required="required" />
        <input type="password" name="password" placeholder="密码" required="required" />
        {{if .captcha}}
        <div class="row">
            <div class="captcha-input">
                <input name="captcha" type="text" placeholder="验证码" required>
            </div>
            <div class="captcha-img">
                {{create_captcha}}
            </div>
        </div>
        {{end}}
        {{ .xsrfdata }}
        <input type="hidden" name="target" value={{.target}}>
        <button type="submit" class="btn btn-primary btn-block btn-large">登&nbsp;录</button>
    </form>
</div>
</body>
</html>
