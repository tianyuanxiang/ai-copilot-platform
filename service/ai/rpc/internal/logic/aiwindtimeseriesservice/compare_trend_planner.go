package aiwindtimeseriesservicelogic

import (
	"ai-copilot-platform/ai-rpc/internal/utils/wind"
)

const (
	queryModeExact         = "全量查询"
	queryModeHourBucket    = "按小时分桶查询"
	queryModeRangeTooMin   = "查询范围太小"
	queryModeRangeTooLarge = "查询范围太大"

	granularityRaw = "raw"
	granularity1h  = "1h"

	sevenDays     = 7 * 24 * 3600
	thirtyOneDays = 31 * 24 * 3600
)

// planQuery 根据时间范围生成查询策略。
func planQuery(database, stable string, fields []string, where string, durationSeconds float64) queryPlan {
	if durationSeconds <= 0 {
		return queryPlan{Mode: queryModeRangeTooMin}
	}

	if durationSeconds <= sevenDays {
		return queryPlan{
			Mode:        queryModeExact,
			Granularity: granularityRaw,
			Database:    database,
			Stable:      stable,
			Fields:      fields,
			Where:       where,
		}
	}

	if durationSeconds <= thirtyOneDays {
		return queryPlan{
			Mode:        queryModeHourBucket,
			Granularity: granularity1h,
			Database:    database,
			Stable:      stable,
			Fields:      fields,
			Where:       where,
		}
	}

	return queryPlan{Mode: queryModeRangeTooLarge}
}

// calcDurationSeconds 计算两个时间戳之间的秒数。
func calcDurationSeconds(startTime, endTime string) float64 {
	return wind.CalcDurationSeconds(startTime, endTime)
}

func (p queryPlan) isExactMode() bool {
	return p.Mode == queryModeExact
}

func (p queryPlan) isHourBucketMode() bool {
	return p.Mode == queryModeHourBucket
}

func (p queryPlan) isRangeTooLarge() bool {
	return p.Mode == queryModeRangeTooLarge
}
