package aichatservicelogic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

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

func (l *GetConversationMessagesLogic) GetConversationMessages(in *pb.GetConversationMessagesReq) (*pb.GetConversationMessagesResp, error) {
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	conversationID, err := parseConversationID(in.ConversationId)
	if err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.AiConversationModel.FindByIDUserID(l.ctx, conversationID, in.UserId); err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "会话不存在")
		}
		return nil, err
	}

	list, total, err := l.svcCtx.AiMessageModel.ListByConversation(l.ctx, conversationID, in.UserId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetConversationMessagesResp{
		Total: total,
		List:  make([]*pb.MessageItem, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, &pb.MessageItem{
			MessageId:      item.MessageId,
			ConversationId: strconv.FormatInt(item.ConversationId, 10),
			Role:           item.Role,
			Content:        item.Content,
			Citations:      parseCitationsJSON(item.Citations),
			TraceId:        item.TraceId,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}

func parseConversationID(raw string) (int64, error) {
	conversationID, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || conversationID <= 0 {
		return 0, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "conversationId 必须是有效数字")
	}
	return conversationID, nil
}
