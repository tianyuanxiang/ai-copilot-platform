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
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	keyword := strings.TrimSpace(req.Keyword)
	hasStatus := req.Status != 0

	list, err := l.svcCtx.AiKnowledgeClient.ListDomain(l.ctx, &aiknowledgeclient.ListDomainReq{
		Page:      int64(req.Page),
		PageSize:  int64(req.PageSize),
		Keyword:   keyword,
		Status:    int64(req.Status),
		HasStatus: hasStatus,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.AiKbDomainItem, 0, len(list.List))
	for _, item := range list.List {
		items = append(items, domainItemFromRPC(item))
	}

	return &types.AiListKbDomainResp{
		Total: list.Total,
		List:  items,
	}, nil
}
