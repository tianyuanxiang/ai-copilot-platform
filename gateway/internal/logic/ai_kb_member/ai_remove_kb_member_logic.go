// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_member

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiRemoveKbMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiRemoveKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiRemoveKbMemberLogic {
	return &AiRemoveKbMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiRemoveKbMemberLogic) AiRemoveKbMember(req *types.AiKbMemberPathReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
