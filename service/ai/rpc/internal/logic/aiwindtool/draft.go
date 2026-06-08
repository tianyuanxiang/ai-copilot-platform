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

func (e *Executor) executeGenerateAlarmAnalysisDraft(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args GenerateAlarmAnalysisDraftArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	startTime, endTime, err := normalizeTimeRange(args.StartTime, args.EndTime)
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}

	payload, err := e.baseDraftPayload(ctx, req, draftPayloadInput{
		FarmCode:      farmCode,
		TowerCode:     args.TowerCode,
		AlarmCode:     args.AlarmCode,
		StartTime:     startTime,
		EndTime:       endTime,
		InputEvidence: args.EvidenceJSON,
		HasStatus:     true,
	})
	if err != nil {
		return nil, err
	}
	draft, err := e.svcCtx.EngineCallClient.WindAlarmSummary(ctx, payload)
	if err != nil {
		return nil, err
	}

	status, content, evidenceSummary := draftResultParts(draft, payload.EvidenceJSON)
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

func (e *Executor) executeGenerateHealthReport(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args GenerateHealthReportArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	startTime, endTime, err := normalizeTimeRange(args.StartTime, args.EndTime)
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	reportType := strings.TrimSpace(args.ReportType)
	if reportType == "" {
		reportType = "health"
	}

	payload, err := e.baseDraftPayload(ctx, req, draftPayloadInput{
		FarmCode:      farmCode,
		TowerCode:     args.TowerCode,
		StartTime:     startTime,
		EndTime:       endTime,
		InputEvidence: args.EvidenceJSON,
		HasStatus:     false,
	})
	if err != nil {
		return nil, err
	}
	payload.ReportType = reportType

	draft, err := e.svcCtx.EngineCallClient.WindHealthReportDraft(ctx, payload)
	if err != nil {
		return nil, err
	}

	status, content, evidenceSummary := draftResultParts(draft, payload.EvidenceJSON)
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

func (e *Executor) executeCreateMaintenanceTicketDraft(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args CreateMaintenanceTicketDraftArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	priority := strings.TrimSpace(args.Priority)
	if priority == "" {
		priority = "normal"
	}
	startTime, endTime, err := normalizeTimeRange("", "")
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}

	payload, err := e.baseDraftPayload(ctx, req, draftPayloadInput{
		FarmCode:      farmCode,
		TowerCode:     args.TowerCode,
		AlarmCode:     args.AlarmCode,
		StartTime:     startTime,
		EndTime:       endTime,
		InputEvidence: args.EvidenceJSON,
		HasStatus:     true,
	})
	if err != nil {
		return nil, err
	}
	payload.Priority = priority

	draft, err := e.svcCtx.EngineCallClient.WindTicketDraft(ctx, payload)
	if err != nil {
		return nil, err
	}

	status, content, evidenceSummary := draftResultParts(draft, payload.EvidenceJSON)
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

type draftPayloadInput struct {
	FarmCode      string
	TowerCode     string
	AlarmCode     string
	StartTime     string
	EndTime       string
	InputEvidence string
	HasStatus     bool
}

func (e *Executor) baseDraftPayload(ctx context.Context, req *pb.WindToolExecuteReq, input draftPayloadInput) (engine.WindDraftRequest, error) {
	towerCode := strings.TrimSpace(input.TowerCode)
	alarmCode := strings.TrimSpace(input.AlarmCode)
	alarmEvidence, err := aiwinddraft.BuildAlarmEvidence(
		ctx, e.svcCtx, input.FarmCode, towerCode, alarmCode,
		input.StartTime, input.EndTime, 0, input.HasStatus,
	)
	if err != nil {
		return engine.WindDraftRequest{}, err
	}
	evidenceItems, evidenceJSON := aiwinddraft.MergeEvidence(input.InputEvidence, alarmEvidence)
	return engine.WindDraftRequest{
		UserID:       aiwinddraft.UserIDString(req.UserId),
		TraceID:      req.TraceId,
		FarmCode:     input.FarmCode,
		TowerCode:    towerCode,
		AlarmCode:    alarmCode,
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		Evidence:     evidenceItems,
		EvidenceJSON: evidenceJSON,
	}, nil
}

func draftResultParts(draft *engine.WindDraftResponse, evidenceJSON string) (string, string, string) {
	status := draftStatus(draft)
	content := aiwinddraft.DraftContent(draft)
	evidenceSummary := aiwinddraft.BuildEvidenceSummary(evidenceJSON)
	return status, content, evidenceSummary
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
