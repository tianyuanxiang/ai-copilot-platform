package aiwindagentservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListToolCallLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListToolCallLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListToolCallLogLogic {
	return &ListToolCallLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListToolCallLogLogic) ListToolCallLog(in *pb.WindListToolCallLogReq) (*pb.WindListToolCallLogResp, error) {
	total, rows, err := l.svcCtx.AiToolCallLogModel.ListByFilter(l.ctx, in.UserId, in.TraceId, in.ToolName, in.Status, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}
	list := make([]*pb.ToolCall, 0, len(rows))
	for _, row := range rows {
		list = append(list, &pb.ToolCall{
			ToolCallId:    row.Id,
			ToolName:      row.ToolName,
			Status:        row.Status,
			ArgumentsJson: row.Arguments,
			ResultJson:    row.Result,
			Message:       row.ErrorMsg,
			CreatedAt:     row.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &pb.WindListToolCallLogResp{Total: total, List: list}, nil
}
