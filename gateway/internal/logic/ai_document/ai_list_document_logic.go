// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_document

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

type AiListDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiListDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiListDocumentLogic {
	return &AiListDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiListDocumentLogic) AiListDocument(req *types.AiListDocumentReq) (resp *types.AiListDocumentResp, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}

	status := strings.TrimSpace(req.Status)
	list, err := l.svcCtx.AiKnowledgeClient.ListDocument(l.ctx, &aiknowledgeclient.ListDocumentReq{
		Page:      int64(req.Page),
		PageSize:  int64(req.PageSize),
		UserId:    userID,
		KbId:      req.KbId,
		Keyword:   strings.TrimSpace(req.Keyword),
		Status:    status,
		HasStatus: status != "",
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.AiDocumentItem, 0, len(list.List))
	for _, doc := range list.List {
		items = append(items, documentItemFromRPC(doc))
	}
	return &types.AiListDocumentResp{
		Total: list.Total,
		List:  items,
	}, nil
}
