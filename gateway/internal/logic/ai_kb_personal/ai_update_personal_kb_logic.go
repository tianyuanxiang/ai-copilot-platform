// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_personal

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiUpdatePersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiUpdatePersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiUpdatePersonalKbLogic {
	return &AiUpdatePersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiUpdatePersonalKbLogic) AiUpdatePersonalKb(req *types.AiUpdatePersonalKbReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
