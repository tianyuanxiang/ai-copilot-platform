package aichatservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConversationLogic {
	return &ListConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询当前用户的 AI 会话列表。
func (l *ListConversationLogic) ListConversation(in *pb.ListConversationReq) (*pb.ListConversationResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ListConversationResp{}, nil
}
