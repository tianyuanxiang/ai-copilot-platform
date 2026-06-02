package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDocumentLogic {
	return &ListDocumentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDocumentLogic) ListDocument(in *pb.ListDocumentReq) (*pb.ListDocumentResp, error) {
	if in.UserId <= 0 || in.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "user_id 和 kb_id 不能为空")
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

	status := ""
	if in.HasStatus {
		status = in.Status
	}
	docs, total, err := l.svcCtx.AiDocumentModel.ListByKbID(l.ctx, model.AiDocumentListQuery{
		KbID:     in.KbId,
		Page:     in.Page,
		PageSize: in.PageSize,
		Keyword:  in.Keyword,
		Status:   status,
	})
	if err != nil {
		l.Logger.Errorf("获取知识库列表错误 %v", err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "获取知识库列表错误")
	}

	items := make([]*pb.DocumentItem, 0, len(docs))
	for i := range docs {
		items = append(items, documentToPB(&docs[i]))
	}
	return &pb.ListDocumentResp{
		Total: total,
		List:  items,
	}, nil
}
