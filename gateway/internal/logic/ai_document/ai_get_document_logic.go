// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_document

import (
	"context"

	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiGetDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiGetDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiGetDocumentLogic {
	return &AiGetDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiGetDocumentLogic) AiGetDocument(req *types.AiDocumentPathReq) (resp *types.AiDocumentItem, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.KbId <= 0 || req.DocumentId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 和 documentId 不能为空")
	}

	doc, err := l.svcCtx.AiKnowledgeClient.GetDocument(l.ctx, &aiknowledgeclient.GetDocumentReq{
		DocumentId: req.DocumentId,
		KbId:       req.KbId,
		UserId:     userID,
	})
	if err != nil {
		return nil, err
	}
	item := documentItemFromRPC(doc)
	return &item, nil
}
