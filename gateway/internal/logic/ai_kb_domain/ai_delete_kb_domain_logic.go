// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_domain

import (
	"context"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiDeleteKbDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiDeleteKbDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiDeleteKbDomainLogic {
	return &AiDeleteKbDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiDeleteKbDomainLogic) AiDeleteKbDomain(req *types.AiKbDomainPathReq) (resp *types.AiCommonResp, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.DomainId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "domainId 不能为空")
	}

	if _, err := l.svcCtx.AiKnowledgeClient.DeleteDomain(l.ctx, &aiknowledgeclient.DeleteDomainReq{
		DomainId:   req.DomainId,
		OperatorId: userID,
	}); err != nil {
		return nil, err
	}

	return &types.AiCommonResp{Message: "ok"}, nil
}
