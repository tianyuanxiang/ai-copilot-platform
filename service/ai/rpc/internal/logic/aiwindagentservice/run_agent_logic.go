package aiwindagentservicelogic

import (
	"context"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunAgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRunAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunAgentLogic {
	return &RunAgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RunAgentLogic) RunAgent(in *pb.WindAgentRunReq) (*pb.WindAgentRunResp, error) {
	traceID := model.WindTraceID()
	tools := []string{
		"search_maintenance_sop",
		"query_sensor_timeseries",
		"query_alarm_events",
		"get_turbine_metadata",
		"compare_sensor_trend",
		"generate_health_report",
		"create_maintenance_ticket_draft",
	}
	calls := make([]*pb.ToolCall, 0, len(tools))
	for _, tool := range tools {
		calls = append(calls, &pb.ToolCall{
			ToolName:      tool,
			Status:        "reserved",
			ArgumentsJson: "{}",
			ResultJson:    "{}",
			Message:       "一期仅注册工具白名单和日志入口，后续补齐受控编排。",
		})
	}
	return &pb.WindAgentRunResp{
		Answer:    "风机混塔 AI Copilot Agent 脚手架已就绪：当前返回白名单工具，不执行复杂编排。",
		ToolCalls: calls,
		TraceId:   traceID,
	}, nil
}
