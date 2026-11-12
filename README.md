# GO语言博客
## 框架
> 框架：gin,go-spring,echo
> 
> golang 1.27
> 数据库：postgreSql 18/mysql/sqlite
> 缓存：bigcache


管理后台的前端代码：[https://github.com/hongmengzhu/xianfu-blog-go-frontend](https://github.com/hongmengzhu/xianfu-blog-go-frontend "管理后台前端代码")

## 功能说明
- [x] 文章 增删该查
- [x] 图片 上传
- [x] markdown 编辑器
- [x] 管理员 增删改查，密码修改
- [x] 每个租户独立域名
- [x] 菜单
- [x] 缓存
- [x] 标签
- [x] 附件
- [x] 角色和权限
- [x] 省市区
- [x] 博客前台显示
- [x] 博客网摘
- [x] 私有书签
- [x] 组织书签
- [ ] 未来功能
  - [ ] 站点属性配置
  - [ ] 七牛云存储


## 编译
### GO环境变量
根据你自己目录设置
```bash
export GOROOT=/usr/local/go
export GOBIN=$GOROOT/bin
export GOPROXY=https://goproxy.cn,direct
export GIT_SSL_NO_VERIFY=true
export PATH=.:$PATH:$GOBIN
```
### 拉取代码
```bash
git clone https://github.com/hongmengzhu/xianfu-blog-go
```
### 进入目录
```go
cd xianfu-blog-go
```
### 下载项目依赖
整理依赖（下载缺失的依赖 + 删除未使用的依赖）
```bash
go mod tidy
```
### 不编译，直接运行
```bash
go run .
```
### 编译打包
编译 可以在 linux 系统中运行的程序
#### 方式一
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o blogGo .
```
#### 方式二
*移除 移除编译路径信息*
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -ldflags '-s -w -extldflags "-static" -X "main.UserName="' \
-gcflags="all=-trimpath=${PWD}" \
-asmflags="all=-trimpath=${PWD}" \
-trimpath \
-o blogGo .
```
>注意
> `-s`：移除符号表。
> `-w`：移除调试信息。
> `-trimpath`：移除路径信息。
> 
### 其他系统编译
#### 编译成 windows 64 可执行文件
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o blogGo .
```

#### 编译成 macOS 64 可执行文件
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o blogGo .
```

#### 编译成 macOS Apple Silicon 可执行文件
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o blogGo .
```

编译命令解释
- GOOS=linux: 指定目标操作系统为 Linux
- GOARCH=amd64: 指定目标架构为 64 位
- CGO_ENABLED=0: 禁用 CGO（静态编译时需要）
- -a: 强制重新编译所有依赖包
- -ldflags '-extldflags "-static"': 设置链接器标志，确保静态链接

## 编译完成后操作 (案例)
新建文件夹 `dist`
1. 复制 `data`目录，到 `dist`目录
2. 复制 `blogGo`目录，到 `dist`目录

把 `dist`目录上传到服务器，然后执行 `blogGo`文件
```bash
# 在服务器上 含有 blogGo 目录下运行
./blogGo
```

## 升级依赖
```bash

go get -u ./

go mod tidy
```


# 账号
## manage 账号
```bash
账号： manage
密码： foxwho.com
```


## system 账号
密码 在 程序同级目录 data文件夹下，account-xxxxxx.md 文件中
按照最新日期，查看密码   （如果已经更改过密码了，那么这里无法使用）
```bash
账号： system
密码： 
```

#  项目结构
```bash
.
├── app
│   └── event                                   事件监听
│   └── manage                                  后台管理
│       ├── domainApi                           api接口权限控制
│       ├── domainBasic                         基础信息，地区，附件，国家，标签等等
│       ├── domainBlog                          博客文章，分类
│       │   ├── controller                          控制器
│       │   ├── model                               模型实体
│       │   │   └── modUser                         用户模型实体
│       │   └── service                             服务
│       ├── domainRam                           后台账号，权限相关
│       └── domainTc                            租户相关
│   └── system                                  系统管理
│       ├── basic                               基础信息，地区，附件，国家，标签等等
│       ├── ram                                 后台账号，权限相关
│       └── tc                                  租户相关
│   └── web                                     前台用户访问和对外api
│       ├── api                                对外api
│       ├── blog                               博客首页
│       │   ├── controller                          控制器
│       │   ├── model                               模型实体
│       │   │   └── modUser                         用户模型实体
│       │   └── service                             服务
│       └── utils                               模块工具类
├── assets                                      web 静态资源
│   └── img                                     图片
│   └── static                                  博客静态文件相关
├── cmd                                         命令文件目录，启动入口文件目录(多项目分别启动时使用)
│   ├── admin                                   admin项目
│   ├── api                                     api项目
│   └── moreServer                              多服务端口 启动
├── data                                        数据文件目录
│   ├── attachment                              上传文件夹目录
│   ├── config                                  配置文件目录
│   └── templates                               模版文件
├── doc                                         文档说明
├── infrastructure                              实体类，结构体，数据库表映射，资源映射
│   ├── entityDemo                              数据库表映射
│   └── repositoryDemo                          资源映射
├── middleware                                  中间件
│   ├── authPg                                    权限认证
│   ├── cachePg                                   缓存
│   ├── components                                其他组件相关
│   ├── dbPg                                      数据库
│   │   └── postgresqlPg
│   ├── runnerPg                                  应用启动后立即执行的一次性任务（初始化等）
│   ├── serverPg                                  服务端，多服务端
│   │   ├── ginServer                           gin 服务端
│   │   └── httpServer                          自定义 http 服务端
│   └── validatorPg                               验证器
├── pkg                                         可以被其他应用导入
│   ├── configPg                                配置结构映射
│   └── tools                                   工具类
│       └── wrapperPg
│           └── rg
├── router                                      路由
└── scripts                                     构建、安装、分析等不同功能的脚本文件
└── test                                        测试目录


```

