package aiwinddraft

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
)

// AlarmSampleLimit 告警证据最大样本条数，超出部分只保留聚合统计
const AlarmSampleLimit = 50

func BuildAlarmEvidence(ctx context.Context, svcCtx *svc.ServiceContext, farmCode, towerCode, alarmCode,
	startTime, endTime string, status int64, hasStatus bool) (map[string]any, error) {

	database := svcCtx.WindFarmModel.FarmDatabase(ctx, farmCode)

	whereParts := append(model.TimeWhere(startTime, endTime), model.DeviceWhere(towerCode, "")...)
	if strings.TrimSpace(alarmCode) != "" {
		whereParts = append(whereParts, fmt.Sprintf("alarm_code=%d", parseInt(alarmCode)))
	}
	if hasStatus {
		whereParts = append(whereParts, fmt.Sprintf("status=%d", status))
	}
	where := model.JoinWhere(whereParts)

	total, err := svcCtx.TdengineModel.QueryCount(ctx, database, "alarm", where)
	if err != nil {
		return nil, err
	}

	// 只查询样本上限条记录，不再拉取全量
	sampleLimit := int64(AlarmSampleLimit)
	if total < sampleLimit {
		sampleLimit = total
	}
	rows, err := svcCtx.TdengineModel.QueryAlarms(ctx, database, where, 1, sampleLimit)
	if err != nil {
		return nil, err
	}

	truncated := total > AlarmSampleLimit

	// 聚合统计：告警等级、状态、风机分布
	levelCounts := map[string]int64{}
	statusCounts := map[string]int64{}
	towerCounts := map[string]int64{}
	for _, row := range rows {
		levelKey := fmt.Sprintf("level_%s", row["alarm_level"])
		levelCounts[levelKey]++

		statusKey := fmt.Sprintf("status_%s", row["status"])
		statusCounts[statusKey]++

		towerKey := row["tower_id"]
		if towerKey == "" {
			towerKey = towerCode
		}
		towerCounts[towerKey]++
	}

	// 如果数据被截断，用全量 total 计算更准确的分布需二次查询；
	// 此处基于样本估算，在返回中标注 truncated=true 表示统计基于样本而非全量
	if truncated {
		levelCounts = nil
		statusCounts = nil
		towerCounts = nil
	}

	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		list = append(list, map[string]any{
			"ts":             row["ts"],
			"farm_code":      farmCode,
			"tower_code":     row["tower_id"],
			"device_channel": parseInt(row["device_channel"]),
			"device_type":    parseInt(row["device_type"]),
			"alarm_location": row["alarm_location"],
			"alarm_level":    parseInt(row["alarm_level"]),
			"alarm_code":     parseInt(row["alarm_code"]),
			"alarm_value":    row["alarm_value"],
			"status":         parseInt(row["status"]),
		})
	}

	evidence := map[string]any{
		"source":           "tdengine.alarm",
		"scaffold":         !svcCtx.TdengineModel.IsConfigured(),
		"farm_code":        farmCode,
		"tower_code":       towerCode,
		"alarm_code":       alarmCode,
		"database":         database,
		"where":            where,
		"alarm_count":      total,
		"returned_records": len(list),
		"truncated":        truncated,
		"time_range": map[string]string{
			"start": startTime,
			"end":   endTime,
		},
		"list": list,
	}
	if levelCounts != nil {
		evidence["level_counts"] = levelCounts
	}
	if statusCounts != nil {
		evidence["status_counts"] = statusCounts
	}
	if towerCounts != nil {
		evidence["tower_counts"] = towerCounts
	}

	return evidence, nil
}

func MergeEvidence(inputEvidence string, autoEvidence ...map[string]any) ([]map[string]any, string) {
	items := make([]map[string]any, 0, len(autoEvidence)+1)
	if strings.TrimSpace(inputEvidence) != "" {
		var parsed any
		if err := json.Unmarshal([]byte(inputEvidence), &parsed); err != nil {
			items = append(items, map[string]any{"type": "invalid_input_evidence_json", "raw": inputEvidence})
		} else {
			switch value := parsed.(type) {
			case []any:
				for _, item := range value {
					if mapped, ok := item.(map[string]any); ok {
						items = append(items, mapped)
					}
				}
			case map[string]any:
				items = append(items, value)
			}
		}
	}
	for _, item := range autoEvidence {
		if item != nil {
			items = append(items, item)
		}
	}
	return items, model.WindEvidenceJSON(items)
}

// BuildEvidenceSummary 从完整 evidence 中提取轻量摘要，用于数据库存储和 API 响应。
// 保留聚合指标和元数据，去除原始 list 数据。
func BuildEvidenceSummary(evidenceJSON string) string {
	if evidenceJSON == "" {
		return ""
	}
	var items []map[string]any
	if err := json.Unmarshal([]byte(evidenceJSON), &items); err != nil {
		// 解析失败时返回原始内容（兼容旧格式）
		return evidenceJSON
	}

	summaries := make([]map[string]any, 0, len(items))
	for _, item := range items {
		summary := map[string]any{}

		// 保留来源和元数据字段
		for _, key := range []string{"source", "scaffold", "farm_code", "tower_code", "alarm_code",
			"database", "where", "device_type_code", "stable", "fields",
			"farmCode", "towerCode", "deviceTypeCode"} {
			if v, ok := item[key]; ok {
				summary[key] = v
			}
		}

		// 保留聚合统计字段
		for _, key := range []string{"alarm_count", "returned_records", "truncated",
			"level_counts", "status_counts", "tower_counts", "time_range",
			"total", "query_mode", "risk", "point_count"} {
			if v, ok := item[key]; ok {
				summary[key] = v
			}
		}

		// 保留 trend 对比指标
		for _, key := range []string{"field_stats", "trend_signal", "risk_signal",
			"data_quality", "threshold_config", "agent_hints", "granularity"} {
			if v, ok := item[key]; ok {
				summary[key] = v
			}
		}

		summaries = append(summaries, summary)
	}

	data, err := json.Marshal(summaries)
	if err != nil {
		return evidenceJSON
	}
	return string(data)
}

func DraftContent(resp *engine.WindDraftResponse) string {
	if resp == nil {
		return "{}"
	}
	data, err := json.Marshal(map[string]any{
		"status":          resp.Status,
		"summary":         resp.Summary,
		"metrics":         resp.Metrics,
		"sections":        resp.Sections,
		"recommendations": resp.Recommendations,
		"todo":            resp.Todo,
		"message":         resp.Message,
	})
	if err != nil {
		return "{}"
	}
	return string(data)
}

// ToolCallLogMaxBytes 工具调用日志单字段最大字节数，超限截断并标记
const ToolCallLogMaxBytes = 32768

func WriteToolCallLog(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, traceID string, toolName string, args any, result any, startedAt time.Time, status string, errorMsg string) {
	if svcCtx == nil || svcCtx.AiToolCallLogModel == nil {
		return
	}

	// 只记录轻量摘要，不写完整 payload 或 draft 内容
	logArgs := slimToolCallArgs(toolName, args)
	logResult := slimToolCallResult(result, status, startedAt)

	_, err := svcCtx.AiToolCallLogModel.Insert(ctx, &model.AiToolCallLog{
		UserId:   userID,
		TraceId:  traceID,
		ToolName: toolName,
		Arguments: truncateJSONBytes(model.WindEvidenceJSON(logArgs), ToolCallLogMaxBytes),
		Result:    truncateJSONBytes(model.WindEvidenceJSON(logResult), ToolCallLogMaxBytes),
		Status:    status,
		ErrorMsg:  truncateRunes(errorMsg, 1000),
		LatencyMs: time.Since(startedAt).Milliseconds(),
	})
	if err != nil {
		fmt.Printf("write ai_tool_call_log failed: %v\n", err)
	}
}

// slimToolCallArgs 从原始参数中提取轻量摘要信息
func slimToolCallArgs(toolName string, args any) map[string]any {
	slim := map[string]any{
		"tool_name": toolName,
	}

	// WindDraftRequest 是结构体，通过 JSON 序列化/反序列化提取关键字段
	raw, err := json.Marshal(args)
	if err != nil {
		return slim
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return slim
	}

	for _, key := range []string{"farm_code", "tower_code", "alarm_code", "report_type",
		"start_time", "end_time", "device_type_code", "field", "priority"} {
		if v, exists := m[key]; exists {
			slim[key] = v
		}
	}
	// 记录 evidence 数据量
	if evidence, ok := m["evidence"].([]any); ok {
		slim["evidence_item_count"] = len(evidence)
	}
	if evidenceJSON, ok := m["evidence_json"].(string); ok {
		slim["evidence_json_bytes"] = len(evidenceJSON)
	}
	return slim
}

// slimToolCallResult 从原始结果中提取轻量摘要信息
func slimToolCallResult(result any, status string, startedAt time.Time) map[string]any {
	slim := map[string]any{
		"status":     status,
		"latency_ms": time.Since(startedAt).Milliseconds(),
	}

	// WindDraftResponse 是结构体，通过 JSON 序列化/反序列化提取关键字段
	raw, err := json.Marshal(result)
	if err != nil {
		return slim
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return slim
	}

	for _, key := range []string{"title", "summary", "status", "message"} {
		if v, exists := m[key]; exists {
			slim[key] = v
		}
	}
	if metrics, ok := m["metrics"]; ok {
		slim["metrics"] = metrics
	}
	return slim
}

// truncateJSONBytes 当 JSON 字符串超过 maxBytes 时，截断并替换为截断标记
func truncateJSONBytes(jsonStr string, maxBytes int) string {
	if maxBytes <= 0 || len(jsonStr) <= maxBytes {
		return jsonStr
	}
	truncated, _ := json.Marshal(map[string]any{
		"truncated":     true,
		"original_bytes": len(jsonStr),
	})
	return string(truncated)
}

func ParseNullableTime(value string) sql.NullTime {
	value = strings.TrimSpace(value)
	if value == "" {
		return sql.NullTime{}
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return sql.NullTime{Time: parsed, Valid: true}
		}
	}
	return sql.NullTime{}
}

func UserIDString(userID int64) string {
	if userID <= 0 {
		return "0"
	}
	return strconv.FormatInt(userID, 10)
}

func parseInt(value string) int64 {
	var n int64
	fmt.Sscan(value, &n)
	return n
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 || utf8.RuneCountInString(value) <= limit {
		return value
	}
	runes := []rune(value)
	return string(runes[:limit])
}
