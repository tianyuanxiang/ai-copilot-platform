package aichatservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationLogic {
	return &DeleteConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteConversationLogic) DeleteConversation(in *pb.DeleteConversationReq) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	conversationID, err := parseConversationID(in.ConversationId)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.AiConversationModel.SoftDeleteByIDUserID(l.ctx, conversationID, in.UserId, in.UserId); err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "会话不存在")
		}
		return nil, err
	}
	return &pb.Empty{}, nil
}
