package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDocumentLogic {
	return &GetDocumentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDocumentLogic) GetDocument(in *pb.GetDocumentReq) (*pb.DocumentItem, error) {
	if in.UserId <= 0 || in.KbId <= 0 || in.DocumentId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "user_id、kb_id、document_id 不能为空")
	}

	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
		}
		return nil, err
	}
	if !canAccessKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有访问该知识库的权限")
	}

	doc, err := l.svcCtx.AiDocumentModel.FindByIDKbID(l.ctx, in.DocumentId, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "文档不存在")
		}
		return nil, err
	}
	return documentToPB(doc), nil
}
