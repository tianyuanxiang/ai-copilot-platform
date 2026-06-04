package aiwindtimeseriesservicelogic

import "ai-copilot-platform/ai-rpc/internal/model"

// trendRequest 趋势对比内部请求参数
type trendRequest struct {
	FarmCode       string
	TowerCode      string
	DeviceTypeCode string
	DeviceCode     string
	Fields         []string
	StartTime      string
	EndTime        string
	IndexID        int64
}

// trendResult 趋势对比结果
type trendResult struct {
	Summary      string
	EvidenceJSON string
	Message      string
}

// fieldStat 单个字段的统计结果（带 json tag，直接序列化进 evidence）
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

// trendSignal 趋势信号
type trendSignal struct {
	Direction string  `json:"direction"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"changePct"`
}

// riskSignal 风险信号
type riskSignal struct {
	Level string `json:"level"`
}

// dataQuality 数据质量
type dataQuality struct {
	MissingRate      float64 `json:"missingRate"`
	LastDataTime     string  `json:"lastDataTime"`
	FreshnessSeconds int64   `json:"freshnessSeconds"`
	SuspectedOffline bool    `json:"suspectedOffline"`
}

// timeRangeInfo evidence 中的时间范围
type timeRangeInfo struct {
	StartTime       string  `json:"startTime"`
	EndTime         string  `json:"endTime"`
	DurationSeconds float64 `json:"durationSeconds"`
}

// thresholdConfigInfo evidence 中的阈值配置
type thresholdConfigInfo struct {
	Status string `json:"status"`
}

// trendEvidenceBody evidence 完整结构
type trendEvidenceBody struct {
	Source          string               `json:"source"`
	Scaffold        bool                 `json:"scaffold"`
	Partial         bool                 `json:"partial"`
	QueryMode       string               `json:"queryMode"`
	Granularity     string               `json:"granularity"`
	FarmCode        string               `json:"farmCode"`
	TowerCode       string               `json:"towerCode"`
	DeviceTypeCode  string               `json:"deviceTypeCode"`
	DeviceCode      string               `json:"deviceCode"`
	Database        string               `json:"database"`
	Stable          string               `json:"stable"`
	Fields          []string             `json:"fields"`
	TimeRange       timeRangeInfo        `json:"timeRange"`
	FieldStats      map[string]fieldStat `json:"fieldStats"`
	TrendSignal     trendSignal          `json:"trendSignal"`
	RiskSignal      riskSignal           `json:"riskSignal"`
	DataQuality     dataQuality          `json:"dataQuality"`
	ThresholdConfig thresholdConfigInfo  `json:"thresholdConfig"`
	AgentHints      []string             `json:"agentHints"`
}

// queryPlan 查询策略
type queryPlan struct {
	Mode        string // exact / minute_bucket / range_too_large
	Granularity string // raw / 1h
	Database    string
	Stable      string
	Fields      []string
	Where       string
	TimeRange   timeRangeInfo
}

// trendCalcParams 计算 FieldStat 的输入参数
type trendCalcParams struct {
	Stats           map[string]model.BasicStat
	First           map[string]string
	Last            map[string]string
	DurationSeconds float64
	DeviceTypeCode  string
	IsMinuteMode    bool
}

// trendBucketParams 从 bucket 数据计算 FieldStat 的输入参数
type trendBucketParams struct {
	Buckets         []model.BucketRow
	DurationSeconds float64
	DeviceTypeCode  string
}
