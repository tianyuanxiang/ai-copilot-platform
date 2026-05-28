package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// FieldThreshold 单字段阈值配置
type FieldThreshold struct {
	Upper             *float64 `json:"upper,omitempty"`
	Lower             *float64 `json:"lower,omitempty"`
	ContinuousSeconds int64    `json:"continuousSeconds,omitempty"`
}

// RadarThresholdConfig 雷达阈值配置，支持字段级 + index 覆盖
type RadarThresholdConfig struct {
	Fields        map[string]FieldThreshold            `json:",inline"`
	IndexOverrides map[string]map[string]FieldThreshold `json:"indexOverrides,omitempty"`
}

// ThresholdConfig 阈值配置结果
type ThresholdConfig struct {
	Status string
	// 普通设备：field -> FieldThreshold
	Fields map[string]FieldThreshold
	// 雷达：字段级 + index 覆盖
	RadarConfig *RadarThresholdConfig
	// 原始 JSON（用于 evidence 输出）
	RawJSON string
}

// ThresholdResolver 阈值解析器接口
type ThresholdResolver interface {
	// ResolveThreshold 读取设备阈值配置
	// 优先从 wind_device_meta.threshold_config 读取，回退 wind_device_type.threshold_config
	ResolveThreshold(ctx context.Context, deviceTypeCode string, towerCode string) (ThresholdConfig, error)
}

type thresholdResolver struct {
	db *sql.DB
}

// NewThresholdResolver 创建阈值解析器
func NewThresholdResolver(db *sql.DB) ThresholdResolver {
	return &thresholdResolver{db: db}
}

// ResolveThreshold 读取阈值配置，优先 wind_device_meta，回退 wind_device_type
func (r *thresholdResolver) ResolveThreshold(ctx context.Context, deviceTypeCode string, towerCode string) (ThresholdConfig, error) {
	if r.db == nil {
		return ThresholdConfig{Status: "not_configured"}, nil
	}

	// 优先查 wind_device_meta 的 threshold_config
	raw, err := r.queryDeviceMetaThreshold(ctx, deviceTypeCode, towerCode)
	if err != nil || raw == "" {
		// 回退到 wind_device_type 的 threshold_config
		raw, err = r.queryDeviceTypeThreshold(ctx, deviceTypeCode)
		if err != nil || raw == "" {
			return ThresholdConfig{Status: "not_configured"}, nil
		}
	}

	return parseThresholdConfig(raw, deviceTypeCode)
}

// queryDeviceMetaThreshold 从 wind_device_meta 查阈值
func (r *thresholdResolver) queryDeviceMetaThreshold(ctx context.Context, deviceTypeCode, towerCode string) (string, error) {
	query := `
		SELECT wdm.threshold_config
		FROM wind_device_meta wdm
		JOIN wind_device wd ON wd.id = wdm.wind_device_id
		JOIN wind_device_type wdt ON wdt.id = wd.device_type_id
		WHERE wdt.code = $1
	`
	args := []interface{}{deviceTypeCode}
	if towerCode != "" {
		query += ` AND wd.tower_id = (SELECT id FROM wind_tower WHERE code = $2 LIMIT 1)`
		args = append(args, towerCode)
	}
	query += ` AND wdm.threshold_config IS NOT NULL AND wdm.threshold_config <> '' LIMIT 1`

	var raw sql.NullString
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&raw)
	if err != nil {
		return "", nil
	}
	if !raw.Valid {
		return "", nil
	}
	return raw.String, nil
}

// queryDeviceTypeThreshold 从 wind_device_type 查阈值
func (r *thresholdResolver) queryDeviceTypeThreshold(ctx context.Context, deviceTypeCode string) (string, error) {
	query := `
		SELECT threshold_config
		FROM wind_device_type
		WHERE code = $1
		  AND threshold_config IS NOT NULL AND threshold_config <> ''
		LIMIT 1
	`
	var raw sql.NullString
	err := r.db.QueryRowContext(ctx, query, deviceTypeCode).Scan(&raw)
	if err != nil {
		return "", nil
	}
	if !raw.Valid {
		return "", nil
	}
	return raw.String, nil
}

// parseThresholdConfig 解析阈值配置 JSON
func parseThresholdConfig(raw string, deviceTypeCode string) (ThresholdConfig, error) {
	cfg := ThresholdConfig{
		Status:  "configured",
		RawJSON: raw,
	}

	if deviceTypeCode == "WPR" {
		// 雷达阈值格式：{"rws":{"upper":25},"indexOverrides":{"3":{"rws":{"upper":22}}}}
		var radar RadarThresholdConfig
		if err := json.Unmarshal([]byte(raw), &radar); err != nil {
			return ThresholdConfig{Status: "parse_error"}, fmt.Errorf("parse radar threshold config: %w", err)
		}
		cfg.RadarConfig = &radar
	} else {
		// 普通传感器格式：{"accel":{"upper":10,"lower":-10,"continuousSeconds":60}}
		var fields map[string]FieldThreshold
		if err := json.Unmarshal([]byte(raw), &fields); err != nil {
			return ThresholdConfig{Status: "parse_error"}, fmt.Errorf("parse threshold config: %w", err)
		}
		cfg.Fields = fields
	}

	return cfg, nil
}

// CheckContinuousExceedance 检测持续超限，返回最长超限秒数和超限次数
// points 必须按时间顺序排列（升序）
func CheckContinuousExceedance(points []TimePoint, threshold FieldThreshold, intervalMs int64) (maxContinuousSec float64, exceedCount int) {
	if len(points) == 0 {
		return 0, 0
	}
	if threshold.Upper == nil && threshold.Lower == nil {
		return 0, 0
	}

	intervalSec := float64(intervalMs) / 1000.0
	currentSec := 0.0
	inExceedance := false

	for _, p := range points {
		exceeded := false
		if threshold.Upper != nil && p.Value > *threshold.Upper {
			exceeded = true
		}
		if threshold.Lower != nil && p.Value < *threshold.Lower {
			exceeded = true
		}

		if exceeded {
			if !inExceedance {
				inExceedance = true
				exceedCount++
				currentSec = intervalSec
			} else {
				currentSec += intervalSec
			}
			if currentSec > maxContinuousSec {
				maxContinuousSec = currentSec
			}
		} else {
			inExceedance = false
			currentSec = 0
		}
	}

	return maxContinuousSec, exceedCount
}

// TimePoint 时间点值，用于超限检测
type TimePoint struct {
	Ts    string
	Value float64
}

// GetRiskLevel 根据超限情况返回风险等级
func GetRiskLevel(exceedCount int, maxContinuousSec float64, threshold FieldThreshold) string {
	if exceedCount == 0 {
		return "normal"
	}
	continuousLimit := float64(threshold.ContinuousSeconds)
	if continuousLimit <= 0 {
		continuousLimit = 60
	}
	if maxContinuousSec >= continuousLimit*3 {
		return "critical"
	}
	if maxContinuousSec >= continuousLimit {
		return "high"
	}
	if exceedCount >= 5 {
		return "warning"
	}
	return "warning"
}
