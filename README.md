# nginx-http-auth

An authentication backend for the Nginx `ngx_http_auth_request_module` module.

[中文](https://github.com/iTraceur/nginx-http-auth/blob/main/README_CN.md)
 | 
[Demo](https://auth-demo.itraceur.com/), user/password: admin/auth-demo, this user has administrator privileges and can operate at will, and the data will be recovered every hour.

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

# User authentication provider, local or remote, default local.
authProvider = local

# Whether to enable User-IP binding to restrict the user to login to the application using a specific IP,
# takes effect when the authProvider is configured to local. once enabled, the manage user needs to
# bind the client IP for each user in the user management page, default false.
ipBinding = false

# User remote authentication API, required when the authProvider is configured to remote.
authAPI = http://127.0.0.1:5000/api/login

# The users who can access user manager page and control API, default admin.
manageUsers = admin;iTraceur;zhaowencheng

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

## Screenshots
### Auth login
![auth login][auth-login]

### Auth captcha
> Captcha is required after login failure.

![auth captcha][auth-captcha]

### User list
![user lig][user-list]

### Add user
![add user][add-user]

### Edit user
![edit user][edit-user]

### Delete user
![delete user][delete-user]


[auth-login]: ./static/screenshot/auth-login.jpg
[auth-captcha]: ./static/screenshot/auth-captcha.jpg
[user-list]: ./static/screenshot/user-list.jpg
[add-user]: ./static/screenshot/add-user.jpg
[edit-user]: ./static/screenshot/edit-user.jpg
[delete-user]: ./static/screenshot/delete-user.jpg
