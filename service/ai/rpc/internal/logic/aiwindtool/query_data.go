package aiwindtool

import (
	"context"
	"encoding/json"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/logic/aiwinddraft"
	aiwindtimeseriesservicelogic "ai-copilot-platform/ai-rpc/internal/logic/aiwindtimeseriesservice"
	"ai-copilot-platform/ai-rpc/pb"
)

func (e *Executor) executeQueryAlarmEvents(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args QueryAlarmEventsArgs
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

	evidence, err := aiwinddraft.BuildAlarmEvidence(
		ctx,
		e.svcCtx,
		farmCode,
		strings.TrimSpace(args.TowerCode),
		strings.TrimSpace(args.AlarmCode),
		startTime,
		endTime,
		args.Status,
		args.HasStatus,
	)
	if err != nil {
		return nil, err
	}
	evidenceJSON := marshalJSON(evidence)

	return &toolResult{
		ResultJSON:   evidenceJSON,
		EvidenceJSON: evidenceJSON,
		Message:      "alarm evidence query completed",
	}, nil
}

func (e *Executor) executeQuerySensorTimeseries(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args QuerySensorTimeseriesArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	normalized, err := normalizeSensorToolArgs(
		args.FarmCode, args.TowerCode, args.DeviceCode, args.DeviceTypeCode,
		args.Field, args.StartTime, args.EndTime, args.IndexID, args.RadarDistanceM,
		args.Page, args.PageSize, true,
	)
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}

	logic := aiwindtimeseriesservicelogic.NewQueryTimeseriesLogic(ctx, e.svcCtx)
	resp, err := logic.QueryTimeseries(&pb.WindTimeseriesQueryReq{
		FarmCode:       normalized.FarmCode,
		TowerCode:      normalized.TowerCode,
		DeviceCode:     normalized.DeviceCode,
		DeviceTypeCode: normalized.DeviceTypeCode,
		Field:          normalized.Fields,
		StartTime:      normalized.StartTime,
		EndTime:        normalized.EndTime,
		Page:           normalized.Page,
		PageSize:       normalized.PageSize,
		IndexId:        normalized.IndexID,
		RadarDistanceM: normalized.RadarDistanceM,
		UserId:         req.UserId,
	})
	if err != nil {
		return &toolResult{Status: statusToolFailed}, err
	}

	result := map[string]any{
		"total":    resp.Total,
		"database": resp.Database,
		"stable":   resp.Stable,
		"fields":   resp.Fields,
		"points":   resp.Points,
		"message":  resp.Message,
	}
	addEvidenceDisplayMeta(result, resp.EvidenceJson)

	if resp.Total == 0 {
		return &toolResult{
			ResultJSON:   marshalJSON(result),
			EvidenceJSON: resp.EvidenceJson,
			Message:      resp.Message,
			Status:       statusNoData,
		}, nil

	}

	return &toolResult{
		ResultJSON:   marshalJSON(result),
		EvidenceJSON: resp.EvidenceJson,
		Message:      resp.Message,
	}, nil
}

func (e *Executor) executeCompareSensorTrend(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args CompareSensorTrendArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	normalized, err := normalizeSensorToolArgs(
		args.FarmCode, args.TowerCode, args.DeviceCode, args.DeviceTypeCode,
		args.Field, args.StartTime, args.EndTime, args.IndexID, args.RadarDistanceM,
		0, 0, false,
	)
	if err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}

	logic := aiwindtimeseriesservicelogic.NewCompareTrendLogic(ctx, e.svcCtx)
	resp, err := logic.CompareTrend(&pb.WindTrendCompareReq{
		FarmCode:       normalized.FarmCode,
		TowerCode:      normalized.TowerCode,
		DeviceCode:     normalized.DeviceCode,
		DeviceTypeCode: normalized.DeviceTypeCode,
		Field:          normalized.Fields,
		StartTime:      normalized.StartTime,
		EndTime:        normalized.EndTime,
		IndexId:        normalized.IndexID,
		RadarDistanceM: normalized.RadarDistanceM,
		UserId:         req.UserId,
	})
	if err != nil {
		return nil, err
	}

	var status = ""
	if resp.Message == statusNoData {
		status = statusNoData
	}
	return &toolResult{
		ResultJSON: marshalJSON(map[string]any{
			"summary": resp.Summary,
			"message": resp.Message,
		}),
		EvidenceJSON: resp.EvidenceJson,
		Message:      resp.Message,
		Status:       status,
	}, nil
}

func addEvidenceDisplayMeta(result map[string]any, evidenceJSON string) {
	var evidence map[string]any
	if err := json.Unmarshal([]byte(evidenceJSON), &evidence); err != nil {
		return
	}
	for _, key := range []string{"deviceTypeName", "fieldLabels", "fieldUnits"} {
		if value, ok := evidence[key]; ok {
			result[key] = value
		}
	}
}
