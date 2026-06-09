// Package aiwindtool 的 executor 文件实现所有 Agent 工具的受控入口。
// RPC 层只需要调用 Executor.Execute，工具白名单、分发和审计都在这里完成。
package aiwindtool

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	statusSuccess          = "success"
	statusInvalidArguments = "invalid_arguments"
	statusPermissionDenied = "permission_denied"
	statusRPCTimeout       = "rpc_timeout"
	statusRPCUnavailable   = "rpc_unavailable"
	statusToolFailed       = "tool_failed"
	statusToolPartial      = "tool_partial"
	statusUnknown          = "unknown"
	statusNoData           = "no_data"
)

// Executor 是 Wind Agent 的 Go 工具执行器。
// 它只持有现有 ServiceContext，不创建新的业务依赖。
type Executor struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewExecutor 创建一个 Wind Agent 工具执行器。
func NewExecutor(svcCtx *svc.ServiceContext) *Executor {
	return &Executor{svcCtx: svcCtx, Logger: logx.WithContext(context.Background())}
}

// Execute 校验 RPC 请求、分发白名单工具、记录最终审计日志，并返回结构化结果。
// 普通工具错误会放进 WindToolExecuteResp，而不是升级成 RPC error。
// 这样 LangGraph 可以根据失败原因补充参数、改用其他工具或结束推理。
func (e *Executor) Execute(ctx context.Context, req *pb.WindToolExecuteReq) (*pb.WindToolExecuteResp, error) {
	startedAt := time.Now()
	if req == nil {
		return failedResponse("", "请求不能为空", startedAt), nil
	}

	normalizedReq := *req
	normalizedReq.ToolName = strings.TrimSpace(req.ToolName)

	if normalizedReq.UserId <= 0 {
		return e.finish(ctx, &normalizedReq, statusInvalidArguments, nil, "userId 必须大于 0", startedAt), nil
	}
	if strings.TrimSpace(normalizedReq.TraceId) == "" {
		return e.finish(ctx, &normalizedReq, statusInvalidArguments, nil, "traceId 不能为空", startedAt), nil
	}
	if !IsAllowedTool(normalizedReq.ToolName) {
		return e.finish(ctx, &normalizedReq, statusPermissionDenied, nil, fmt.Sprintf("工具 %q 不在白名单中", normalizedReq.ToolName), startedAt), nil
	}
	e.Logger.Infof("开始执行工具: [%s]", normalizedReq.ToolName)
	e.Logger.Infof("请求参数为: [%s]", normalizedReq.ArgumentsJson)
	var (
		result *toolResult
		err    error
	)
	switch normalizedReq.ToolName {
	case ToolGetTurbineMetadata:
		result, err = e.executeGetTurbineMetadata(ctx, &normalizedReq)
	case ToolSearchMaintenanceSOP:
		result, err = e.executeSearchMaintenanceSOP(ctx, &normalizedReq)
	case ToolQueryAlarmEvents:
		result, err = e.executeQueryAlarmEvents(ctx, &normalizedReq)
	case ToolQuerySensorTimeseries:
		result, err = e.executeQuerySensorTimeseries(ctx, &normalizedReq)
	case ToolCompareSensorTrend:
		result, err = e.executeCompareSensorTrend(ctx, &normalizedReq)
	case ToolGenerateAlarmAnalysisDraft:
		result, err = e.executeGenerateAlarmAnalysisDraft(ctx, &normalizedReq)
	case ToolGenerateHealthReport:
		result, err = e.executeGenerateHealthReport(ctx, &normalizedReq)
	case ToolCreateMaintenanceTicketDraft:
		result, err = e.executeCreateMaintenanceTicketDraft(ctx, &normalizedReq)
	}
	if err != nil {
		status := statusToolFailed
		if result != nil && result.Status != "" {
			status = result.Status
		}
		e.Logger.Errorf("tool:[%s] execute err:%v", normalizedReq.ToolName, err)
		return e.finish(ctx, &normalizedReq, status, result, err.Error(), startedAt), nil
	}

	status := statusSuccess
	if result != nil && result.Status != "" {
		status = result.Status
	}
	return e.finish(ctx, &normalizedReq, status, result, result.Message, startedAt), nil
}

func (e *Executor) finish(
	ctx context.Context,
	req *pb.WindToolExecuteReq,
	status string,
	result *toolResult,
	message string,
	startedAt time.Time,
) *pb.WindToolExecuteResp {
	if result == nil {
		result = &toolResult{ResultJSON: "{}"}
	}
	if result.ResultJSON == "" {
		result.ResultJSON = "{}"
	}
	if message == "" {
		message = result.Message
	}

	latencyMs := time.Since(startedAt).Milliseconds()
	logID := writeToolCallLog(ctx, e.svcCtx, req, status, result.ResultJSON, message, latencyMs)
	return &pb.WindToolExecuteResp{
		ToolCallId:   logID,
		ToolName:     req.ToolName,
		Status:       status,
		ResultJson:   result.ResultJSON,
		EvidenceJson: result.EvidenceJSON,
		Citations:    result.Citations,
		Message:      message,
		LatencyMs:    latencyMs,
	}
}

func failedResponse(toolName, message string, startedAt time.Time) *pb.WindToolExecuteResp {
	return &pb.WindToolExecuteResp{
		ToolName:   toolName,
		Status:     statusToolFailed,
		ResultJson: "{}",
		Message:    message,
		LatencyMs:  time.Since(startedAt).Milliseconds(),
	}
}
