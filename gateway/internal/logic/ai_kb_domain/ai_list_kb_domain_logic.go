// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_domain

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiListKbDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListKbDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListKbDomainLogic {
	return &AiListKbDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListKbDomainLogic) AiListKbDomain(req *types.AiListKbDomainReq) (resp *types.AiListKbDomainResp, err error) {
	// todo: add your logic here and delete this line

	return
}
