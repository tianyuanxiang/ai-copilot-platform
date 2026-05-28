package model

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// AggStat 单个字段的聚合统计结果
type AggStat struct {
	Count   int64
	Min     float64
	Max     float64
	Avg     float64
	Stddev  float64
	First   float64
	Last    float64
	FirstTs string
	LastTs  string
}

// FieldStats 单个字段的完整趋势统计
type FieldStats struct {
	Count         int64   `json:"count"`
	ExpectedCount int64   `json:"expectedCount"`
	MissingRate   float64 `json:"missingRate"`
	Min           float64 `json:"min"`
	Max           float64 `json:"max"`
	Avg           float64 `json:"avg"`
	Stddev        float64 `json:"stddev"`
	Amplitude     float64 `json:"amplitude"`
	First         float64 `json:"first"`
	Last          float64 `json:"last"`
	Change        float64 `json:"change"`
	ChangePct     float64 `json:"changePct"`
	Trend         string  `json:"trend"`
	LatestValue   float64 `json:"latestValue"`
	LatestTs      string  `json:"latestTs"`
	Slope         float64 `json:"slope"`
}

// DataQuality 数据质量指标
type DataQuality struct {
	MissingRate      float64 `json:"missingRate"`
	LastDataTime     string  `json:"lastDataTime"`
	FreshnessSeconds int64   `json:"freshnessSeconds"`
	SuspectedOffline bool    `json:"suspectedOffline"`
}

// TrendSignal 趋势信号
type TrendSignal struct {
	Direction string  `json:"direction"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"changePct"`
	Slope     float64 `json:"slope"`
}

// RiskSignal 风险信号
type RiskSignal struct {
	ExceedCount      int    `json:"exceedCount"`
	MaxContinuousSec int64  `json:"maxContinuousSec"`
	Exceeded         bool   `json:"exceeded"`
	Level            string `json:"level"`
}

// CompareTrendEvidence compare 接口 evidence 完整结构
type CompareTrendEvidence struct {
	Source         string                         `json:"source"`
	Scaffold       bool                           `json:"scaffold"`
	FarmCode       string                         `json:"farmCode"`
	TowerCode      string                         `json:"towerCode"`
	Database       string                         `json:"database"`
	Stable         string                         `json:"stable"`
	Fields         []string                       `json:"fields"`
	SamplingPolicy map[string]interface{}         `json:"samplingPolicy"`
	DataQuality    *DataQuality                   `json:"dataQuality,omitempty"`
	FieldStats     map[string]*FieldStats         `json:"fieldStats,omitempty"`
	RadarStats     map[int]map[string]*FieldStats `json:"radarStats,omitempty"`
	TrendSignal    *TrendSignal                   `json:"trendSignal,omitempty"`
	RiskSignal     *RiskSignal                    `json:"riskSignal,omitempty"`
	ThresholdCfg   map[string]interface{}         `json:"thresholdConfig,omitempty"`
	AgentHints     []string                       `json:"agentHints"`
}

// CalcFieldStats 根据 AggStat 和期望数量计算 FieldStats
func CalcFieldStats(agg AggStat, expectedCount int64, durationSeconds float64) *FieldStats {
	fs := &FieldStats{
		Count:         agg.Count,
		ExpectedCount: expectedCount,
		Min:           agg.Min,
		Max:           agg.Max,
		Avg:           agg.Avg,
		Stddev:        agg.Stddev,
		Amplitude:     agg.Max - agg.Min,
		First:         agg.First,
		Last:          agg.Last,
		Change:        agg.Last - agg.First,
		LatestValue:   agg.Last,
		LatestTs:      agg.LastTs,
	}

	// 缺测率
	if expectedCount > 0 {
		missing := expectedCount - agg.Count
		if missing < 0 {
			missing = 0
		}
		fs.MissingRate = float64(missing) / float64(expectedCount)
	}

	// 变化百分比
	if agg.First != 0 {
		fs.ChangePct = (agg.Last - agg.First) / math.Abs(agg.First) * 100
	}

	// 趋势方向
	fs.Trend = calcTrend(agg.Last-agg.First, fs.Amplitude)

	// 简单斜率（每秒变化量）
	if durationSeconds > 0 {
		fs.Slope = (agg.Last - agg.First) / durationSeconds
	}

	return fs
}

// calcTrend 判断趋势方向
func calcTrend(change, amplitude float64) string {
	if amplitude == 0 {
		return "stable"
	}
	threshold := amplitude * 0.1
	if change > threshold {
		return "rising"
	}
	if change < -threshold {
		return "falling"
	}
	return "stable"
}

// CalcDataQuality 计算数据质量指标
func CalcDataQuality(lastTs string, missingRate float64) *DataQuality {
	dq := &DataQuality{
		MissingRate:  missingRate,
		LastDataTime: lastTs,
	}

	if lastTs != "" {
		t, err := time.Parse("2006-01-02 15:04:05.000", lastTs)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", lastTs)
		}
		if err == nil {
			dq.FreshnessSeconds = int64(time.Since(t).Seconds())
			// 超过 1 小时没有数据视为疑似离线
			dq.SuspectedOffline = dq.FreshnessSeconds > 3600
		}
	}

	return dq
}

// CalcRiskFromThreshold 根据阈值配置和字段统计计算风险信号
func CalcRiskFromThreshold(fs *FieldStats, tc *FieldThreshold) *RiskSignal {
	rs := &RiskSignal{Level: "normal"}
	if tc == nil {
		return rs
	}

	// 检查是否超限
	exceeded := false
	if tc.Upper != nil && fs.Max > *tc.Upper {
		exceeded = true
		rs.ExceedCount++
	}
	if tc.Lower != nil && fs.Min < *tc.Lower {
		exceeded = true
		rs.ExceedCount++
	}
	rs.Exceeded = exceeded

	// 风险等级
	if exceeded && tc.Upper != nil {
		upperVal := *tc.Upper
		divisor := math.Abs(upperVal)
		if divisor == 0 {
			divisor = 1
		}
		exceedRatio := math.Abs(fs.Max-upperVal) / divisor
		if exceedRatio > 0.5 {
			rs.Level = "critical"
		} else if exceedRatio > 0.2 {
			rs.Level = "high"
		} else {
			rs.Level = "warning"
		}
	} else if exceeded {
		rs.Level = "warning"
	}

	return rs
}

// BuildAgentHints 根据统计结果生成 agentHints
func BuildAgentHints(towerCode, deviceTypeCode string, fieldStats map[string]*FieldStats, radarStats map[int]map[string]*FieldStats, riskSignals map[string]*RiskSignal) []string {
	var hints []string

	label := towerCode
	if label == "" {
		label = deviceTypeCode
	}

	for field, fs := range fieldStats {
		if fs == nil {
			continue
		}
		if fs.MissingRate > 0.05 {
			hints = append(hints, fmt.Sprintf("%s %s 缺测率 %.1f%%，数据质量需关注", label, field, fs.MissingRate*100))
		}
		if rs, ok := riskSignals[field]; ok && rs.Exceeded {
			hints = append(hints, fmt.Sprintf("%s %s 存在超限，最大值 %.4g，风险等级 %s", label, field, fs.Max, rs.Level))
		}
		if fs.MissingRate <= 0.01 && !func() bool {
			rs, ok := riskSignals[field]
			return ok && rs.Exceeded
		}() {
			hints = append(hints, fmt.Sprintf("%s %s 均值 %.4g，缺测率 %.1f%%，未发现持续超限", label, field, fs.Avg, fs.MissingRate*100))
		}
	}

	// 雷达分层提示
	for idx, fields := range radarStats {
		for field, fs := range fields {
			if fs != nil && fs.Amplitude > 0.3 {
				hints = append(hints, fmt.Sprintf("%s 雷达 index=%d 的 %s 波动较大（幅度 %.4g）", label, idx, field, fs.Amplitude))
			}
		}
	}

	if len(hints) == 0 {
		hints = append(hints, fmt.Sprintf("%s %s 数据状态正常", label, deviceTypeCode))
	}

	return hints
}

// BuildSummary 生成 summary 文案
func BuildSummary(towerCode, deviceTypeCode string, fieldStats map[string]*FieldStats, radarStats map[int]map[string]*FieldStats) string {
	label := towerCode
	if label == "" {
		label = deviceTypeCode
	}

	if len(radarStats) > 0 {
		// 雷达设备
		indexCount := len(radarStats)
		maxTiField := ""
		var maxTiAmp float64
		for idx, fields := range radarStats {
			if ti, ok := fields["ti"]; ok && ti != nil {
				if ti.Amplitude > maxTiAmp {
					maxTiAmp = ti.Amplitude
					maxTiField = fmt.Sprintf("index=%d", idx)
				}
			}
		}
		if maxTiField != "" {
			return fmt.Sprintf("%s 测风雷达 %d 个距离层已分析，%s 的 ti 波动最大（幅度 %.4g），未配置持续超限阈值。", label, indexCount, maxTiField, maxTiAmp)
		}
		return fmt.Sprintf("%s 测风雷达 %d 个距离层已分析，数据状态正常。", label, indexCount)
	}

	if len(fieldStats) == 0 {
		return fmt.Sprintf("%s %s 暂无统计数据。", label, deviceTypeCode)
	}

	// 普通传感器，取第一个字段
	for field, fs := range fieldStats {
		if fs == nil {
			continue
		}
		return fmt.Sprintf("%s %s 在该时段均值 %.4g，最大 %.4g，缺测率 %.1f%%，未发现持续超限。",
			label, field, fs.Avg, fs.Max, fs.MissingRate*100)
	}

	return fmt.Sprintf("%s %s 趋势分析完成。", label, deviceTypeCode)
}

// ParseFloat 安全解析字符串为 float64
func ParseFloat(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	return v, err == nil
}
