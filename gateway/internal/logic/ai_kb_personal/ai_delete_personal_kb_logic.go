// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_personal

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiDeletePersonalKbLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiDeletePersonalKbLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiDeletePersonalKbLogic {
	return &AiDeletePersonalKbLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiDeletePersonalKbLogic) AiDeletePersonalKb(req *types.AiPersonalKbPathReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
