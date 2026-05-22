// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_public

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListPublicKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListPublicKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListPublicKbLogic {
	return &AiListPublicKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListPublicKbLogic) AiListPublicKb(req *types.AiListPublicKbReq) (resp *types.AiListPublicKbResp, err error) {
	// todo: add your logic here and delete this line

	return
}
