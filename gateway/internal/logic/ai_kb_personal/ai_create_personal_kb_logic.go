// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_personal

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiCreatePersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiCreatePersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiCreatePersonalKbLogic {
	return &AiCreatePersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiCreatePersonalKbLogic) AiCreatePersonalKb(req *types.AiCreatePersonalKbReq) (resp *types.AiCreatePersonalKbResp, err error) {
	// todo: add your logic here and delete this line

	return
}
