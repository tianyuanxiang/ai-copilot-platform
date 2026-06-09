package aiwindtimeseriesservicelogic

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/utils/wind"
)

const (
	trendStable  = "stable"
	trendRising  = "rising"
	trendFalling = "falling"

	riskNormal = "normal"

	trendChangePctThreshold = 5.0
	wprDeviceTypeCode       = "WPR"
	wprDistanceField        = "d"
)

func calcFieldStat(stat model.BasicStat, expectedCount int64, firstVal, firstTs string, lastVal, lastTs string) fieldStat {
	fs := fieldStat{
		Count:         stat.Count.Int64,
		ExpectedCount: expectedCount,
		Min:           stat.Min.Float64,
		Max:           stat.Max.Float64,
		Avg:           stat.Avg.Float64,
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
		missing := expectedCount - stat.Count.Int64
		if missing < 0 {
			missing = 0
		}
		fs.MissingRate = float64(missing) / float64(expectedCount)
	}
	return fs
}

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
		counts = append(counts, b.Count.Int64)
		avgs = append(avgs, b.Avg.Float64)
		mins = append(mins, b.Min.Float64)
		maxs = append(maxs, b.Max.Float64)
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

	totalCount := int64(0)
	for _, c := range counts {
		totalCount += c
	}
	overallMin := mins[0]
	overallMax := maxs[0]
	sumAvg := 0.0
	for i, a := range avgs {
		sumAvg += a
		if mins[i] < overallMin {
			overallMin = mins[i]
		}
		if maxs[i] > overallMax {
			overallMax = maxs[i]
		}
	}

	first := parseFloatZero(firstBucketAvg)
	last := parseFloatZero(lastBucketAvg)
	fs := fieldStat{
		Count:   totalCount,
		Min:     overallMin,
		Max:     overallMax,
		Avg:     sumAvg / float64(len(avgs)),
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

	if threshold.Upper == nil {
		rs.Level = "warning"
		return rs
	}
	upperVal := *threshold.Upper
	divisor := math.Abs(upperVal)
	if divisor == 0 {
		divisor = 1
	}
	exceedRatio := math.Abs(fs.Max-upperVal) / divisor
	switch {
	case exceedRatio > 0.5:
		rs.Level = "critical"
	case exceedRatio > 0.2:
		rs.Level = "high"
	default:
		rs.Level = "warning"
	}
	return rs
}

func calcExpectedCount(deviceTypeCode string, durationSeconds float64, isMinuteMode bool) int64 {
	if isMinuteMode {
		return int64(durationSeconds / 60.0)
	}
	policy := wind.ResolveSamplingPolicy(deviceTypeCode)
	return policy.ExpectedCount(durationSeconds, 1)
}

func buildSummary(towerCode string, displayMeta model.WindDeviceDisplayMeta, fields []string, fieldStats map[string]fieldStat) string {
	label := towerCode
	if label == "" {
		label = displayMeta.DeviceTypeName
	}
	deviceName := displayMeta.DeviceTypeName
	if deviceName == "" {
		deviceName = displayMeta.DeviceTypeCode
	}
	if len(fieldStats) == 0 {
		return fmt.Sprintf("%s号风机 %s 暂无统计数据。", label, deviceName)
	}

	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		fs, ok := fieldStats[field]
		if !ok {
			continue
		}
		fieldName := displayMeta.FieldLabel(field)
		unit := displayMeta.FieldUnit(field)
		avgText := fmt.Sprintf("%.4g", fs.Avg)
		maxText := fmt.Sprintf("%.4g", fs.Max)
		if unit != "" {
			avgText += unit
			maxText += unit
		}
		parts = append(parts, fmt.Sprintf("%s均值%s、最大%s、趋势%s、缺测率 %.1f%%",
			fieldName, avgText, maxText, trendDesc(fs.Trend), fs.MissingRate*100))
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s号风机 %s 暂无统计数据。", label, deviceName)
	}
	return fmt.Sprintf("%s号风机 %s 趋势分析：%s。", label, deviceName, strings.Join(parts, "；"))
}

func buildAgentHints(towerCode string, displayMeta model.WindDeviceDisplayMeta, fields []string, fieldStats map[string]fieldStat, riskSignals map[string]riskSignal) []string {
	hints := make([]string, 0, len(fields))
	label := towerCode
	if label == "" {
		label = displayMeta.DeviceTypeName
	}
	for _, field := range fields {
		fs, ok := fieldStats[field]
		if !ok {
			continue
		}
		fieldName := displayMeta.FieldLabel(field)
		if fs.MissingRate > 0.05 {
			hints = append(hints, fmt.Sprintf("%s号风机 %s 缺测率 %.1f%%，数据质量需要关注", label, fieldName, fs.MissingRate*100))
		}
		if rs, ok := riskSignals[field]; ok && rs.Level != riskNormal {
			hints = append(hints, fmt.Sprintf("%s号风机 %s 存在超限，最大值 %.4g，风险等级 %s", label, fieldName, fs.Max, rs.Level))
		}
		if fs.MissingRate <= 0.01 {
			if rs, ok := riskSignals[field]; !ok || rs.Level == riskNormal {
				hints = append(hints, fmt.Sprintf("%s号风机 %s 均值 %.4g，缺测率 %.1f%%，未发现持续超限", label, fieldName, fs.Avg, fs.MissingRate*100))
			}
		}
	}
	if len(hints) == 0 {
		hints = append(hints, fmt.Sprintf("%s号风机 %s 数据状态正常", label, displayMeta.DeviceTypeName))
	}
	return hints
}

func splitRequestedTrendFields(fields []string) []string {
	result := make([]string, 0, len(fields))
	for _, field := range fields {
		for _, item := range strings.Split(field, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				result = append(result, item)
			}
		}
	}
	return result
}

func resolveTrendAnalysisFields(deviceTypeCode string, displayMeta model.WindDeviceDisplayMeta, requestedFields, selectedFields []string) ([]string, []string, bool) {
	if deviceTypeCode != wprDeviceTypeCode {
		return selectedFields, nil, false
	}

	analysisFields, excludedFields := excludeTrendFields(selectedFields, wprDistanceField)
	usedDefaultFields := len(requestedFields) == 0
	if len(analysisFields) == 0 {
		analysisFields, _ = excludeTrendFields(displayMeta.Fields, wprDistanceField)
		usedDefaultFields = true
	}
	return analysisFields, excludedFields, usedDefaultFields
}

func excludeTrendFields(fields []string, excluded ...string) ([]string, []string) {
	excludedSet := make(map[string]struct{}, len(excluded))
	for _, field := range excluded {
		excludedSet[field] = struct{}{}
	}
	kept := make([]string, 0, len(fields))
	removed := make([]string, 0)
	for _, field := range fields {
		if _, ok := excludedSet[field]; ok {
			removed = append(removed, field)
			continue
		}
		kept = append(kept, field)
	}
	return kept, removed
}

func firstFieldStat(fields []string, fieldStats map[string]fieldStat) (fieldStat, bool) {
	for _, field := range fields {
		if fs, ok := fieldStats[field]; ok {
			return fs, true
		}
	}
	return fieldStat{}, false
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
