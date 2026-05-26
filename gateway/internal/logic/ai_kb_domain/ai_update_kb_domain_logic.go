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
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.DomainId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "domainId 不能为空")
	}

	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	description := req.Description

	if _, err := l.svcCtx.AiKnowledgeClient.UpdateDomain(l.ctx, &aiknowledgeclient.UpdateDomainReq{
		DomainId:       req.DomainId,
		Name:           name,
		HasName:        name != "",
		Code:           code,
		HasCode:        code != "",
		Description:    description,
		HasDescription: description != "",
		Sort:           int64(req.Sort),
		HasSort:        req.Sort != 0,
		Status:         int64(req.Status),
		HasStatus:      req.Status != 0,
		OperatorId:     userID,
	}); err != nil {
		return nil, err
	}

	return &types.AiCommonResp{Message: "ok"}, nil
}
