// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_llmops

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListLlmCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListLlmCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListLlmCallLogic {
	return &AiListLlmCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListLlmCallLogic) AiListLlmCall(req *types.AiListLlmCallReq) (resp *types.AiListLlmCallResp, err error) {
	// todo: add your logic here and delete this line

	return
}
