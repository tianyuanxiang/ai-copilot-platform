package aichatservicelogic

import (
	"context"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConversationTitleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConversationTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationTitleLogic {
	return &UpdateConversationTitleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateConversationTitleLogic) UpdateConversationTitle(in *pb.UpdateConversationTitleReq) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	conversationID, err := parseConversationID(in.ConversationId)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "title 不能为空")
	}
	if err := l.svcCtx.AiConversationModel.UpdateTitleByIDUserID(l.ctx, conversationID, in.UserId, truncateRunes(title, 80)); err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "会话不存在")
		}
		return nil, err
	}
	return &pb.Empty{}, nil
}
