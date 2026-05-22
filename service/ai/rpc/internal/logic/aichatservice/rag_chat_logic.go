package aichatservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RagChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRagChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RagChatLogic {
	return &RagChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 执行普通 RAG 问答，返回完整答案。
func (l *RagChatLogic) RagChat(in *pb.RagChatReq) (*pb.RagChatResp, error) {
	// todo: add your logic here and delete this line

	return &pb.RagChatResp{}, nil
}
