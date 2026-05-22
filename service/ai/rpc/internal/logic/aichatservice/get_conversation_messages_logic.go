package aichatservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationMessagesLogic {
	return &GetConversationMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询指定会话的消息历史。
func (l *GetConversationMessagesLogic) GetConversationMessages(in *pb.GetConversationMessagesReq) (*pb.GetConversationMessagesResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetConversationMessagesResp{}, nil
}
