// Package aiwindtool 的 draft 文件实现三类只落草稿、不触发真实业务动作的工具。
package aiwindtool

import (
	"context"
	"fmt"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/logic/aiwinddraft"
	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/pb"
)

// executeGenerateAlarmAnalysisDraft 生成告警分析草稿。
// 这里不能直接调用现有 AnalyzeAlarmLogic：该 logic 会创建新的 traceId 并单独写工具日志。
// Agent 工具层需要沿用父 traceId，并由 Executor 统一只写一条审计日志。
func (e *Executor) executeGenerateAlarmAnalysisDraft(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args GenerateAlarmAnalysisDraftArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return nil, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return nil, err
	}
	startTime, endTime, err := normalizeTimeRange(args.StartTime, args.EndTime)
	if err != nil {
		return nil, err
	}

	alarmEvidence, err := aiwinddraft.BuildAlarmEvidence(
		ctx, e.svcCtx, farmCode, strings.TrimSpace(args.TowerCode), strings.TrimSpace(args.AlarmCode),
		startTime, endTime, 0, true,
	)
	if err != nil {
		return nil, err
	}
	evidenceItems, evidenceJSON := aiwinddraft.MergeEvidence(args.EvidenceJSON, alarmEvidence)
	payload := engine.WindDraftRequest{
		UserID:       aiwinddraft.UserIDString(req.UserId),
		TraceID:      req.TraceId,
		FarmCode:     farmCode,
		TowerCode:    strings.TrimSpace(args.TowerCode),
		AlarmCode:    strings.TrimSpace(args.AlarmCode),
		StartTime:    startTime,
		EndTime:      endTime,
		Evidence:     evidenceItems,
		EvidenceJSON: evidenceJSON,
	}
	draft, err := e.svcCtx.EngineCallClient.WindAlarmSummary(ctx, payload)
	if err != nil {
		return nil, err
	}

	status := draftStatus(draft)
	content := aiwinddraft.DraftContent(draft)
	evidenceSummary := aiwinddraft.BuildEvidenceSummary(evidenceJSON)
	id, err := e.svcCtx.AiAlarmAnalysisModel.InsertReturningID(ctx, &model.AiAlarmAnalysis{
		UserId:    req.UserId,
		TraceId:   req.TraceId,
		FarmCode:  farmCode,
		TowerCode: payload.TowerCode,
		AlarmCode: payload.AlarmCode,
		Title:     draft.Title,
		Content:   content,
		Evidence:  evidenceSummary,
		Status:    status,
	})
	if err != nil {
		return nil, err
	}
	return buildDraftToolResult("alarm_analysis", id, draft.Title, content, evidenceSummary, status, draft.Message), nil
}

// executeGenerateHealthReport 生成健康报告草稿。
// 这里不能直接调用现有 GenerateHealthReportLogic，原因同告警分析草稿：
// 必须保留 Agent 父 traceId，并避免重复写入工具日志。
func (e *Executor) executeGenerateHealthReport(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args GenerateHealthReportArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return nil, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return nil, err
	}
	startTime, endTime, err := normalizeTimeRange(args.StartTime, args.EndTime)
	if err != nil {
		return nil, err
	}
	reportType := strings.TrimSpace(args.ReportType)
	if reportType == "" {
		reportType = "health"
	}

	alarmEvidence, err := aiwinddraft.BuildAlarmEvidence(
		ctx, e.svcCtx, farmCode, strings.TrimSpace(args.TowerCode), "",
		startTime, endTime, 0, false,
	)
	if err != nil {
		return nil, err
	}
	evidenceItems, evidenceJSON := aiwinddraft.MergeEvidence(args.EvidenceJSON, alarmEvidence)
	payload := engine.WindDraftRequest{
		UserID:       aiwinddraft.UserIDString(req.UserId),
		TraceID:      req.TraceId,
		FarmCode:     farmCode,
		TowerCode:    strings.TrimSpace(args.TowerCode),
		ReportType:   reportType,
		StartTime:    startTime,
		EndTime:      endTime,
		Evidence:     evidenceItems,
		EvidenceJSON: evidenceJSON,
	}
	draft, err := e.svcCtx.EngineCallClient.WindHealthReportDraft(ctx, payload)
	if err != nil {
		return nil, err
	}

	status := draftStatus(draft)
	content := aiwinddraft.DraftContent(draft)
	evidenceSummary := aiwinddraft.BuildEvidenceSummary(evidenceJSON)
	id, err := e.svcCtx.AiHealthReportModel.InsertReturningID(ctx, &model.AiHealthReport{
		UserId:     req.UserId,
		TraceId:    req.TraceId,
		ReportType: reportType,
		FarmCode:   farmCode,
		TowerCode:  payload.TowerCode,
		StartTime:  aiwinddraft.ParseNullableTime(startTime),
		EndTime:    aiwinddraft.ParseNullableTime(endTime),
		Title:      draft.Title,
		Content:    content,
		Evidence:   evidenceSummary,
		Status:     status,
	})
	if err != nil {
		return nil, err
	}
	return buildDraftToolResult("health_report", id, draft.Title, content, evidenceSummary, status, draft.Message), nil
}

// executeCreateMaintenanceTicketDraft 生成维修工单草稿。
// 这里不能直接调用现有 CreateTicketDraftLogic，因为 Agent 链路必须沿用父 traceId，
// 且本工具只创建草稿，不创建正式工单、不派单、不通知。
func (e *Executor) executeCreateMaintenanceTicketDraft(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args CreateMaintenanceTicketDraftArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return nil, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return nil, err
	}
	priority := strings.TrimSpace(args.Priority)
	if priority == "" {
		priority = "normal"
	}
	// 工单工具参数没有显式时间窗口。默认取最近 24 小时告警作为自动证据，
	// 避免无边界扫描 TDengine。
	startTime, endTime, err := normalizeTimeRange("", "")
	if err != nil {
		return nil, err
	}

	alarmEvidence, err := aiwinddraft.BuildAlarmEvidence(
		ctx, e.svcCtx, farmCode, strings.TrimSpace(args.TowerCode), strings.TrimSpace(args.AlarmCode),
		startTime, endTime, 0, true,
	)
	if err != nil {
		return nil, err
	}
	evidenceItems, evidenceJSON := aiwinddraft.MergeEvidence(args.EvidenceJSON, alarmEvidence)
	payload := engine.WindDraftRequest{
		UserID:       aiwinddraft.UserIDString(req.UserId),
		TraceID:      req.TraceId,
		FarmCode:     farmCode,
		TowerCode:    strings.TrimSpace(args.TowerCode),
		AlarmCode:    strings.TrimSpace(args.AlarmCode),
		Priority:     priority,
		Evidence:     evidenceItems,
		EvidenceJSON: evidenceJSON,
	}
	draft, err := e.svcCtx.EngineCallClient.WindTicketDraft(ctx, payload)
	if err != nil {
		return nil, err
	}

	status := draftStatus(draft)
	content := aiwinddraft.DraftContent(draft)
	evidenceSummary := aiwinddraft.BuildEvidenceSummary(evidenceJSON)
	id, err := e.svcCtx.AiMaintenanceTicketDraftModel.InsertReturningID(ctx, &model.AiMaintenanceTicketDraft{
		UserId:    req.UserId,
		TraceId:   req.TraceId,
		FarmCode:  farmCode,
		TowerCode: payload.TowerCode,
		AlarmCode: payload.AlarmCode,
		Title:     draft.Title,
		Content:   content,
		Evidence:  evidenceSummary,
		Status:    status,
	})
	if err != nil {
		return nil, err
	}
	return buildDraftToolResult("ticket", id, draft.Title, content, evidenceSummary, status, draft.Message), nil
}

func draftStatus(draft *engine.WindDraftResponse) string {
	if draft == nil || strings.TrimSpace(draft.Status) == "" {
		return "draft"
	}
	return draft.Status
}

func buildDraftToolResult(draftType string, id int64, title, content, evidenceJSON, status, message string) *toolResult {
	resultJSON := marshalJSON(map[string]any{
		"draftType":    draftType,
		"draftId":      id,
		"title":        title,
		"content":      content,
		"evidenceJson": evidenceJSON,
		"status":       status,
		"message":      message,
	})
	return &toolResult{
		ResultJSON:   resultJSON,
		EvidenceJSON: evidenceJSON,
		Message:      fmt.Sprintf("%s; status=%s", message, status),
	}
}
