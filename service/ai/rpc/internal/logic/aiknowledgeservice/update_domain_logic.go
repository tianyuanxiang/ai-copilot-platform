package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDomainLogic {
	return &UpdateDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新公共知识库领域基础信息。
func (l *UpdateDomainLogic) UpdateDomain(in *pb.UpdateDomainReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
