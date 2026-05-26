package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetKnowledgeBaseLogic {
	return &GetKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询知识库详情。
func (l *GetKnowledgeBaseLogic) GetKnowledgeBase(in *pb.GetKnowledgeBaseReq) (*pb.KnowledgeBaseItem, error) {
	if in.KbId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "kbId 不能为空")
	}
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
		}
		return nil, err
	}

	// 检查访问权限
	if !canAccessKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有访问该知识库的权限")
	}

	// 获取文档数量
	docCount, _ := l.svcCtx.AiKnowledgeBaseModel.CountDocumentsByKbID(l.ctx, kb.Id)

	// 获取领域名称
	domainName := model.GetDomainNameByID(l.ctx, l.svcCtx.Orm, model.NullInt64Value(kb.DomainId))

	return knowledgeBaseToPB(kb, docCount, domainName), nil
}
