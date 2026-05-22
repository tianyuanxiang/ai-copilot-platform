// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_agent

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListToolCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListToolCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListToolCallLogic {
	return &AiListToolCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListToolCallLogic) AiListToolCall(req *types.AiListToolCallReq) (resp *types.AiListToolCallResp, err error) {
	// todo: add your logic here and delete this line

	return
}
