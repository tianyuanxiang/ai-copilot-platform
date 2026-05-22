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
	// todo: add your logic here and delete this line

	return &pb.ListDomainResp{}, nil
}
