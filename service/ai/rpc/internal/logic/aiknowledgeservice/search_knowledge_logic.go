package aiknowledgeservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchKnowledgeLogic {
	return &SearchKnowledgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 执行知识库检索，支持 quick/deep 模式和不同检索范围。
func (l *SearchKnowledgeLogic) SearchKnowledge(in *pb.SearchKnowledgeReq) (*pb.SearchKnowledgeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SearchKnowledgeResp{}, nil
}
