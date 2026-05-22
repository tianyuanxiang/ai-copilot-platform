package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteKnowledgeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteKnowledgeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKnowledgeBaseLogic {
	return &DeleteKnowledgeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除知识库，并由实现层清理关联文档和索引。
func (l *DeleteKnowledgeBaseLogic) DeleteKnowledgeBase(in *pb.DeleteKnowledgeBaseReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
