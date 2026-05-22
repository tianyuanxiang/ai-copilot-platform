// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_member

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiAddKbMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiAddKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiAddKbMemberLogic {
	return &AiAddKbMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiAddKbMemberLogic) AiAddKbMember(req *types.AiAddKbMemberReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
