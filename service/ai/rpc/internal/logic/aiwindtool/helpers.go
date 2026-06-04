// Package aiwindtool 的 helpers 文件收纳轻量参数校验和 JSON 辅助函数。
// 这些规则只保护工具入口，不替代已有业务 logic 内部的校验。
package aiwindtool

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	defaultPageSize = int64(500)
	maxPageSize     = int64(1000)
	maxFieldCount   = 8
	defaultTopK     = int64(5)
	maxTopK         = int64(10)
	maxQueryDays    = 31
	toolTimeLayout  = "2006-01-02 15:04:05"
)

// decodeArgs 将 Agent 提供的 JSON 参数解码到目标结构体。
// 空字符串等价于空 JSON 对象，方便无参数工具直接调用。
func decodeArgs(raw string, target any) error {
	if strings.TrimSpace(raw) == "" {
		raw = "{}"
	}
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		return fmt.Errorf("argumentsJson 不是合法 JSON: %w", err)
	}
	return nil
}

// marshalJSON 将内部结果编码为 JSON。
// 工具返回值必须始终可交给 Agent 继续推理，编码异常时返回空对象。
func marshalJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// normalizeFarmCode 统一清理风场编号并转成大写。
func normalizeFarmCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

// normalizePage 统一分页参数：页码默认 1，每页默认 50，最大 100。
func normalizePage(page, pageSize int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

// normalizeTopK 统一 SOP 检索条数：默认 5，最大 10。
func normalizeTopK(value int64) int64 {
	if value <= 0 {
		return defaultTopK
	}
	if value > maxTopK {
		return maxTopK
	}
	return value
}

// normalizeTimeRange 统一查询时间范围。
// 两端都未填写时默认最近 24 小时；只填写一端、结束时间不晚于开始时间、
// 或范围超过 31 天时返回错误，防止 Agent 发起无边界查询。
func normalizeTimeRange(startTime, endTime string) (string, string, error) {
	startTime = strings.TrimSpace(startTime)
	endTime = strings.TrimSpace(endTime)
	if startTime == "" && endTime == "" {
		end := time.Now()
		return end.Add(-24 * time.Hour).Format(toolTimeLayout), end.Format(toolTimeLayout), nil
	}
	if startTime == "" || endTime == "" {
		return "", "", fmt.Errorf("startTime 和 endTime 必须同时提供")
	}

	start, err := parseToolTime(startTime)
	if err != nil {
		return "", "", fmt.Errorf("startTime 格式错误: %w", err)
	}
	end, err := parseToolTime(endTime)
	if err != nil {
		return "", "", fmt.Errorf("endTime 格式错误: %w", err)
	}
	if !end.After(start) {
		return "", "", fmt.Errorf("endTime 必须晚于 startTime")
	}
	if end.Sub(start) > maxQueryDays*24*time.Hour {
		return "", "", fmt.Errorf("查询时间范围不能超过 %d 天", maxQueryDays)
	}
	return start.Format(toolTimeLayout), end.Format(toolTimeLayout), nil
}

// validateFields 限制一次工具调用最多查询八个测点字段。
func validateFields(fields []string) error {
	if len(fields) > maxFieldCount {
		return fmt.Errorf("field 最多允许 %d 个", maxFieldCount)
	}
	return nil
}

func parseToolTime(value string) (time.Time, error) {
	for _, layout := range []string{
		toolTimeLayout,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000",
		time.RFC3339,
	} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("应使用 yyyy-MM-dd HH:mm:ss 或 RFC3339")
}

func requireFarmCode(value string) (string, error) {
	farmCode := normalizeFarmCode(value)
	if farmCode == "" {
		return "", fmt.Errorf("farmCode 不能为空")
	}
	return farmCode, nil
}
