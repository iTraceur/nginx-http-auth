# nginx-http-auth

An authentication backend for the Nginx `ngx_http_auth_request_module` module.

[中文](https://github.com/iTraceur/nginx-http-auth/blob/main/README_CN.md)

## usage
### Build from source
```bash
# git clone https://github.com/iTraceur/nginx-http-auth.git
# go get ./...
# chmod +x control
# ./control build
# ./control pack
```

### Install
```bash
# tar -zxvf nginx-http-auth-0.1.0.tar.gz 
# cp conf/nginx.conf /etc/nginx/nginx.conf  # and change it to suit your needs
# mv conf/app.example.conf conf/app.conf  # and change it to suit your needs 
# service nginx reload
# ./control start
```

### Configuration instructions
```ini
# Basic configuration
appname = nginx-http-auth  # App name
httpaddr = 127.0.0.1  # HTTP listening address
httpport = 8080  # HTTP listening port
runmode = dev  # Use "prod" for production environment or "test" for test environment

# Session configuration
sessionname = SessionID  # The name of the cookie that stored on the client
sessiongcmaxlifetime = 86400  # Session expiration time
sessioncookielifetime = 86400  # Cookie expiration time
sessionprovider = redis  # The provider of the session, you can also use memory, file, mysql, etc
sessionproviderconfig = "127.0.0.1:6379"  # The provider's path or link address


# XSRF configuration
xsrfkey = 4b6774f328ee1a2f24fcb62842fc0cfc  # XSRF key
xsrfexpire = 86400  # XSRF expiration time

# User remote authentication API
authAPI = http://127.0.0.1:5000/api/login

# The users who can access control API, default admin
controlUsers = admin;iTraceur;zhaowencheng

# Clinet IP control configuration
[ipControl]
direct = 127.0.0.1;192.168.1.5  # IPs that are allowed to access
deny =  # IPs that are denied to access, This configuration takes precedence over the direct

# Time control configuration
[timeControl]
direct = 09:00-21:00  # Time range that are denied to access
deny =  # Time range that are denied to access, This configuration takes precedence over the direct

# User control configuration
[userControl]
allow =    # Users that are allowed to access, This configuration takes precedence over the deny
deny = test;demo  # Users that are denied to access
```

### Nginx configuration
```bash
# cp conf/nginx.example.conf /etc/nginx/conf.d/nginx-http-auth.conf  # and change it to suit your needs
# service nginx reload
```
