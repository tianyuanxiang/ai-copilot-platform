// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_public

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiDeletePublicKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiDeletePublicKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiDeletePublicKbLogic {
	return &AiDeletePublicKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiDeletePublicKbLogic) AiDeletePublicKb(req *types.AiPublicKbPathReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
