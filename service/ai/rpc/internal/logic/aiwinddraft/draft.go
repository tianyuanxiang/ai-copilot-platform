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

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/xerrors"
)

// AlarmSamplePool 告警分层抽样的样本池大小
const AlarmSamplePool = 500

// BuildAlarmEvidence 构建告警分析 evidence，使用 TDengine 服务端聚合查询 + 分层抽样。
func BuildAlarmEvidence(ctx context.Context, svcCtx *svc.ServiceContext, farmCode, towerCode, alarmCode,
	startTime, endTime string, status int64, hasStatus bool) (map[string]any, error) {

	logger := logx.WithContext(ctx)
	database := svcCtx.WindFarmModel.FarmDatabase(ctx, farmCode)
	logger.Infof("[BuildAlarmEvidence] 开始构建告警证据 farm_code=%s tower_code=%s alarm_code=%s time=[%s, %s]",
		farmCode, towerCode, alarmCode, startTime, endTime)

	whereParts := append(model.TimeWhere(startTime, endTime), model.DeviceWhere(towerCode, "")...)
	if strings.TrimSpace(alarmCode) != "" {
		whereParts = append(whereParts, fmt.Sprintf("alarm_code=%d", parseInt(alarmCode)))
	}
	if hasStatus {
		whereParts = append(whereParts, fmt.Sprintf("status=%d", status))
	}
	where := model.JoinWhere(whereParts)

	// 选择时间聚合粒度
	granularity := selectAlarmGranularity(startTime, endTime)

	// 查询 1: 一条 SQL 拿到 total + first_ts + last_ts
	aggResult, err := svcCtx.TdengineModel.QueryAlarmAggregates(ctx, database, where)
	if err != nil {
		logger.Errorf("[BuildAlarmEvidence] QueryAlarmAggregates 失败 database=%s where=%s err=%v", database, where, err)
		return nil, xerrors.Errorf("failed to QueryAlarmAggregates: %v", err)
	}

	total := parseInt(aggResult["total"])
	firstTs := aggResult["first_ts"]
	lastTs := aggResult["last_ts"]

	// 查询 2-5: GROUP BY 分布统计
	levelRows, err := svcCtx.TdengineModel.QueryAlarmGroupBy(ctx, database, "alarm_level", where, 0)
	if err != nil {
		logger.Errorf("[BuildAlarmEvidence] QueryAlarmGroupBy(alarm_level) 失败 database=%s err=%v", database, err)
		return nil, err
	}
	statusRows, err := svcCtx.TdengineModel.QueryAlarmGroupBy(ctx, database, "status", where, 0)
	if err != nil {
		logger.Errorf("[BuildAlarmEvidence] QueryAlarmGroupBy(status) 失败 database=%s err=%v", database, err)
		return nil, err
	}
	alarmCodeRows, err := svcCtx.TdengineModel.QueryAlarmGroupBy(ctx, database, "alarm_code", where, 20)
	if err != nil {
		logger.Errorf("[BuildAlarmEvidence] QueryAlarmGroupBy(alarm_code) 失败 database=%s err=%v", database, err)
		return nil, err
	}
	towerRows, err := svcCtx.TdengineModel.QueryAlarmGroupBy(ctx, database, "tower_id", where, 20)
	if err != nil {
		logger.Errorf("[BuildAlarmEvidence] QueryAlarmGroupBy(tower_id) 失败 database=%s err=%v", database, err)
		return nil, err
	}

	// 查询 6: 时间桶分布
	bucketRows, err := svcCtx.TdengineModel.QueryAlarmTimeBuckets(ctx, database, where, granularity)
	if err != nil {
		logger.Errorf("[BuildAlarmEvidence] QueryAlarmTimeBuckets 失败 database=%s granularity=%s err=%v", database, granularity, err)
		return nil, err
	}

	// 构建聚合结果
	byLevel := groupRowsToMap(levelRows)
	byStatus := groupRowsToMap(statusRows)
	byAlarmCodeTop := groupRowsToList(alarmCodeRows)
	byTowerTop := groupRowsToList(towerRows)
	byTimeBucket := timeBucketsToList(bucketRows)
	peakBucket := findPeakBucket(byTimeBucket)

	// 查询 7: 拉取样本池用于分层抽样
	poolSize := int64(AlarmSamplePool)
	if total < poolSize {
		poolSize = total
	}
	poolRows, err := svcCtx.TdengineModel.QueryAlarms(ctx, database, where, 1, poolSize)
	if err != nil {
		logger.Errorf("[BuildAlarmEvidence] QueryAlarms(样本池) 失败 database=%s poolSize=%d err=%v", database, poolSize, err)
		return nil, err
	}

	// 提取 top alarm_code 和 tower_code 列表用于抽样
	topAlarmCodes := extractTopKeys(alarmCodeRows, 10)
	topTowers := extractTopKeys(towerRows, 5)
	peakBucketStart := ""
	if peakBucket != nil {
		peakBucketStart = peakBucket["bucket_start"].(string)
	}
	// 分层抽样
	samples := selectAlarmSamples(poolRows, topAlarmCodes, topTowers, peakBucketStart)

	logger.Infof("[BuildAlarmEvidence] 构建完成 total=%d level_counts=%v sampled=%d",
		total, byLevel, len(samples))

	return map[string]any{
		"source":            "tdengine.alarm",
		"scaffold":          !svcCtx.TdengineModel.IsConfigured(),
		"farm_code":         farmCode,
		"tower_code":        towerCode,
		"alarm_code":        alarmCode,
		"database":          database,
		"where":             where,
		"start_time":        startTime,
		"end_time":          endTime,
		"granularity":       granularity,
		"total":             total,
		"sampled":           len(samples),
		"truncated":         total > int64(len(samples)),
		"by_level":          byLevel,
		"by_status":         byStatus,
		"by_alarm_code_top": byAlarmCodeTop,
		"by_tower_top":      byTowerTop,
		"time_bucket_count": len(byTimeBucket),
		"first_ts":          firstTs,
		"last_ts":           lastTs,
		"peak_bucket":       peakBucket,
		"samples":           samples,
	}, nil
}

// selectAlarmGranularity 根据时间范围自动选择聚合粒度。
func selectAlarmGranularity(startTime, endTime string) string {
	layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05", time.RFC3339}
	var start, end time.Time
	for _, layout := range layouts {
		if t, err := time.Parse(layout, startTime); err == nil {
			start = t
			break
		}
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, endTime); err == nil {
			end = t
			break
		}
	}
	if start.IsZero() || end.IsZero() || !end.After(start) {
		return "1h"
	}
	duration := end.Sub(start)
	switch {
	case duration <= 24*time.Hour:
		return "1m"
	case duration <= 7*24*time.Hour:
		return "1h"
	case duration <= 90*24*time.Hour:
		return "1d"
	default:
		return "1w"
	}
}

// groupRowsToMap 将 GROUP BY 查询结果转为 map[string]int64。
func groupRowsToMap(rows []map[string]string) map[string]int64 {
	m := make(map[string]int64, len(rows))
	for _, row := range rows {
		key := ""
		cnt := int64(0)
		for k, v := range row {
			if k == "cnt" {
				cnt = parseInt(v)
			} else {
				key = v
			}
		}
		if key != "" {
			m[key] = cnt
		}
	}
	return m
}

// groupRowsToList 将 GROUP BY 查询结果转为 []map[string]any [{key,count}]。
func groupRowsToList(rows []map[string]string) []map[string]any {
	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		key := ""
		cnt := int64(0)
		for k, v := range row {
			if k == "cnt" {
				cnt = parseInt(v)
			} else {
				key = v
			}
		}
		list = append(list, map[string]any{"key": key, "count": cnt})
	}
	return list
}

// timeBucketsToList 将 INTERVAL 查询结果转为 []map[string]any [{bucket_start,count}]。
func timeBucketsToList(rows []map[string]string) []map[string]any {
	list := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		bucket := map[string]any{}
		for k, v := range row {
			if k == "_wstart" {
				bucket["bucket_start"] = v
			} else if k == "cnt" {
				bucket["count"] = parseInt(v)
			}
		}
		list = append(list, bucket)
	}
	return list
}

// findPeakBucket 从时间桶列表中找到告警最多的桶。
func findPeakBucket(buckets []map[string]any) map[string]any {
	if len(buckets) == 0 {
		return nil
	}
	var peak map[string]any
	var maxCount int64
	for _, b := range buckets {
		if cnt, ok := b["count"].(int64); ok && cnt > maxCount {
			maxCount = cnt
			peak = b
		}
	}
	return peak
}

// extractTopKeys 从 GROUP BY 结果中提取前 n 个 key。
func extractTopKeys(rows []map[string]string, n int) []string {
	keys := make([]string, 0, n)
	for i, row := range rows {
		if i >= n {
			break
		}
		for k, v := range row {
			if k != "cnt" {
				keys = append(keys, v)
			}
		}
	}
	return keys
}

// selectAlarmSamples 从样本池中分层抽样，返回带 sample_reason 的样本列表。
func selectAlarmSamples(rows []map[string]string, topAlarmCodes []string, topTowers []string, peakBucketStart string) []map[string]any {
	if len(rows) == 0 {
		return nil
	}
	// 索引 -> 已选原因
	selected := make(map[int]bool)
	reasons := make(map[int][]string)

	addReason := func(idx int, reason string) {
		if idx < 0 || idx >= len(rows) {
			return
		}
		if !selected[idx] {
			selected[idx] = true
			reasons[idx] = []string{reason}
		} else {
			// 避免重复 reason
			for _, r := range reasons[idx] {
				if r == reason {
					return
				}
			}
			reasons[idx] = append(reasons[idx], reason)
		}
	}

	// 1. 高等级告警 (level >= 3)
	highCount := 0
	for i, row := range rows {
		if highCount >= 20 {
			break
		}
		if parseInt(row["alarm_level"]) >= 3 {
			addReason(i, "high_level")
			highCount++
		}
	}

	// 2. Top alarm_code 每类 2 条
	topCodeSet := make(map[string]bool, len(topAlarmCodes))
	for _, code := range topAlarmCodes {
		topCodeSet[code] = true
	}
	codeCount := make(map[string]int)
	for i, row := range rows {
		code := row["alarm_code"]
		if !topCodeSet[code] {
			continue
		}
		if codeCount[code] >= 2 {
			continue
		}
		addReason(i, "top_alarm_code")
		codeCount[code]++
	}

	// 3. 峰值时间桶 10 条
	if peakBucketStart != "" {
		peakCount := 0
		peakPrefix := peakBucketStart[:min(10, len(peakBucketStart))] // 取日期部分做前缀匹配
		for i, row := range rows {
			if peakCount >= 10 {
				break
			}
			if strings.HasPrefix(row["ts"], peakPrefix) {
				addReason(i, "peak_bucket")
				peakCount++
			}
		}
	}

	// 4. 边界样本: 最早 5 + 最新 5
	for i := 0; i < 5 && i < len(rows); i++ {
		addReason(i, "boundary")
	}
	for i := len(rows) - 1; i >= len(rows)-5 && i >= 0; i-- {
		addReason(i, "boundary")
	}

	// 5. Top tower 每台 2 条
	topTowerSet := make(map[string]bool, len(topTowers))
	for _, t := range topTowers {
		topTowerSet[t] = true
	}
	towerCount := make(map[string]int)
	for i, row := range rows {
		tid := row["tower_id"]
		if !topTowerSet[tid] {
			continue
		}
		if towerCount[tid] >= 2 {
			continue
		}
		addReason(i, "top_tower")
		towerCount[tid]++
	}

	// 构建最终样本列表
	samples := make([]map[string]any, 0, len(selected))
	for idx := range selected {
		row := rows[idx]
		samples = append(samples, map[string]any{
			"ts":            row["ts"],
			"tower_code":    row["tower_id"],
			"alarm_level":   parseInt(row["alarm_level"]),
			"alarm_code":    parseInt(row["alarm_code"]),
			"status":        parseInt(row["status"]),
			"sample_reason": reasons[idx],
		})
	}
	return samples
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
			"total", "query_mode", "risk", "point_count",
			"granularity", "start_time", "end_time",
			"by_level", "by_status", "by_alarm_code_top", "by_tower_top",
			"time_bucket_count", "first_ts", "last_ts", "peak_bucket",
			"sampled"} {
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
		UserId:    userID,
		TraceId:   traceID,
		ToolName:  toolName,
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
		"start_time", "end_time", "device_type_code", "field", "priority", "granularity"} {
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
		"truncated":      true,
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
