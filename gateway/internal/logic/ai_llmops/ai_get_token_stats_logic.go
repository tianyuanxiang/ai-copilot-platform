// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_llmops

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGetTokenStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetTokenStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetTokenStatsLogic {
	return &AiGetTokenStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetTokenStatsLogic) AiGetTokenStats(req *types.AiTokenStatsReq) (resp *types.AiTokenStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
