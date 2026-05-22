package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.KnowledgeBaseItem{}, nil
}
