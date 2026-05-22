// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_public

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiCreatePublicKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiCreatePublicKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiCreatePublicKbLogic {
	return &AiCreatePublicKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiCreatePublicKbLogic) AiCreatePublicKb(req *types.AiCreatePublicKbReq) (resp *types.AiCreatePublicKbResp, err error) {
	// todo: add your logic here and delete this line

	return
}
