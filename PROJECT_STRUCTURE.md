# Go-Blog 后端项目结构说明

## 一、项目概述

本项目是一个基于 Go 的个人博客后端服务，采用经典的分层架构（接口层 → 服务层 → 数据层），使用 Gin 作为 Web 框架，MySQL 存储关系数据，Elasticsearch 存储文章内容并支持全文检索，Redis 作为缓存与热点数据存储。

## 二、技术栈

| 类别 | 技术 | 用途 |
| ---- | ---- | ---- |
| Web 框架 | gin-gonic/gin | HTTP 路由与请求处理 |
| ORM | gorm.io/gorm + mysql 驱动 | 操作 MySQL 关系数据库 |
| 搜索引擎 | go-elasticsearch/v8 | 操作 Elasticsearch，存储文章、支持全文检索 |
| 缓存 | go-redis/v9 | 缓存、登录会话、热点新闻缓存 |
| 认证 | golang-jwt/jwt/v4 | JWT 登录认证与管理员鉴权 |
| 定时任务 | robfig/cron/v3 | 定时同步文章浏览量等任务 |
| 日志 | go.uber.org/zap + lumberjack | 结构化日志输出与滚动切割 |
| 验证码 | mojocn/base64Captcha | 图形验证码 |
| 邮件 | jordan-wright/email | 发送邮件（注册验证等） |
| 第三方 | qiniu/go-sdk、高德 API、QQ 登录、hotSearch | 七牛云存储、高德地图、QQ 登录、热点新闻抓取 |
| CLI | urfave/cli | 命令行参数（初始化库、导入导出数据等） |

## 三、整体目录结构

```
go-blog/
├── server/                     # 后端服务主目录
│   ├── main.go                 # 程序入口
│   ├── config.yaml             # 配置文件
│   ├── go.mod                  # Go 模块定义
│   ├── api/                    # 接口层（HTTP 处理函数）
│   ├── assets/                 # 静态资源
│   ├── config/                 # 配置结构体定义
│   ├── core/                   # 核心初始化（配置、日志、服务器启动）
│   ├── flag/                   # 命令行工具（初始化数据库、导入导出）
│   ├── global/                 # 全局变量与模型
│   ├── initialize/             # 各组件初始化（路由、数据库、ES、Redis、定时任务）
│   ├── log/                    # 日志输出目录
│   ├── middleware/             # 中间件（JWT、管理员、日志、登录记录）
│   ├── model/                  # 数据模型（数据库表、ES 文档、请求/响应结构体）
│   ├── router/                 # 路由层（路由注册）
│   ├── service/                # 服务层（业务逻辑）
│   ├── task/                   # 定时任务
│   ├── uploads/                # 本地图片上传目录
│   └── utils/                  # 工具函数
└── temp/                       # 临时测试代码
```

## 四、分层架构说明

请求链路：`HTTP 请求 → middleware 中间件 → router 路由 → api 接口层 → service 服务层 → model 数据层（MySQL / Elasticsearch / Redis）`

### 1. api 接口层

> 目录：`server/api/`

只负责接收请求参数、调用服务层方法、返回统一格式响应，不包含业务逻辑。

- 通过 `enter.go` 中的 `ApiGroup` 聚合所有 API 结构体，各 API 通过 `service.ServiceGroupApp.XXX` 调用服务层
- 统一响应格式：`response.OkWithData(data, c)` 包装为 `{code, data, msg}` JSON

| 文件 | 职责 |
| ---- | ---- |
| base.go | 基础接口（验证码、发送邮箱验证码等） |
| user.go | 用户接口（注册、登录、修改密码等） |
| image.go | 图片接口（上传、列表、分类管理等） |
| article.go | 文章接口（增删改查、点赞、浏览量等） |
| comment.go | 评论接口（发布、回复、删除等） |
| feedback.go | 反馈接口 |
| advertisement.go | 广告位接口 |
| friend_link.go | 友链接口 |
| website.go | 网站配置接口（轮播图、热点新闻、日历等） |
| config.go | 系统配置接口 |

### 2. service 服务层

> 目录：`server/service/`

承载核心业务逻辑，负责跨库查询组装、数据加工、调用第三方服务。

| 文件 | 职责 |
| ---- | ---- |
| base.go | 基础服务（验证码、邮箱等） |
| user.go | 用户注册、登录、信息维护 |
| jwt.go | JWT 签发、校验、黑名单管理 |
| qq.go | QQ 第三方登录 |
| gaode.go | 高德地图 API 调用 |
| image.go | 图片上传、分类标记、分类计数 |
| article.go | 文章增删改查（MySQL + ES 两段式查询） |
| article_helper.go | 文章查询辅助（列表组装、字段减载） |
| article_stat.go | 文章统计（浏览量、点赞量更新） |
| comment.go / comment_helpers.go | 评论与评论树递归组装 |
| feedback.go | 反馈业务 |
| advertisement.go | 广告位业务 |
| friend_link.go | 友链业务 |
| website.go | 网站配置、轮播图、热点新闻 |
| hot_search.go | 热点新闻获取（Redis 缓存 + 策略模式） |
| calendar.go | 日历业务 |
| es_index.go | ES 索引管理 |
| config.go | 系统配置业务 |

### 3. router 路由层

> 目录：`server/router/`

注册路由，按权限等级将接口挂载到三个路由组：

| 路由组 | 权限要求 | 说明 |
| ---- | ---- | ---- |
| publicGroup | 无 | 公开接口，无需登录 |
| privateGroup | JWT 登录校验 | 需登录后访问 |
| adminGroup | JWT + 管理员角色双重校验 | 仅管理员可访问 |

### 4. middleware 中间件

> 目录：`server/middleware/`

| 文件 | 职责 |
| ---- | ---- |
| jwt.go | JWT 登录认证，校验 Token 并注入用户信息 |
| admin.go | 管理员角色鉴权 |
| logger.go | Gin 请求日志（GinLogger / GinRecovery） |
| login_record.go | 登录记录中间件 |

### 5. model 数据模型层

> 目录：`server/model/`

| 子目录 | 职责 |
| ---- | ---- |
| database/ | MySQL 表结构（user、image、comment、article_tag、feedback、advertisement、friend_link、login、jwt_blacklist、article_like、article_categor、footer_link 等） |
| elasticsearch/ | ES 索引文档结构（article 等） |
| request/ | 接口请求参数结构体 |
| response/ | 接口响应结构体（统一响应格式定义） |
| appTypes/ | 业务常量（图片分类、用户角色、存储方式、注册方式等） |
| other/ | 其他结构体（分页、热点新闻、日历、QQ、高德、ES 索引等） |

### 6. initialize 初始化

> 目录：`server/initialize/`

| 文件 | 职责 |
| ---- | ---- |
| router.go | 路由初始化（挂载静态文件、创建三个路由组） |
| gorm.go | MySQL 连接初始化 |
| es.go | Elasticsearch 客户端连接初始化 |
| redis.go | Redis 连接初始化 |
| corn.go | 定时任务初始化（注册 task 包任务） |
| other.go | 其他组件初始化（本地缓存等） |

### 7. core 核心层

> 目录：`server/core/`

| 文件 | 职责 |
| ---- | ---- |
| conf.go | 读取 config.yaml 初始化配置 |
| zap.go | 初始化 zap 日志 |
| server.go | 服务器启动入口（组装路由、加载 JWT 黑名单缓存、监听端口） |
| server_win.go / server_other.go | 按操作系统区分服务器启动方式（endless 平滑重启） |

### 8. global 全局变量

> 目录：`server/global/`

通过 `global.go` 暴露全局单例：`Config`（配置）、`Log`（日志）、`DB`（MySQL）、`EsClient`（ES）、`Redis`、`BlackCache`（JWT 黑名单本地缓存）。

### 9. flag 命令行工具

> 目录：`server/flag/`

基于 urfave/cli 实现的命令行工具，用于数据维护：

| 命令 | 功能 |
| ---- | ---- |
| `--sql` | 初始化 MySQL 数据库结构 |
| `--sql-export` | 导出 SQL 数据到文件 |
| `--sql-import` | 从文件导入 SQL 数据 |
| `--es` | 初始化 Elasticsearch 索引 |
| `--es-export` | 从 ES 导出数据 |
| `--es-import` | 向 ES 导入数据 |
| `--admin` | 按 config.yaml 创建管理员账号 |

### 10. task 定时任务

> 目录：`server/task/`

| 文件 | 任务 | 说明 |
| ---- | ---- | ---- |
| enter.go | 任务注册入口 | 注册定时触发器 |
| article_views.go | 同步文章浏览量 | 每小时将 Redis 中的浏览量增量同步到数据库 |
| hot_search.go | 热点新闻更新 | 定时抓取各平台热点新闻 |
| calendar.go | 日历数据 | 日历相关定时任务 |

### 11. utils 工具函数

> 目录：`server/utils/`

| 文件 / 子目录 | 职责 |
| ---- | ---- |
| upload/ | 图片上传实现（local.go 本地存储、qiniu.go 七牛云、upload.go 入口） |
| hotSearch/ | 热点新闻来源抓取（百度、头条、知乎、快手），通过策略模式分发 |
| image.go / illustration.go | 图片分类标记（封面 / 正文插图 / 未使用） |
| paginaton.go | 分页工具 |
| jwt.go / claim.go | JWT 解析与用户信息提取 |
| hash.go | 哈希工具 |
| email.go | 邮件发送工具 |
| random.go | 随机数生成 |
| http.go | HTTP 请求工具 |
| statistics.go | 统计工具 |
| diff_arrays.go | 数组差集计算（用于分类计数增减） |
| yaml.go / pares.go | 配置解析辅助 |

### 12. enter.go 聚合结构体 + 全局单例（核心设计）

每层目录下的 `enter.go` 采用「聚合结构体 + 全局单例」模式：把该层所有组件结构体内嵌（组合）进一个总结构体，再暴露一个包级单例变量。

```
router/enter.go     type RouterGroup struct { FriendLinkRouter ... }   var RouterGroupApp = new(RouterGroup)
api/enter.go        type ApiGroup struct { FriendLinkApi ... }         var ApiGroupApp = new(ApiGroup)
                    + 12 个 service 引用变量：var friendLinkService = service.ServiceGroupApp.FriendLinkService
service/enter.go    type ServiceGroup struct { FriendLinkService ... } var ServiceGroupApp = new(ServiceGroup)
```

层间连接方式（代码中的实际取值）：

| 连接 | 代码 |
| ---- | ---- |
| 路由层 → 接口层 | `friendLinkApi := api.ApiGroupApp.FriendLinkApi`（在 router 方法内） |
| 接口层 → 服务层 | `var friendLinkService = service.ServiceGroupApp.FriendLinkService`（在 api/enter.go 文件级变量） |

新增一个功能只需三步：新建 `XxxApi / XxxService / XxxRouter` 结构体 → 在对应 `enter.go` 中嵌入一行 → 在 `api/enter.go` 增加一行 service 引用变量，无需改动任何调用点。

## 五、程序启动流程

```
main.go
  ├── core.InitConf()           # 1. 读取 config.yaml
  ├── core.InitLogger()         # 2. 初始化 zap 日志
  ├── initialize.OtherInit()    # 3. 其他组件（本地缓存）
  ├── initialize.InitGorm()     # 4. 连接 MySQL
  ├── initialize.ConnectES()    # 5. 连接 Elasticsearch
  ├── initialize.ConnectRedis() # 6. 连接 Redis
  ├── flag.InitFlag()           # 7. 处理命令行参数（无参数则跳过）
  ├── initialize.InitCron()     # 8. 启动定时任务
  └── core.RunServer()          # 9. 初始化路由并启动 HTTP 服务器
```

## 六、关键设计约定

- **跨库两段式查询**：关系数据（用户点赞、分类等）存 MySQL，文章内容存 ES；列表接口先通过 MySQL 分页获取目标文档 ID，再批量从 ES 拉取内容组装返回，并对大字段（正文、关键词）进行减载
- **统一响应格式**：所有接口通过 `response.OkWithData / OkWithMessage` 返回 `{code, data, msg}` 结构
- **enter.go 聚合结构体 + 全局单例**：router / api / service 每层通过 `enter.go` 将所有组件嵌入一个总结构体，并暴露 `XxxGroupApp` 全局单例；api 层另持有 12 个 service 引用变量完成层间解耦，新增功能仅需新建结构体 + 嵌入 + 引用三步
- **权限路由组**：同一功能的不同接口按权限要求分别挂载到 publicGroup / privateGroup / adminGroup
- **高并发计数**：点赞、浏览量更新使用 ES Painless 脚本 `+=` 相对累加，避免并发丢失
- **热点新闻**：接口 + 策略模式，`hotSearch.NewSource` 按来源标识 switch 分发到对应抓取实现
- **静态文件**：`Router.StaticFS` 将磁盘 uploads 目录映射到 `/uploads` URL 前缀，供前端访问本地存储图片
- **安全**：数据库密码、JWT 密钥等敏感配置通过环境变量注入，不写入代码或配置文件
