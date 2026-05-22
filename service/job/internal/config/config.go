package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcClientConf

	SysRpc zrpc.RpcClientConf

	BizRedis BizRedis

	Jobs JobsConf
}

type JobsConf struct {
	ClearLoginLog ClearLoginLogJobConf
	ClearOperLog  ClearOperLogJobConf
}

type ClearLoginLogJobConf struct {
	Enable            bool
	Cron              string
	LockExpireSeconds int
}

type ClearOperLogJobConf struct {
	Enable            bool
	Cron              string
	LockExpireSeconds int
}

// Redis Redis缓存配置
type BizRedis struct {
	// Host Redis地址，格式：host:port
	Host string
	// Pass Redis密码，无密码则留空
	Pass string
	// DB Redis数据库编号，默认0
	DB int
}
