package aiagentservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRunAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunAgentLogic {
	return &RunAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 运行受控 Agent，只允许调用白名单工具。
func (l *RunAgentLogic) RunAgent(in *pb.AgentRunReq) (*pb.AgentRunResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AgentRunResp{}, nil
}
