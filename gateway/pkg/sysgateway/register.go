package sysgateway

import (
	"ai-copilot-platform/gateway/internal/config"
	"ai-copilot-platform/gateway/internal/handler"
	"ai-copilot-platform/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Auth      AuthConfig
	SysRpc    zrpc.RpcClientConf
	Upload    UploadConfig
	WebSocket WebSocketConfig
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}

type UploadConfig struct {
	Path        string
	MaxSize     int64
	AllowedExts []string
}

type WebSocketConfig struct {
	HeartbeatInterval int
	MaxConnections    int
}

// 可由外部导入的注册包

func Register(server *rest.Server, c Config) {
	serverCtx := svc.NewServiceContext(toInternalConfig(c))
	handler.RegisterHandlers(server, serverCtx)
}

func toInternalConfig(c Config) config.Config {
	var internal config.Config
	internal.Auth.AccessSecret = c.Auth.AccessSecret
	internal.Auth.AccessExpire = c.Auth.AccessExpire
	internal.SysRpc = c.SysRpc
	internal.Upload.Path = c.Upload.Path
	internal.Upload.MaxSize = c.Upload.MaxSize
	internal.Upload.AllowedExts = c.Upload.AllowedExts
	internal.WebSocket.HeartbeatInterval = c.WebSocket.HeartbeatInterval
	internal.WebSocket.MaxConnections = c.WebSocket.MaxConnections
	return internal
}
