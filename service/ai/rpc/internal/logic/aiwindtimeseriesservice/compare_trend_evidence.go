package aiwindtimeseriesservicelogic

import "ai-copilot-platform/ai-rpc/internal/model"

func buildEvidence(plan queryPlan, req trendRequest, fieldStats map[string]fieldStat, ts trendSignal,
	rs riskSignal, dq dataQuality, thresholdStatus string, agentHints []string, scaffold bool, partial bool) string {
	return model.WindEvidenceJSON(newTrendEvidenceBody(plan, req, fieldStats, ts, rs, dq, thresholdStatus, agentHints, scaffold, partial))
}

func buildScaffoldEvidence(req trendRequest, queryMode string, message string) string {
	body := newTrendEvidenceBody(queryPlan{Mode: queryMode, Fields: req.AnalysisFields}, req, map[string]fieldStat{},
		trendSignal{}, riskSignal{Level: riskNormal}, dataQuality{}, "not_configured", []string{}, true, false)
	if message != "" {
		body.AgentHints = append(body.AgentHints, message)
	}
	return model.WindEvidenceJSON(body)
}

func buildPartialEvidence(plan queryPlan, req trendRequest, fieldStats map[string]fieldStat, ts trendSignal, rs riskSignal, dq dataQuality, agentHints []string) string {
	return buildEvidence(plan, req, fieldStats, ts, rs, dq, "not_configured", agentHints, false, true)
}

func newTrendEvidenceBody(plan queryPlan, req trendRequest, fieldStats map[string]fieldStat, ts trendSignal,
	rs riskSignal, dq dataQuality, thresholdStatus string, agentHints []string, scaffold bool, partial bool) trendEvidenceBody {
	fields := plan.Fields
	if len(fields) == 0 {
		fields = req.AnalysisFields
	}
	return trendEvidenceBody{
		Source:             "tdengine",
		Scaffold:           scaffold,
		Partial:            partial,
		QueryMode:          plan.Mode,
		Granularity:        plan.Granularity,
		FarmCode:           req.FarmCode,
		TowerCode:          req.TowerCode,
		DeviceTypeCode:     req.DeviceTypeCode,
		DeviceTypeName:     req.DisplayMeta.DeviceTypeName,
		DeviceCode:         req.DeviceCode,
		Database:           plan.Database,
		Stable:             plan.Stable,
		Fields:             fields,
		FieldLabels:        req.DisplayMeta.FieldLabels,
		FieldUnits:         req.DisplayMeta.FieldUnits,
		RequestedFields:    req.RequestedFields,
		AnalysisFields:     req.AnalysisFields,
		ExcludedFields:     req.ExcludedFields,
		UsedDefaultFields:  req.UsedDefaultFields,
		TimeRange:          trendTimeRange(req),
		FieldStats:         fieldStats,
		TrendSignal:        ts,
		RiskSignal:         rs,
		DataQuality:        dq,
		ThresholdConfig:    thresholdConfigInfo{Status: thresholdStatus},
		AgentHints:         agentHints,
		IndexID:            req.ResolvedIndexID,
		RequestedDistanceM: req.RadarDistanceM,
		MatchedDistanceM:   req.MatchedDistanceM,
	}
}

func trendTimeRange(req trendRequest) timeRangeInfo {
	return timeRangeInfo{
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationSeconds: calcDurationSeconds(req.StartTime, req.EndTime),
	}
}
