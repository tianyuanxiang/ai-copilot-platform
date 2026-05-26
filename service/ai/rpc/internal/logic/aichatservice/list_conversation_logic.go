package aichatservicelogic

import (
	"context"
	"strconv"
	"time"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

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

func (l *ListConversationLogic) ListConversation(in *pb.ListConversationReq) (*pb.ListConversationResp, error) {
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	list, total, err := l.svcCtx.AiConversationModel.ListByUser(l.ctx, model.ConversationListQuery{
		Page:     in.Page,
		PageSize: in.PageSize,
		UserID:   in.UserId,
		KbID:     in.KbId,
		HasKbID:  in.HasKbId && in.KbId > 0,
		Keyword:  in.Keyword,
	})
	if err != nil {
		return nil, err
	}

	resp := &pb.ListConversationResp{
		Total: total,
		List:  make([]*pb.ConversationItem, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, &pb.ConversationItem{
			ConversationId: strconv.FormatInt(item.Id, 10),
			KbId:           item.KbId,
			Title:          item.Title,
			LatestMessage:  item.LatestMessage,
			CreatedAt:      item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      item.UpdatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}
