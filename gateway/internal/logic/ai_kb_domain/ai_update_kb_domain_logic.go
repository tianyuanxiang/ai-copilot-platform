// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_domain

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiUpdateKbDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiUpdateKbDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiUpdateKbDomainLogic {
	return &AiUpdateKbDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiUpdateKbDomainLogic) AiUpdateKbDomain(req *types.AiUpdateKbDomainReq) (resp *types.AiCommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
