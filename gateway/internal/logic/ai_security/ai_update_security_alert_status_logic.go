// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_security

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiUpdateSecurityAlertStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiUpdateSecurityAlertStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiUpdateSecurityAlertStatusLogic {
	return &AiUpdateSecurityAlertStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiUpdateSecurityAlertStatusLogic) AiUpdateSecurityAlertStatus(req *types.AiUpdateSecurityAlertStatusReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
