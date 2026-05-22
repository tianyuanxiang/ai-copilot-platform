package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDomainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDomainLogic {
	return &CreateDomainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建公共知识库领域。
func (l *CreateDomainLogic) CreateDomain(in *pb.CreateDomainReq) (*pb.CreateDomainResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateDomainResp{}, nil
}
