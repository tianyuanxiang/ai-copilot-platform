// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_agent

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiRunAgentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiRunAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiRunAgentLogic {
	return &AiRunAgentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiRunAgentLogic) AiRunAgent(req *types.AiAgentRunReq) (resp *types.AiAgentRunResp, err error) {
	// todo: add your logic here and delete this line

	return
}
