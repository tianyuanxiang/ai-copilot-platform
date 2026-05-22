// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	SysRpc zrpc.RpcClientConf

	AiRpc zrpc.RpcClientConf

	Upload struct {
		Path        string
		MaxSize     int64
		AllowedExts []string
	}

	WebSocket struct {
		HeartbeatInterval int
		MaxConnections    int
	}
}

type AIEngineConf struct {
	BaseURL string
	Timeout int64
}
