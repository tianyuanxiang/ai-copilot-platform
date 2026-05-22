package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.ListKnowledgeBaseResp{}, nil
}
