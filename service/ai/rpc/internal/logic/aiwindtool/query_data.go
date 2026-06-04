// Package aiwindtool 的 query_data 文件实现告警、时序和趋势三类事实查询工具。
package aiwindtool

import (
	"context"
	"fmt"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/logic/aiwinddraft"
	aiwindtimeseriesservicelogic "ai-copilot-platform/ai-rpc/internal/logic/aiwindtimeseriesservice"
	"ai-copilot-platform/ai-rpc/pb"
)

// executeQueryAlarmEvents 复用 aiwinddraft.BuildAlarmEvidence。
// BuildAlarmEvidence 已经提供告警聚合、时间桶和分层抽样，避免向 Agent 塞入大量原始记录。

func (e *Executor) executeQueryAlarmEvents(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args QueryAlarmEventsArgs
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

// executeQuerySensorTimeseries 复用现有 QueryTimeseriesLogic。
// TDengine stable、字段白名单和 WPR 雷达 indexId 行为继续由现有 logic 负责。
func (e *Executor) executeQuerySensorTimeseries(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args QuerySensorTimeseriesArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return nil, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return nil, err
	}
	args.DeviceTypeCode = strings.ToUpper(strings.TrimSpace(args.DeviceTypeCode))
	if args.DeviceTypeCode == "" {
		return nil, fmt.Errorf("deviceTypeCode 不能为空")
	}
	if err := validateFields(args.Field); err != nil {
		return nil, err
	}
	args.Page, args.PageSize = normalizePage(args.Page, args.PageSize)
	args.StartTime, args.EndTime, err = normalizeTimeRange(args.StartTime, args.EndTime)
	if err != nil {
		return nil, err
	}

	logic := aiwindtimeseriesservicelogic.NewQueryTimeseriesLogic(ctx, e.svcCtx)
	resp, err := logic.QueryTimeseries(&pb.WindTimeseriesQueryReq{
		FarmCode:       farmCode,
		TowerCode:      strings.TrimSpace(args.TowerCode),
		DeviceCode:     strings.TrimSpace(args.DeviceCode),
		DeviceTypeCode: args.DeviceTypeCode,
		Field:          args.Field,
		StartTime:      args.StartTime,
		EndTime:        args.EndTime,
		Page:           args.Page,
		PageSize:       args.PageSize,
		IndexId:        args.IndexID,
		RadarDistanceM: args.RadarDistanceM,
		UserId:         req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &toolResult{
		ResultJSON: marshalJSON(map[string]any{
			"total":    resp.Total,
			"database": resp.Database,
			"stable":   resp.Stable,
			"fields":   resp.Fields,
			"points":   resp.Points,
			"message":  resp.Message,
		}),
		EvidenceJSON: resp.EvidenceJson,
		Message:      resp.Message,
	}, nil
}

// executeCompareSensorTrend 复用现有 CompareTrendLogic。
// 趋势统计、降采样、阈值解析和风险信号继续由已有业务 logic 负责。
func (e *Executor) executeCompareSensorTrend(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args CompareSensorTrendArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return nil, err
	}
	farmCode, err := requireFarmCode(args.FarmCode)
	if err != nil {
		return nil, err
	}
	args.DeviceTypeCode = strings.ToUpper(strings.TrimSpace(args.DeviceTypeCode))
	if args.DeviceTypeCode == "" {
		return nil, fmt.Errorf("deviceTypeCode 不能为空")
	}
	if err := validateFields(args.Field); err != nil {
		return nil, err
	}
	args.StartTime, args.EndTime, err = normalizeTimeRange(args.StartTime, args.EndTime)
	if err != nil {
		return nil, err
	}

	logic := aiwindtimeseriesservicelogic.NewCompareTrendLogic(ctx, e.svcCtx)
	resp, err := logic.CompareTrend(&pb.WindTrendCompareReq{
		FarmCode:       farmCode,
		TowerCode:      strings.TrimSpace(args.TowerCode),
		DeviceCode:     strings.TrimSpace(args.DeviceCode),
		DeviceTypeCode: args.DeviceTypeCode,
		Field:          args.Field,
		StartTime:      args.StartTime,
		EndTime:        args.EndTime,
		IndexId:        args.IndexID,
		RadarDistanceM: args.RadarDistanceM,
		UserId:         req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &toolResult{
		ResultJSON: marshalJSON(map[string]any{
			"summary": resp.Summary,
			"message": resp.Message,
		}),
		EvidenceJSON: resp.EvidenceJson,
		Message:      resp.Message,
	}, nil
}
