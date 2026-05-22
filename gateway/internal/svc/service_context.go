// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	aistatusclient "ai-copilot-platform/ai-rpc/client/aistatusservice"
	"ai-copilot-platform/gateway/internal/config"
	"ai-copilot-platform/gateway/internal/ws"
	"go-zero-rpc/common/middleware"
	authclient "go-zero-rpc/sys-rpc/client/authservice"
	permclient "go-zero-rpc/sys-rpc/client/permissionservice"
	systemclient "go-zero-rpc/sys-rpc/client/systemservice"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	AuthRpc authclient.AuthService
	SysRpc  systemclient.SystemService
	PermRpc permclient.PermissionService

	AiStatusClient    aistatusclient.AiStatusService
	AiKnowledgeClient aiknowledgeclient.AiKnowledgeService

	WsHub *ws.Hub

	AuthMiddleware    rest.Middleware
	CasbinMiddleware  rest.Middleware
	OperLogMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Create downstream RPC clients used by gateway logic.
	sysCli := zrpc.MustNewClient(c.SysRpc)
	aiCli := zrpc.MustNewClient(c.AiRpc)

	authSvc := authclient.NewAuthService(sysCli)
	sysSvc := systemclient.NewSystemService(sysCli)
	permSvc := permclient.NewPermissionService(sysCli)

	aiStatusSvc := aistatusclient.NewAiStatusService(aiCli)
	aiKnowledgeSvc := aiknowledgeclient.NewAiKnowledgeService(aiCli)

	wsHub := ws.NewHub()
	go wsHub.Run()

	return &ServiceContext{
		Config:  c,
		AuthRpc: authSvc,
		SysRpc:  sysSvc,
		PermRpc: permSvc,

		AiStatusClient:    aiStatusSvc,
		AiKnowledgeClient: aiKnowledgeSvc,

		WsHub: wsHub,

		AuthMiddleware:   middleware.NewAuthMiddleware(permSvc).Handle(c.Auth.AccessSecret),
		CasbinMiddleware: middleware.NewCasbinMiddleware(permSvc).Handle,
	}
}
