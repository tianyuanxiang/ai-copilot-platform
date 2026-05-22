// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_member

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiUpdateKbMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiUpdateKbMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiUpdateKbMemberLogic {
	return &AiUpdateKbMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiUpdateKbMemberLogic) AiUpdateKbMember(req *types.AiUpdateKbMemberReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
