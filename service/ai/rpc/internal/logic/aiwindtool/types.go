package aiwindtool

import "ai-copilot-platform/ai-rpc/pb"

const (
	ToolGetTurbineMetadata           = "get_turbine_metadata"
	ToolSearchMaintenanceSOP         = "search_maintenance_sop"
	ToolQueryAlarmEvents             = "query_alarm_events"
	ToolQuerySensorTimeseries        = "query_sensor_timeseries"
	ToolCompareSensorTrend           = "compare_sensor_trend"
	ToolGenerateAlarmAnalysisDraft   = "generate_alarm_analysis_draft"
	ToolGenerateHealthReport         = "generate_health_report"
	ToolCreateMaintenanceTicketDraft = "create_maintenance_ticket_draft"
)

type GetTurbineMetadataArgs struct {
	FarmCode       string `json:"farmCode,omitempty"`
	TowerCode      string `json:"towerCode,omitempty"`
	DeviceTypeCode string `json:"deviceTypeCode,omitempty"`
}

type SearchMaintenanceSOPArgs struct {
	Query       string  `json:"query"`
	KBID        int64   `json:"kbId,omitempty"`
	SearchScope string  `json:"searchScope,omitempty"`
	DomainID    int64   `json:"domainId,omitempty"`
	DocumentIDs []int64 `json:"documentIds,omitempty"`
	TopK        int64   `json:"topK,omitempty"`
	AnswerMode  string  `json:"answerMode,omitempty"`
}

type QueryAlarmEventsArgs struct {
	FarmCode  string `json:"farmCode"`
	TowerCode string `json:"towerCode,omitempty"`
	AlarmCode string `json:"alarmCode,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Status    int64  `json:"status,omitempty"`
	HasStatus bool   `json:"hasStatus,omitempty"`
}

type QuerySensorTimeseriesArgs struct {
	FarmCode       string   `json:"farmCode"`
	TowerCode      string   `json:"towerCode,omitempty"`
	DeviceCode     string   `json:"deviceCode,omitempty"`
	DeviceTypeCode string   `json:"deviceTypeCode"`
	Field          []string `json:"field,omitempty"`
	StartTime      string   `json:"startTime,omitempty"`
	EndTime        string   `json:"endTime,omitempty"`
	Page           int64    `json:"page,omitempty"`
	PageSize       int64    `json:"pageSize,omitempty"`
	IndexID        int64    `json:"indexId,omitempty"`
	RadarDistanceM int64    `json:"radarDistanceM,omitempty"`
}

type CompareSensorTrendArgs struct {
	FarmCode       string   `json:"farmCode"`
	TowerCode      string   `json:"towerCode,omitempty"`
	DeviceCode     string   `json:"deviceCode,omitempty"`
	DeviceTypeCode string   `json:"deviceTypeCode"`
	Field          []string `json:"field,omitempty"`
	StartTime      string   `json:"startTime,omitempty"`
	EndTime        string   `json:"endTime,omitempty"`
	IndexID        int64    `json:"indexId,omitempty"`
	RadarDistanceM int64    `json:"radarDistanceM,omitempty"`
}

type GenerateAlarmAnalysisDraftArgs struct {
	FarmCode     string `json:"farmCode"`
	TowerCode    string `json:"towerCode,omitempty"`
	AlarmCode    string `json:"alarmCode,omitempty"`
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
	EvidenceJSON string `json:"evidenceJson,omitempty"`
}

type GenerateHealthReportArgs struct {
	ReportType   string `json:"reportType,omitempty"`
	FarmCode     string `json:"farmCode"`
	TowerCode    string `json:"towerCode,omitempty"`
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
	EvidenceJSON string `json:"evidenceJson,omitempty"`
}

type CreateMaintenanceTicketDraftArgs struct {
	FarmCode     string `json:"farmCode"`
	TowerCode    string `json:"towerCode,omitempty"`
	AlarmCode    string `json:"alarmCode,omitempty"`
	Priority     string `json:"priority,omitempty"`
	EvidenceJSON string `json:"evidenceJson,omitempty"`
}

type toolResult struct {
	ResultJSON   string
	EvidenceJSON string
	Citations    []*pb.Citation
	Message      string
	Status       string
}
