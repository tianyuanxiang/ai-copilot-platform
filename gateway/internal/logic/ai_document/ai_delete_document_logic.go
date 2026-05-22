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

type AiDeleteDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiDeleteDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiDeleteDocumentLogic {
	return &AiDeleteDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiDeleteDocumentLogic) AiDeleteDocument(req *types.AiDocumentPathReq) (resp *types.AiCommonResp, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if req.KbId <= 0 || req.DocumentId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 和 documentId 不能为空")
	}

	if _, err := l.svcCtx.AiKnowledgeClient.DeleteDocument(l.ctx, &aiknowledgeclient.DeleteDocumentReq{
		DocumentId: req.DocumentId,
		KbId:       req.KbId,
		UserId:     userID,
		OperatorId: userID,
	}); err != nil {
		return nil, err
	}

	return &types.AiCommonResp{Message: "document deleted"}, nil
}
