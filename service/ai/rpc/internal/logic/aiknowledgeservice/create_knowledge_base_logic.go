package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateKnowledgeBaseLogic {
	return &CreateKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建个人或公共知识库。
func (l *CreateKnowledgeBaseLogic) CreateKnowledgeBase(in *pb.CreateKnowledgeBaseReq) (*pb.CreateKnowledgeBaseResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateKnowledgeBaseResp{}, nil
}
