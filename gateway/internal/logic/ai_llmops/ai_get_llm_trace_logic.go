// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_llmops

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGetLlmTraceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetLlmTraceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetLlmTraceLogic {
	return &AiGetLlmTraceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetLlmTraceLogic) AiGetLlmTrace(req *types.AiLlmTracePathReq) (resp *types.AiLlmTraceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
