// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_member

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListKbMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListKbMemberLogic {
	return &AiListKbMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListKbMemberLogic) AiListKbMember(req *types.AiListKbMemberReq) (resp *types.AiListKbMemberResp, err error) {
	// todo: add your logic here and delete this line

	return
}
