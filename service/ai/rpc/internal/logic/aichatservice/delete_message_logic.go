package aichatservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
	return &DeleteMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMessageLogic) DeleteMessage(in *pb.DeleteMessageReq) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if in.MessageId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "messageId 不能为空")
	}
	conversationID, err := parseConversationID(in.ConversationId)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.AiMessageModel.SoftDeleteByIDUserID(l.ctx, conversationID, in.MessageId, in.UserId, in.UserId); err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "消息不存在")
		}
		return nil, err
	}
	if err := NewRagChatLogic(l.ctx, l.svcCtx).refreshConversationSummary(conversationID, in.UserId); err != nil {
		l.Errorf("refresh conversation summary after delete message failed: %v", err)
	}
	return &pb.Empty{}, nil
}
