package aiwindtimeseriesservicelogic

import "ai-copilot-platform/ai-rpc/internal/utils/wind"

const (
	queryModeExact         = "exact"
	queryModeHourBucket    = "hour_bucket"
	queryModeRangeTooMin   = "range_too_small"
	queryModeRangeTooLarge = "range_too_large"

	granularityRaw = "raw"
	granularity1h  = "1h"

	sevenDays     = 7 * 24 * 3600
	thirtyOneDays = 31 * 24 * 3600
)

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
