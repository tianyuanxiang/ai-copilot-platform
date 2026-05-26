// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_kb_domain

import (
	"context"
	"strings"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

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
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "name 不能为空")
	}
	if strings.TrimSpace(req.Code) == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "code 不能为空")
	}

	rpcResp, err := l.svcCtx.AiKnowledgeClient.CreateDomain(l.ctx, &aiknowledgeclient.CreateDomainReq{
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.TrimSpace(req.Code),
		Description: req.Description,
		Sort:        int64(req.Sort),
		Status:      int64(req.Status),
		OperatorId:  userID,
	})
	if err != nil {
		return nil, err
	}

	return &types.AiCreateKbDomainResp{
		DomainId: rpcResp.DomainId,
	}, nil
}
