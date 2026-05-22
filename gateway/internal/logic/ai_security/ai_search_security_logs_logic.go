// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_security

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiSearchSecurityLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiSearchSecurityLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiSearchSecurityLogsLogic {
	return &AiSearchSecurityLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiSearchSecurityLogsLogic) AiSearchSecurityLogs(req *types.AiSecurityLogSearchReq) (resp *types.AiSecurityLogSearchResp, err error) {
	// todo: add your logic here and delete this line

	return
}
