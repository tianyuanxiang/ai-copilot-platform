package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	DB struct {
		DataSource string
	}

	BizRedis struct {
		Host string
		Pass string
		DB   int
	}

	Engine struct {
		BaseURL   string
		TimeoutMs int
	}

	Storage struct {
		UploadPath string
	}

	Elasticsearch ElasticsearchConf
}

type ElasticsearchConf struct {
	Addresses         []string
	Username          string
	Password          string
	KbChunksIndex     string
	SecurityLogsIndex string
}
