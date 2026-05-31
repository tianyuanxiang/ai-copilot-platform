package aiwindagentservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/logic/aiwindtool"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteToolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExecuteToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteToolLogic {
	return &ExecuteToolLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExecuteToolLogic) ExecuteTool(in *pb.WindToolExecuteReq) (*pb.WindToolExecuteResp, error) {
	// ExecuteToolLogic 只负责接收 RPC 请求。
	// 工具白名单、参数校验、业务执行和审计日志统一交给 aiwindtool.Executor。
	executor := aiwindtool.NewExecutor(l.svcCtx)
	return executor.Execute(l.ctx, in)
}
