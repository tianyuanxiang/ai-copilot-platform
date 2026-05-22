// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_personal

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGetPersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetPersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetPersonalKbLogic {
	return &AiGetPersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetPersonalKbLogic) AiGetPersonalKb(req *types.AiPersonalKbPathReq) (resp *types.AiPersonalKbItem, err error) {
	// todo: add your logic here and delete this line

	return
}
