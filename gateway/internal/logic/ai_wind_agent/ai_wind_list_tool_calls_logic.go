// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package ai_wind_agent

import (
	"context"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/svc"
	"ai-copilot-platform/gateway/internal/types"
	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiWindListToolCallsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiWindListToolCallsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiWindListToolCallsLogic {
	return &AiWindListToolCallsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiWindListToolCallsLogic) AiWindListToolCalls(req *types.AiWindListToolCallReq) (resp *types.AiWindListToolCallResp, err error) {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	result, err := l.svcCtx.AiWindAgentClient.ListToolCallLog(l.ctx, &pb.WindListToolCallLogReq{
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		UserId:   userID,
		TraceId:  req.TraceId,
		ToolName: req.ToolName,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.AiWindToolCallItem, 0, len(result.List))
	for _, item := range result.List {
		if item == nil {
			continue
		}
		list = append(list, types.AiWindToolCallItem{
			ToolCallId:    item.ToolCallId,
			ToolName:      item.ToolName,
			Status:        item.Status,
			ArgumentsJson: item.ArgumentsJson,
			ResultJson:    item.ResultJson,
			Message:       item.Message,
			CreatedAt:     item.CreatedAt,
		})
	}
	return &types.AiWindListToolCallResp{Total: result.Total, List: list}, nil
}
