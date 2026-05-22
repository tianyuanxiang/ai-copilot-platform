# go-zero-rpc 架构仓

`go-zero-rpc` 是基于 **go-zero v1.10** 的通用微服务架构仓，提供后台系统常用的认证、RBAC 权限、系统管理、文件管理、日志管理、WebSocket、定时任务等基础能力。

当前推荐用法不是把整套代码复制到业务项目里改，而是把架构仓作为独立依赖接入业务仓：

- 架构仓继续维护通用能力：`common`、`sys-rpc`、`gateway/pkg/sysgateway`。
- 业务仓只维护自己的业务 RPC 和业务 HTTP 路由。
- 业务网关通过一行注册，把架构仓的通用 HTTP 路由挂载到自己的 `rest.Server` 上。

## 架构演进思路

本项目经过了几轮拆分和优化：

1. **单体架构**
   - 系统能力和业务能力都放在同一个项目里。
   - 缺点：架构代码和业务代码耦合，多个项目复用困难。

2. **嵌入微服务**
   - 把认证、用户、角色、菜单、权限等系统能力拆成 `sys-rpc` 微服务。
   - 业务服务通过 RPC 调用系统能力，不再直接复制系统服务逻辑。

3. **抽取服务逻辑后遇到网关重复维护问题**
   - 服务层逻辑已经通用化，但业务网关仍需要复制架构仓的 `.api`、`handler`、`logic`、`types`。
   - 架构仓新增或修改系统接口后，业务仓还要同步修改网关层代码，维护成本仍然很高。

4. **当前优化：架构仓导出可嵌入网关模块**
   - 架构仓新增 `gateway/pkg/sysgateway`。
   - 业务网关只需要调用 `sysgateway.Register(server, config)`，即可挂载架构仓的 `/api/v1/auth` 和 `/api/v1/system` 路由。
   - 架构仓后续修改系统接口时，业务仓只需要升级依赖，不再复制或重写 sys/auth 的网关代码。

## 在自己的业务项目中推荐架构

```text
HTTP Client
    |
    v
业务仓 Gateway，一个 HTTP Server
    |
    |-- /api/v1/<biz>      -> 业务仓自己的 handler/logic 提供
    |-- /api/v1/ws         -> 可由业务仓按需提供
    |
    +-- 调用架构仓 sys-rpc
    +-- 调用业务仓业务 RPC

架构仓 sys-rpc
    |
    +-- PostgreSQL
    +-- Redis
    +-- Casbin
```

### 关键边界

- `/api/v1/auth` 和 `/api/v1/system` 归架构仓维护。
- 业务仓不再复制架构仓的 `auth/sys` 网关代码。
- 业务仓只维护自己的业务 `.api`、`handler`、`logic`、`types`。
- `gateway/pkg/sysgateway` 是对外公共包；架构仓的 `gateway/internal/*` 仍然保持内部实现，不暴露给业务仓。

## 目录结构

```text
go-zero-rpc/
├── common/                    # 通用能力：响应、错误、JWT、中间件等
├── gateway/                   # 架构仓自己的 HTTP 网关
│   ├── gateway.go             # 独立启动模式入口
│   ├── api/                   # 架构仓 auth/system API 定义
│   ├── internal/              # 架构仓网关内部 handler/logic/svc/types
│   └── pkg/sysgateway/        # 对外导出的可嵌入网关注册包
├── service/sys/rpc/           # 系统 RPC 服务：认证、权限、用户、角色、菜单等
├── service/job/               # 定时任务服务
├── deploy/                    # 数据库初始化脚本
├── etc/                       # 通用配置，例如 Casbin rbac_model.conf
└── go.work                    # 本地 workspace
```

## 两种运行模式

### 模式一：架构仓独立运行

用于单独开发、调试架构仓自身能力。

此时启动架构仓自己的 gateway：

```bash
# 终端 1：启动 sys-rpc
cd D:/GoProject/project_new/go-zero-rpc/service/sys/rpc
go run sys.go -f etc/sys.yaml

# 终端 2：启动 job，可选
cd D:/GoProject/project_new/go-zero-rpc/service/job
go run job.go -f etc/job.yaml

# 终端 3：启动架构仓 gateway
cd D:/GoProject/project_new/go-zero-rpc/gateway
go run gateway.go -f etc/gateway-api.yaml
```

此模式下，架构仓 gateway 自己创建 `rest.Server`，并调用：

```go
ctx := svc.NewServiceContext(c)
handler.RegisterHandlers(server, ctx)
server.Start()
```

### 模式二：业务仓嵌入架构仓网关模块（推荐）

用于业务项目接入架构仓通用能力。

此时业务仓启动自己的 gateway，不需要再启动架构仓 gateway。业务仓 gateway 会把架构仓路由和业务路由注册到同一个 `rest.Server` 上。

#### 新业务项目入口示例：

```go
package main

import (
    "new_project/gateway/internal/config"
    "new_project/gateway/internal/handler"
    "new_project/gateway/internal/svc"
    "go-zero-rpc/gateway/pkg/sysgateway"

    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/rest"
)

func main() {
    var c config.Config
    conf.MustLoad("etc/gateway-api.yaml", &c)

    server := rest.MustNewServer(c.RestConf)
    defer server.Stop()

    // 业务仓自己的上下文，用于业务路由。
    ctx := svc.NewServiceContext(c)

    // 挂载架构仓通用路由：/api/v1/auth 和 /api/v1/system。
    sysgateway.Register(server, sysgateway.Config{
        Auth: sysgateway.AuthConfig{
            AccessSecret: c.Auth.AccessSecret,
            AccessExpire: c.Auth.AccessExpire,
        },
        SysRpc: c.SysRpc,
        Upload: sysgateway.UploadConfig{
            Path:        c.Upload.Path,
            MaxSize:     c.Upload.MaxSize,
            AllowedExts: c.Upload.AllowedExts,
        },
        WebSocket: sysgateway.WebSocketConfig{
            HeartbeatInterval: c.WebSocket.HeartbeatInterval,
            MaxConnections:    c.WebSocket.MaxConnections,
        },
    })

    // 挂载业务仓自己的路由，例如 /api/v1/new_project。
    handler.RegisterHandlers(server, ctx)

    server.Start()
}
```

最终同一个业务 gateway 上会有两类路由：

```text
/api/v1/auth            -> 架构仓 sysgateway 注册
/api/v1/system          -> 架构仓 sysgateway 注册
/api/v1/new_project     -> 业务仓 handler.RegisterHandlers 注册
/api/v1/ws              -> 业务仓按需注册
```

## 对 sysgateway 的正确理解

`sysgateway.Register(server, config)` 做的事情是：

1. 根据外部传入的 `sysgateway.Config` 创建架构仓内部 `ServiceContext`。
2. 复用架构仓已有的 `internal/handler.RegisterHandlers`。
3. 把架构仓的 auth/system 路由注册到调用方传入的 `server` 上。

它不会启动新的 HTTP Server，也不会启动架构仓 `gateway.go`。

架构仓 `gateway.go` 和 `pkg/sysgateway` 是两种使用方式：

- `gateway.go`：架构仓独立运行模式。
- `pkg/sysgateway`：业务网关嵌入运行模式。

如果业务仓 gateway 已经通过 `sysgateway.Register` 挂载了架构仓路由，通常只需要启动业务仓 gateway 和架构仓 `sys-rpc`，不需要再启动架构仓 gateway。

如果确实同时启动业务仓 gateway 和架构仓 gateway，也不会发生代码层面的重复注册，因为它们是两个不同进程、两个不同 `rest.Server`。但它们不能监听同一个端口，而且业务上通常没有必要暴露两套相同的 auth/system HTTP 入口。

## 对sysgateway的通俗理解

- #### 架构仓的改进：

把架构仓网关的路由方法从`gateway.go`中包装一份到`pkg/sysgateway.Register `函数中，**Register**函数可由`自定义项目`引入，传入相关配置，例如：

```go
sysgateway.Register(server, sysgateway.Config{
		Auth: sysgateway.AuthConfig{
			AccessSecret: c.Auth.AccessSecret,
			AccessExpire: c.Auth.AccessExpire,
		},
		SysRpc: c.SysRpc,
		......
})
```

这样`自定义项目`就可以把架构仓的 `/api/v1/auth`、`/api/v1/system` 等通用路由注册到`自定义项目`的server 中。

但是`需要注意`：架构仓的gateway.go文件中还是有：

```go
ctx := svc.NewServiceContext(c) handler.RegisterHandlers(server, ctx)
```

启动业务仓的`gateway.go`后，路由注册会进来，如果这时再同时启动架构仓的`gateway.go`，路由会重复注册，所以启动业务仓的`gateway.go`后，架构仓就不能启动`gateway.go`了，除非换个端口。 

- #### 业务仓的改进

业务仓启动`gateway.go`的时候会先调用

```
"go-zero-rpc/gateway/pkg/sysgateway"
```

把定义好的的server和相关config通过`sysgateway.Register(server, config)`去架构仓中注册路由到自己的server中。拿到架构仓的路由后，再通过

```go
handler.RegisterHandlers(server, ctx)
server.Start() 
```

把`自定义项目`的路由再注册进去。

## 业务仓接入步骤

### 1. 添加模块依赖

业务网关 `go.mod` 引入架构仓模块：

```go
require (
    go-zero-rpc/common v1.0.0
    go-zero-rpc/gateway v1.0.0
    go-zero-rpc/sys-rpc v1.0.0
)

replace (
    go-zero-rpc/common => D:/xxx/xxx/go-zero-rpc/common
    go-zero-rpc/gateway => D:/xxx/xxx/go-zero-rpc/gateway
    go-zero-rpc/sys-rpc => D:/xxx/xxx/go-zero-rpc/service/sys/rpc
)
```

本地开发阶段可以使用 `replace` 指向本地架构仓目录；发布版本后，可改为 Git tag 或私有模块版本。

### 2. 业务网关配置

业务网关的配置需要包含业务自身配置，以及 sysgateway 运行所需字段。

示例：

```yaml
Name: gateway
Host: 0.0.0.0
Port: 8888
Mode: dev
Timeout: 30000

Auth:
  AccessSecret: "your-jwt-secret-key"
  AccessExpire: 1200

SysRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: sys.rpc

Upload:
  Path: uploads
  MaxSize: 10485760
  AllowedExts:
    - .png
    - .jpg
    - .jpeg
    - .pdf

WebSocket:
  HeartbeatInterval: 30
  MaxConnections: 10000
```

说明：

- `RestConf` 不需要传给 `sysgateway`，因为业务仓已经创建了 `rest.Server`。
- `Auth`、`SysRpc`、`Upload`、`WebSocket` 是架构仓通用路由运行所需配置。
- 业务仓配置不要求和架构仓 `internal/config.Config` 完全一一对应，只需要能映射到 `sysgateway.Config`。

### 3. 调整业务仓 API 定义

业务仓 `gateway.api` 不再 import 架构仓的 auth/system API：

```api
import (
    "desc/base.api"
    "desc/new_project/project.api"
    "desc/new_project/parameter_set.api"
    "desc/new_project/file.api"
    "desc/new_project/job.api"
)
```

不要再引入：

```api
"desc/auth.api"
"desc/sys/user.api"
"desc/sys/role.api"
"desc/sys/menu.api"
"desc/sys/api.api"
"desc/sys/casbin.api"
"desc/sys/dict.api"
"desc/sys/log.api"
"desc/sys/file.api"
```

否则业务仓会再次生成并维护 sys/auth 网关代码，回到重复维护的问题。

### 4. 调整业务仓路由注册

业务仓 `handler.RegisterHandlers` 只注册业务路由，例如：

```go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    server.AddRoutes(..., rest.WithPrefix("/api/v1/new_project"))
    server.AddRoutes(..., rest.WithPrefix("/api/v1")) // 例如 /ws
}
```

不要在业务仓里再注册：

```text
/api/v1/auth
/api/v1/system
```

这两个前缀由 `sysgateway.Register` 负责。

## 架构仓自身改造方案

### 已新增公共包

```text
gateway/pkg/sysgateway/register.go
```

公共接口：

```go
type Config struct {
    Auth      AuthConfig
    SysRpc    zrpc.RpcClientConf
    Upload    UploadConfig
    WebSocket WebSocketConfig
}

func Register(server *rest.Server, c Config)
```

设计原则：

- 公共包只暴露稳定的配置结构和注册函数。
- 不向业务仓暴露 `gateway/internal/svc`、`gateway/internal/handler`、`gateway/internal/types`。
- 内部继续复用架构仓 goctl 生成的 handler/logic/types。
- 架构仓接口变化时，只改架构仓；业务仓通过升级依赖获得新路由（使用本地目录**replace**的方式 不需要手动升级）。

### 架构仓独立网关保留

`gateway/gateway.go` 保留，用于架构仓单独运行和调试。

这不会和业务仓嵌入模式冲突，因为：

- 独立模式：架构仓自己创建 server，自己启动。
- 嵌入模式：业务仓创建 server，架构仓只把路由注册进去。

是否启动架构仓 gateway 是部署选择，不是代码依赖要求。

## 启动完整业务系统

推荐业务项目实际运行时启动：

```bash
# 1. 启动 etcd、PostgreSQL、Redis

# 2. 启动架构仓 sys-rpc
cd D:/xxx/xxx/go-zero-rpc/service/sys/rpc
go run sys.go -f etc/sys.yaml

# 3. 启动架构仓 job，可选
cd D:/xxx/xxx/go-zero-rpc/service/job
go run job.go -f etc/job.yaml

# 4. 启动业务仓自己的 RPC 服务
cd D:/xxx/xxx/new_project/backend/service/blade/rpc
go run new_project.go -f etc/blade.yaml

# 5. 启动业务仓 gateway
cd D:/xxx/xxx/new_project/backend/gateway
go run gateway.go -f etc/gateway-api.yaml
```

此时不需要再启动：

```bash
D:/xxx/xxx/go-zero-rpc/gateway/gateway.go
```

## 验证接口

登录：

```bash
curl -X POST http://localhost:8888/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

访问系统接口：

```bash
curl http://localhost:8888/api/v1/system/user \
  -H "Authorization: Bearer <access_token>"
```

访问业务接口：

```bash
curl http://localhost:8888/api/v1/blade/project \
  -H "Authorization: Bearer <access_token>"
```

## 配置和数据初始化

### 数据库初始化

```bash
createdb system
psql -d system -f deploy/system_pgsql.sql
```

### Casbin 管理员策略

```sql
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
VALUES ('p', 'admin', '/*', '*', '', '', '');
```

### sys-rpc 配置示例

```yaml
Name: sys.rpc
ListenOn: 0.0.0.0:9100
Mode: dev

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: sys.rpc

JwtAuth:
  AccessSecret: "your-jwt-secret-key"
  AccessExpire: 1200
  RefreshExpire: 10080

DB:
  DataSource: "postgres://user:password@127.0.0.1:5432/system?sslmode=disable&TimeZone=Asia/Shanghai"

CacheRedis:
  - Host: "127.0.0.1:6379"
    Pass: ""
    Type: node

BizRedis:
  Host: 127.0.0.1:6379
  Pass: ""
  DB: 10

CasbinModelPath: "etc/rbac_model.conf"
UploadPath: "uploads"
```

## 接口归属

| 路径前缀 | 归属 | 说明 |
|---|---|---|
| `/api/v1/auth` | 架构仓 `sysgateway` | 登录、刷新 Token、退出、当前用户、修改密码 |
| `/api/v1/system` | 架构仓 `sysgateway` | 用户、角色、菜单、API、字典、文件、日志、Casbin 规则 |
| `/api/v1/<biz>` | 业务仓 | 业务项目自己的接口 |
| `/api/v1/ws` | 业务仓或架构仓按需选择 | 当前业务项目可自行保留 |

## 后续维护约定

- 架构仓新增或修改 auth/system 接口时，只在架构仓中更新 `.api`、`handler`、`logic`、`types`。
- 业务仓不要手工复制架构仓 `gateway/internal` 下的代码。
- 业务仓升级架构仓能力时，更新 `go-zero-rpc/gateway`、`go-zero-rpc/common`、`go-zero-rpc/sys-rpc` 的依赖版本。
- 如果业务项目确实需要覆盖某个系统接口，优先在架构仓提供可配置开关或扩展点，不建议直接复制整套 sys/auth 网关代码。

## 常见问题

### 业务仓启动后，还要启动架构仓 gateway 吗？

通常不需要。业务仓 gateway 已经通过 `sysgateway.Register` 挂载了架构仓的 auth/system 路由。

### 同时启动业务仓 gateway 和架构仓 gateway 会重复注册吗？

不会发生同一个 server 内的重复注册，因为它们是两个进程、两个 `rest.Server`。但不能监听同一个端口，业务上也通常没必要同时暴露两套 HTTP 入口。

### 业务仓配置必须和架构仓 config.Config 完全一样吗？

不需要。业务仓只要能提供 `sysgateway.Config` 需要的字段即可：`Auth`、`SysRpc`、`Upload`、`WebSocket`。

### 为什么不直接导入架构仓 internal/handler？

Go 的 `internal` 包有访问边界，业务仓不应该也不能稳定依赖架构仓 `gateway/internal/*`。所以架构仓提供 `gateway/pkg/sysgateway` 作为公共边界。

## License

MIT License.
