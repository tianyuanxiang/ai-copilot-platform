package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDomainLogic {
	return &DeleteDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除公共知识库领域。
func (l *DeleteDomainLogic) DeleteDomain(in *pb.DeleteDomainReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
