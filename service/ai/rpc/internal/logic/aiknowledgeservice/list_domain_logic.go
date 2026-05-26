package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDomainLogic {
	return &ListDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询公共知识库领域列表。
func (l *ListDomainLogic) ListDomain(in *pb.ListDomainReq) (*pb.ListDomainResp, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	domains, total, err := l.svcCtx.AiKbDomainModel.List(l.ctx, page, pageSize, in.Keyword, in.Status, in.HasStatus)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.DomainItem, 0, len(domains))
	for i := range domains {
		items = append(items, domainToPB(&domains[i]))
	}
	return &pb.ListDomainResp{
		Total: total,
		List:  items,
	}, nil
}
