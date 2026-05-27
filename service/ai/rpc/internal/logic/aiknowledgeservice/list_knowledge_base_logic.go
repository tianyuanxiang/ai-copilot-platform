package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListKnowledgeBaseLogic {
	return &ListKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询当前用户可访问的知识库列表。
func (l *ListKnowledgeBaseLogic) ListKnowledgeBase(in *pb.ListKnowledgeBaseReq) (*pb.ListKnowledgeBaseResp, error) {
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 获取用户可访问的知识库 ID 列表
	scope := "all"
	hasDomainId := in.HasDomainId && in.DomainId > 0
	if in.HasKbType {
		scope = in.KbType
	}
	accessibleIDs, err := l.svcCtx.AiKnowledgeBaseModel.FindAccessibleKnowledgeBaseIDsByScope(l.ctx, in.UserId, scope, in.DomainId, hasDomainId)

	if err != nil {
		return nil, err
	}
	if len(accessibleIDs) == 0 {
		return &pb.ListKnowledgeBaseResp{Total: 0, List: []*pb.KnowledgeBaseItem{}}, nil
	}

	// 按条件过滤并分页
	kbs, total, err := l.svcCtx.AiKnowledgeBaseModel.ListByCondition(l.ctx, model.KbListQuery{
		KbType:        in.KbType,
		HasKbType:     in.HasKbType,
		DomainId:      in.DomainId,
		HasDomainId:   in.HasDomainId,
		Keyword:       in.Keyword,
		Visibility:    in.Visibility,
		HasVisibility: in.HasVisibility,
		Status:        in.Status,
		HasStatus:     in.HasStatus,
		Page:          page,
		PageSize:      pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 过滤只保留用户可访问的知识库
	accessibleSet := make(map[int64]bool, len(accessibleIDs))
	for _, id := range accessibleIDs {
		accessibleSet[id] = true
	}

	items := make([]*pb.KnowledgeBaseItem, 0)
	for i := range kbs {
		if !accessibleSet[kbs[i].Id] {
			continue
		}
		docCount, _ := l.svcCtx.AiKnowledgeBaseModel.CountDocumentsByKbID(l.ctx, kbs[i].Id)
		domainName := model.GetDomainNameByID(l.ctx, l.svcCtx.Orm, model.NullInt64Value(kbs[i].DomainId))
		items = append(items, knowledgeBaseToPB(&kbs[i], docCount, domainName))
	}

	return &pb.ListKnowledgeBaseResp{
		Total: total,
		List:  items,
	}, nil
}
