# nginx-http-auth

一个用于Nginx `ngx_http_auth_request_module` 模块的认证后端。

## 安装使用
### 从源码构建

```bash
# git clone https://github.com/iTraceur/nginx-http-auth.git
# go get ./...
# chmod +x control
# ./control build
# ./control pack
```

### 安装
```bash
# tar -zxvf nginx-http-auth-0.1.0.tar.gz 
# mv conf/app.example.conf conf/app.conf  # 按需要更改相应的配置
# ./control start
```

### 配置说明
```ini
# 基本配置
appname = nginx-http-auth
httpaddr = 127.0.0.1  # HTTP 监听地址，按需要进行更改
httpport = 8080  # HTTP 监听地址，按需要进行更改
runmode = dev  # 生产环境使用 prod，测试环境使用 test

# Session 配置
sessionname = SessionID  # 存储在客户端的 cookie 名称
sessiongcmaxlifetime = 86400  # Session 过期时间
sessioncookielifetime = 86400  # Cookie 过期时间
sessionprovider = redis  # Session 存储引擎, 还支持 memory, file, mysql 等
sessionproviderconfig = "127.0.0.1:6379"  # Session 存储引擎的路径或链接地址


# XSRF 配置
xsrfkey = 4b6774f328ee1a2f24fcb62842fc0cfc  # XSRF key
xsrfexpire = 86400  # XSRF 过期时间

# 远程用户认证接口
authAPI = http://127.0.0.1:5000/api/login

# 可访问控制接口的用户，默认为admin
controlUsers = admin;iTraceur;zhaowencheng

# 客户端 IP 访问控制
[ipControl]
direct = 127.0.0.1;192.168.1.5  # 允许访问的 IP
deny =  # 拒绝访问的IP，这个配置优先于 direct 的配置

# 访问时间段控制
[timeControl]
direct = 09:00-21:00  # 允许访问的时间段
deny =  # 拒绝访问的时间段，这个配置优先于 direct 的配置

# 用户访问控制
[userControl]
allow =    # 允许访问的用户，这个配置优先于 deny 的配置
deny = test;demo  # 拒绝访问的用户
```

### 配置Nginx
```bash
# cp conf/nginx.example.conf /etc/nginx/conf.d/nginx-http-auth.conf  # 按需要更改相应的配置
# service nginx reload
```
