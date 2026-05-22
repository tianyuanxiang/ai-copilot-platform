// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_security

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListSecurityEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListSecurityEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListSecurityEventsLogic {
	return &AiListSecurityEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListSecurityEventsLogic) AiListSecurityEvents(req *types.AiListSecurityEventsReq) (resp *types.AiListSecurityEventsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
