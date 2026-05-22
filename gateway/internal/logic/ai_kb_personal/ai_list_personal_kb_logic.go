// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_personal

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListPersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListPersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListPersonalKbLogic {
	return &AiListPersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListPersonalKbLogic) AiListPersonalKb(req *types.AiListPersonalKbReq) (resp *types.AiListPersonalKbResp, err error) {
	// todo: add your logic here and delete this line

	return
}
