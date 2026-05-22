// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_public

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGetPublicKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetPublicKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetPublicKbLogic {
	return &AiGetPublicKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetPublicKbLogic) AiGetPublicKb(req *types.AiPublicKbPathReq) (resp *types.AiPublicKbItem, err error) {
	// todo: add your logic here and delete this line

	return
}
