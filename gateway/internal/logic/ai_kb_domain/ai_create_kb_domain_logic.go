// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_domain

import (
	"context"

	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiCreateKbDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiCreateKbDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiCreateKbDomainLogic {
	return &AiCreateKbDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiCreateKbDomainLogic) AiCreateKbDomain(req *types.AiCreateKbDomainReq) (resp *types.AiCreateKbDomainResp, err error) {
	// todo: add your logic here and delete this line

	return
}
