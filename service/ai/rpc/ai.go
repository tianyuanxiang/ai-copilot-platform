package main

import (
	"flag"
	"fmt"
	"time"

	"ai-copilot-platform/ai-rpc/internal/config"
	aiagentservice "ai-copilot-platform/ai-rpc/internal/server/aiagentservice"
	aichatservice "ai-copilot-platform/ai-rpc/internal/server/aichatservice"
	aiknowledgeservice "ai-copilot-platform/ai-rpc/internal/server/aiknowledgeservice"
	aiobservabilityservice "ai-copilot-platform/ai-rpc/internal/server/aiobservabilityservice"
	aireportservice "ai-copilot-platform/ai-rpc/internal/server/aireportservice"
	aisecurityservice "ai-copilot-platform/ai-rpc/internal/server/aisecurityservice"
	aistatusservice "ai-copilot-platform/ai-rpc/internal/server/aistatusservice"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/ai.yaml", "the config file")

func main() {
	flag.Parse()

	logx.DisableStat()
	// 将慢查询阈值设为 1ms，使所有 SQL 语句都被记录到日志中（仅开发调试用）
	sqlx.SetSlowThreshold(time.Millisecond)

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterAiAgentServiceServer(grpcServer, aiagentservice.NewAiAgentServiceServer(ctx))
		pb.RegisterAiChatServiceServer(grpcServer, aichatservice.NewAiChatServiceServer(ctx))
		pb.RegisterAiKnowledgeServiceServer(grpcServer, aiknowledgeservice.NewAiKnowledgeServiceServer(ctx))
		pb.RegisterAiObservabilityServiceServer(grpcServer, aiobservabilityservice.NewAiObservabilityServiceServer(ctx))
		pb.RegisterAiReportServiceServer(grpcServer, aireportservice.NewAiReportServiceServer(ctx))
		pb.RegisterAiSecurityServiceServer(grpcServer, aisecurityservice.NewAiSecurityServiceServer(ctx))
		pb.RegisterAiStatusServiceServer(grpcServer, aistatusservice.NewAiStatusServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
