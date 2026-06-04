package aiwindtimeseriesservicelogic

import "ai-copilot-platform/ai-rpc/internal/model"

// buildEvidence 根据查询计划、统计结果等构建完整 evidence 并序列化为 JSON。
func buildEvidence(plan queryPlan, req trendRequest, fieldStats map[string]fieldStat, ts trendSignal, rs riskSignal, dq dataQuality, thresholdStatus string, agentHints []string, scaffold bool, partial bool) string {
	body := trendEvidenceBody{
		Source:         "tdengine",
		Scaffold:       scaffold,
		Partial:        partial,
		QueryMode:      plan.Mode,
		Granularity:    plan.Granularity,
		FarmCode:       req.FarmCode,
		TowerCode:      req.TowerCode,
		DeviceTypeCode: req.DeviceTypeCode,
		Database:       plan.Database,
		Stable:         plan.Stable,
		Fields:         plan.Fields,
		TimeRange: timeRangeInfo{
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationSeconds: calcDurationSeconds(req.StartTime, req.EndTime),
		},
		FieldStats:  fieldStats,
		TrendSignal: ts,
		RiskSignal:  rs,
		DataQuality: dq,
		ThresholdConfig: thresholdConfigInfo{
			Status: thresholdStatus,
		},
		AgentHints: agentHints,
	}

	return model.WindEvidenceJSON(body)
}

// buildScaffoldEvidence 构建 scaffold evidence（TDengine 未配置或范围过大时使用）。
func buildScaffoldEvidence(req trendRequest, queryMode string, message string) string {
	body := trendEvidenceBody{
		Source:         "tdengine",
		Scaffold:       true,
		Partial:        false,
		QueryMode:      queryMode,
		FarmCode:       req.FarmCode,
		TowerCode:      req.TowerCode,
		DeviceTypeCode: req.DeviceTypeCode,
		DeviceCode:     req.DeviceCode,
		Fields:         req.Fields,
		TimeRange: timeRangeInfo{
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationSeconds: calcDurationSeconds(req.StartTime, req.EndTime),
		},
		FieldStats:      make(map[string]fieldStat),
		AgentHints:      []string{},
		ThresholdConfig: thresholdConfigInfo{Status: "not_configured"},
	}

	return model.WindEvidenceJSON(body)
}

// buildPartialEvidence 构建部分 evidence（超时等场景）。
func buildPartialEvidence(plan queryPlan, req trendRequest, fieldStats map[string]fieldStat, ts trendSignal, rs riskSignal, dq dataQuality, agentHints []string) string {
	body := trendEvidenceBody{
		Source:         "tdengine",
		Scaffold:       false,
		Partial:        true,
		QueryMode:      plan.Mode,
		Granularity:    plan.Granularity,
		FarmCode:       req.FarmCode,
		TowerCode:      req.TowerCode,
		DeviceTypeCode: req.DeviceTypeCode,
		Database:       plan.Database,
		Stable:         plan.Stable,
		Fields:         plan.Fields,
		TimeRange: timeRangeInfo{
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationSeconds: calcDurationSeconds(req.StartTime, req.EndTime),
		},
		FieldStats:      fieldStats,
		TrendSignal:     ts,
		RiskSignal:      rs,
		DataQuality:     dq,
		ThresholdConfig: thresholdConfigInfo{Status: "not_configured"},
		AgentHints:      agentHints,
	}

	return model.WindEvidenceJSON(body)
}
