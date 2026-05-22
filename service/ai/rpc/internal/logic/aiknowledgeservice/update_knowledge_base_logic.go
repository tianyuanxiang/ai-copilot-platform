package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateKnowledgeBaseLogic {
	return &UpdateKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新知识库名称、描述、领域、可见性或状态。
func (l *UpdateKnowledgeBaseLogic) UpdateKnowledgeBase(in *pb.UpdateKnowledgeBaseReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
