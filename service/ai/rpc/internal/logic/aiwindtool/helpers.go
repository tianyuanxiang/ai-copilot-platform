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

type normalizedSensorToolArgs struct {
	FarmCode       string
	TowerCode      string
	DeviceCode     string
	DeviceTypeCode string
	Fields         []string
	StartTime      string
	EndTime        string
	Page           int64
	PageSize       int64
	IndexID        int64
	RadarDistanceM int64
}

func decodeArgs(raw string, target any) error {
	if strings.TrimSpace(raw) == "" {
		raw = "{}"
	}
	if err := json.Unmarshal([]byte(raw), target); err != nil {
		return fmt.Errorf("argumentsJson 不是合法 JSON: %w", err)
	}
	return nil
}

func marshalJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func normalizeFarmCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

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

func normalizeTopK(value int64) int64 {
	if value <= 0 {
		return defaultTopK
	}
	if value > maxTopK {
		return maxTopK
	}
	return value
}

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

func normalizeSensorToolArgs(farmCode, towerCode, deviceCode, deviceTypeCode string, fields []string,
	startTime, endTime string, indexID, radarDistanceM, page, pageSize int64, withPage bool) (normalizedSensorToolArgs, error) {
	normalizedFarmCode, err := requireFarmCode(farmCode)
	if err != nil {
		return normalizedSensorToolArgs{}, err
	}
	deviceTypeCode = strings.ToUpper(strings.TrimSpace(deviceTypeCode))
	if deviceTypeCode == "" {
		return normalizedSensorToolArgs{}, fmt.Errorf("deviceTypeCode 不能为空")
	}
	if err := validateFields(fields); err != nil {
		return normalizedSensorToolArgs{}, err
	}
	startTime, endTime, err = normalizeTimeRange(startTime, endTime)
	if err != nil {
		return normalizedSensorToolArgs{}, err
	}
	if withPage {
		page, pageSize = normalizePage(page, pageSize)
	}
	return normalizedSensorToolArgs{
		FarmCode:       normalizedFarmCode,
		TowerCode:      strings.TrimSpace(towerCode),
		DeviceCode:     strings.TrimSpace(deviceCode),
		DeviceTypeCode: deviceTypeCode,
		Fields:         fields,
		StartTime:      startTime,
		EndTime:        endTime,
		Page:           page,
		PageSize:       pageSize,
		IndexID:        indexID,
		RadarDistanceM: radarDistanceM,
	}, nil
}
