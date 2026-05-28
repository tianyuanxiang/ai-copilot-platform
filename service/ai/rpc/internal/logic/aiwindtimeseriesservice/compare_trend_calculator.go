package aiwindtimeseriesservicelogic

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/utils/wind"
)

const (
	trendStable = "stable"
	trendRising = "rising"
	trendFalling = "falling"

	riskNormal = "normal"

	trendChangePctThreshold = 5.0
)

// calcFieldStat 根据精确模式基础统计 + 边界行计算单个字段的 fieldStat。
func calcFieldStat(stat model.BasicStat, expectedCount int64, firstVal, firstTs string, lastVal, lastTs string) fieldStat {
	fs := fieldStat{
		Count:         stat.Count,
		ExpectedCount: expectedCount,
		Min:           stat.Min,
		Max:           stat.Max,
		Avg:           stat.Avg,
		FirstTs:       firstTs,
		LastTs:        lastTs,
	}

	first := parseFloatZero(firstVal)
	last := parseFloatZero(lastVal)
	fs.First = first
	fs.Last = last
	fs.Change = last - first

	if first != 0 {
		fs.ChangePct = (last - first) / math.Abs(first) * 100
	}

	fs.Trend = calcTrendDir(fs.ChangePct, first, fs.Change, fs.Max-fs.Min)

	if expectedCount > 0 {
		missing := expectedCount - stat.Count
		if missing < 0 {
			missing = 0
		}
		fs.MissingRate = float64(missing) / float64(expectedCount)
	}

	return fs
}

// calcFieldStatFromBuckets 根据分钟级 bucket 数据计算 fieldStat。
// 从 BucketRow 切片中提取指定 field 的数据，内部计算首尾 bucket。
func calcFieldStatFromBuckets(buckets []model.BucketRow, field string, durationMinutes float64) fieldStat {
	var counts []int64
	var avgs []float64
	var mins []float64
	var maxs []float64
	firstBucketAvg := ""
	firstBucketTs := ""
	lastBucketAvg := ""
	lastBucketTs := ""

	for _, b := range buckets {
		if b.Field != field {
			continue
		}
		counts = append(counts, b.Count)
		avgs = append(avgs, b.Avg)
		mins = append(mins, b.Min)
		maxs = append(maxs, b.Max)
		if firstBucketTs == "" {
			firstBucketTs = b.Wstart
			firstBucketAvg = fmt.Sprintf("%g", b.Avg)
		}
		lastBucketTs = b.Wstart
		lastBucketAvg = fmt.Sprintf("%g", b.Avg)
	}

	if len(counts) == 0 || len(avgs) == 0 {
		return fieldStat{}
	}

	var totalCount int64
	for _, c := range counts {
		totalCount += c
	}
	var overallMin = maxs[0]
	var overallMax = mins[0]
	var sumAvg float64
	for i, a := range avgs {
		sumAvg += a
		if mins[i] < overallMin {
			overallMin = mins[i]
		}
		if maxs[i] > overallMax {
			overallMax = maxs[i]
		}
	}
	overallAvg := sumAvg / float64(len(avgs))

	first := parseFloatZero(firstBucketAvg)
	last := parseFloatZero(lastBucketAvg)
	fs := fieldStat{
		Count:   totalCount,
		Min:     overallMin,
		Max:     overallMax,
		Avg:     overallAvg,
		First:   first,
		Last:    last,
		FirstTs: firstBucketTs,
		LastTs:  lastBucketTs,
		Change:  last - first,
	}

	if first != 0 {
		fs.ChangePct = (last - first) / math.Abs(first) * 100
	}
	fs.Trend = calcTrendDir(fs.ChangePct, first, fs.Change, overallMax-overallMin)

	expectedCount := int64(durationMinutes)
	if expectedCount > 0 {
		actualBucketCount := int64(len(counts))
		missing := expectedCount - actualBucketCount
		if missing < 0 {
			missing = 0
		}
		fs.MissingRate = float64(missing) / float64(expectedCount)
		fs.ExpectedCount = expectedCount
	}

	return fs
}

// calcTrendDir 计算趋势方向。
func calcTrendDir(changePct, first, change, amplitude float64) string {
	if first == 0 {
		threshold := amplitude * 0.1
		if change > threshold {
			return trendRising
		}
		if change < -threshold {
			return trendFalling
		}
		return trendStable
	}

	if changePct >= trendChangePctThreshold {
		return trendRising
	}
	if changePct <= -trendChangePctThreshold {
		return trendFalling
	}
	return trendStable
}

// calcDataQuality 计算数据质量指标。
func calcDataQuality(lastTs string, missingRate float64) dataQuality {
	dq := dataQuality{
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
			dq.SuspectedOffline = dq.FreshnessSeconds > 3600
		}
	}

	return dq
}

// calcRiskSignal 根据阈值配置和字段统计计算风险信号。
func calcRiskSignal(fs fieldStat, threshold *model.FieldThreshold) riskSignal {
	rs := riskSignal{Level: riskNormal}

	if threshold == nil {
		return rs
	}

	exceeded := false
	if threshold.Upper != nil && fs.Max > *threshold.Upper {
		exceeded = true
	}
	if threshold.Lower != nil && fs.Min < *threshold.Lower {
		exceeded = true
	}

	if !exceeded {
		return rs
	}

	if threshold.Upper != nil {
		upperVal := *threshold.Upper
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
	} else {
		rs.Level = "warning"
	}

	return rs
}

// calcExpectedCount 根据采样策略计算期望数据量。
func calcExpectedCount(deviceTypeCode string, durationSeconds float64, isMinuteMode bool) int64 {
	if isMinuteMode {
		return int64(durationSeconds / 60.0)
	}
	policy := wind.ResolveSamplingPolicy(deviceTypeCode)
	return policy.ExpectedCount(durationSeconds, 1)
}

// buildSummary 生成趋势分析 summary 文案。
func buildSummary(towerCode, deviceTypeCode string, fieldStats map[string]fieldStat) string {
	label := towerCode
	if label == "" {
		label = deviceTypeCode
	}

	if len(fieldStats) == 0 {
		return fmt.Sprintf("%s %s 暂无统计数据。", label, deviceTypeCode)
	}

	for field, fs := range fieldStats {
		return fmt.Sprintf("%s %s 在该时段均值 %.4g，最大 %.4g，趋势 %s，缺测率 %.1f%%",
			label, field, fs.Avg, fs.Max, trendDesc(fs.Trend), fs.MissingRate*100)
	}

	return fmt.Sprintf("%s %s 趋势分析完成。", label, deviceTypeCode)
}

// buildAgentHints 根据统计结果生成 agent 提示。
func buildAgentHints(towerCode, deviceTypeCode string, fieldStats map[string]fieldStat, riskSignals map[string]riskSignal) []string {
	var hints []string

	label := towerCode
	if label == "" {
		label = deviceTypeCode
	}

	for field, fs := range fieldStats {
		if fs.MissingRate > 0.05 {
			hints = append(hints, fmt.Sprintf("%s %s 缺测率 %.1f%%，数据质量需关注", label, field, fs.MissingRate*100))
		}
		if rs, ok := riskSignals[field]; ok && rs.Level != riskNormal {
			hints = append(hints, fmt.Sprintf("%s %s 存在超限，最大值 %.4g，风险等级 %s", label, field, fs.Max, rs.Level))
		}
		if fs.MissingRate <= 0.01 {
			if rs, ok := riskSignals[field]; !ok || rs.Level == riskNormal {
				hints = append(hints, fmt.Sprintf("%s %s 均值 %.4g，缺测率 %.1f%%，未发现持续超限", label, field, fs.Avg, fs.MissingRate*100))
			}
		}
	}

	if len(hints) == 0 {
		hints = append(hints, fmt.Sprintf("%s %s 数据状态正常", label, deviceTypeCode))
	}

	return hints
}

func trendDesc(trend string) string {
	switch trend {
	case trendRising:
		return "上升"
	case trendFalling:
		return "下降"
	default:
		return "平稳"
	}
}

func parseFloatZero(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func riskLevelOrder(level string) int {
	switch level {
	case "critical":
		return 4
	case "high":
		return 3
	case "warning":
		return 2
	case "normal":
		return 1
	default:
		return 0
	}
}
