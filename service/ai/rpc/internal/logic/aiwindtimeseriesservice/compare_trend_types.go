package aiwindtimeseriesservicelogic

import "ai-copilot-platform/ai-rpc/internal/model"

type trendRequest struct {
	FarmCode          string
	TowerCode         string
	DeviceTypeCode    string
	DeviceCode        string
	Fields            []string
	StartTime         string
	EndTime           string
	IndexID           int64
	RadarDistanceM    int64
	ResolvedIndexID   int64
	MatchedDistanceM  float64
	DisplayMeta       model.WindDeviceDisplayMeta
	RequestedFields   []string
	AnalysisFields    []string
	ExcludedFields    []string
	UsedDefaultFields bool
}

type trendResult struct {
	Summary      string
	EvidenceJSON string
	Message      string
}

type fieldStat struct {
	Count         int64   `json:"count"`
	ExpectedCount int64   `json:"expectedCount"`
	MissingRate   float64 `json:"missingRate"`
	Min           float64 `json:"min"`
	Max           float64 `json:"max"`
	Avg           float64 `json:"avg"`
	First         float64 `json:"first"`
	FirstTs       string  `json:"firstTs"`
	Last          float64 `json:"last"`
	LastTs        string  `json:"lastTs"`
	Change        float64 `json:"change"`
	ChangePct     float64 `json:"changePct"`
	Trend         string  `json:"trend"`
}

type trendSignal struct {
	Direction string  `json:"direction"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"changePct"`
}

type riskSignal struct {
	Level string `json:"level"`
}

type dataQuality struct {
	MissingRate      float64 `json:"missingRate"`
	LastDataTime     string  `json:"lastDataTime"`
	FreshnessSeconds int64   `json:"freshnessSeconds"`
	SuspectedOffline bool    `json:"suspectedOffline"`
}

type timeRangeInfo struct {
	StartTime       string  `json:"startTime"`
	EndTime         string  `json:"endTime"`
	DurationSeconds float64 `json:"durationSeconds"`
}

type thresholdConfigInfo struct {
	Status string `json:"status"`
}

type trendEvidenceBody struct {
	Source             string               `json:"source"`
	Scaffold           bool                 `json:"scaffold"`
	Partial            bool                 `json:"partial"`
	QueryMode          string               `json:"queryMode"`
	Granularity        string               `json:"granularity"`
	FarmCode           string               `json:"farmCode"`
	TowerCode          string               `json:"towerCode"`
	DeviceTypeCode     string               `json:"deviceTypeCode"`
	DeviceTypeName     string               `json:"deviceTypeName,omitempty"`
	DeviceCode         string               `json:"deviceCode"`
	Database           string               `json:"database"`
	Stable             string               `json:"stable"`
	Fields             []string             `json:"fields"`
	FieldLabels        map[string]string    `json:"fieldLabels,omitempty"`
	FieldUnits         map[string]string    `json:"fieldUnits,omitempty"`
	RequestedFields    []string             `json:"requestedFields,omitempty"`
	AnalysisFields     []string             `json:"analysisFields,omitempty"`
	ExcludedFields     []string             `json:"excludedAnalysisFields,omitempty"`
	UsedDefaultFields  bool                 `json:"usedDefaultAnalysisFields,omitempty"`
	TimeRange          timeRangeInfo        `json:"timeRange"`
	FieldStats         map[string]fieldStat `json:"fieldStats"`
	TrendSignal        trendSignal          `json:"trendSignal"`
	RiskSignal         riskSignal           `json:"riskSignal"`
	DataQuality        dataQuality          `json:"dataQuality"`
	ThresholdConfig    thresholdConfigInfo  `json:"thresholdConfig"`
	AgentHints         []string             `json:"agentHints"`
	IndexID            int64                `json:"indexId,omitempty"`
	RequestedDistanceM int64                `json:"requestedDistanceM,omitempty"`
	MatchedDistanceM   float64              `json:"matchedDistanceM,omitempty"`
}

type queryPlan struct {
	Mode        string
	Granularity string
	Database    string
	Stable      string
	Fields      []string
	Where       string
	TimeRange   timeRangeInfo
}

type trendCalcParams struct {
	Stats           map[string]model.BasicStat
	First           map[string]string
	Last            map[string]string
	DurationSeconds float64
	DeviceTypeCode  string
	IsMinuteMode    bool
}

type trendBucketParams struct {
	Buckets         []model.BucketRow
	DurationSeconds float64
	DeviceTypeCode  string
}
