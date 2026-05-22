package svc

import (
	"ai-copilot-platform/job/internal/config"
	"ai-copilot-platform/sys-rpc/client/systemservice"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	SysRpc systemservice.SystemService
	RDB    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.BizRedis.Host,
		Password: c.BizRedis.Pass,
		DB:       c.BizRedis.DB,
	})

	return &ServiceContext{
		Config: c,
		SysRpc: systemservice.NewSystemService(
			zrpc.MustNewClient(c.SysRpc),
		),
		RDB: rdb,
	}
}
