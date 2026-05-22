package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocumentLogic {
	return &DeleteDocumentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteDocument removes one document and its chunks/indexes.
func (l *DeleteDocumentLogic) DeleteDocument(in *pb.DeleteDocumentReq) (*pb.Empty, error) {
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
	if !canMaintainKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有维护该知识库的权限")
	}

	doc, err := l.svcCtx.AiDocumentModel.FindByIDKbID(l.ctx, in.DocumentId, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "文档不存在")
		}
		return nil, err
	}
	if doc.Status == documentStatusParsing || doc.Status == documentStatusIndexing {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "文档正在处理，不能删除")
	}

	if err := deleteDocumentFromElasticsearch(l.ctx, l.svcCtx, doc.Id); err != nil {
		return nil, err
	}

	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err := l.svcCtx.AiDocumentChunkModel.DeleteByDocumentIDTrans(l.ctx, tx, doc.Id); err != nil {
			return err
		}
		if err := l.svcCtx.AiDocumentParentChunkModel.DeleteByDocumentIDTrans(l.ctx, tx, doc.Id); err != nil {
			return err
		}
		return l.svcCtx.AiDocumentModel.DeleteByIDKbIDTrans(l.ctx, tx, doc.Id, in.KbId)
	})
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
