// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_security

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListSecurityAlertsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListSecurityAlertsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListSecurityAlertsLogic {
	return &AiListSecurityAlertsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListSecurityAlertsLogic) AiListSecurityAlerts(req *types.AiListSecurityAlertsReq) (resp *types.AiListSecurityAlertsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
