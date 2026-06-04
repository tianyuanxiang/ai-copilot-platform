// Package aiwindtool 提供 Wind ReAct Agent 的 Go 工具执行层。
// 本文件只定义工具名称、入参和内部统一结果，避免业务文件重复声明协议细节。
package aiwindtool

import "ai-copilot-platform/ai-rpc/pb"

const (
	// ToolGetTurbineMetadata 查询风场、风机和设备元数据。
	ToolGetTurbineMetadata = "get_turbine_metadata"
	// ToolSearchMaintenanceSOP 从知识库检索维护 SOP。
	ToolSearchMaintenanceSOP = "search_maintenance_sop"
	// ToolQueryAlarmEvents 查询告警聚合证据。
	ToolQueryAlarmEvents = "query_alarm_events"
	// ToolQuerySensorTimeseries 查询传感器时序数据。
	ToolQuerySensorTimeseries = "query_sensor_timeseries"
	// ToolCompareSensorTrend 对比传感器趋势。
	ToolCompareSensorTrend = "compare_sensor_trend"
	// ToolGenerateAlarmAnalysisDraft 生成告警分析草稿。
	ToolGenerateAlarmAnalysisDraft = "generate_alarm_analysis_draft"
	// ToolGenerateHealthReport 生成健康报告草稿。
	ToolGenerateHealthReport = "generate_health_report"
	// ToolCreateMaintenanceTicketDraft 生成维修工单草稿。
	ToolCreateMaintenanceTicketDraft = "create_maintenance_ticket_draft"
)

// GetTurbineMetadataArgs 是风机元数据查询参数。
// farmCode 为空时返回风场列表；有 farmCode 时返回风机；
// 进一步提供 towerCode 时，还会返回该风机下的设备。
type GetTurbineMetadataArgs struct {
	FarmCode       string `json:"farmCode,omitempty"`
	TowerCode      string `json:"towerCode,omitempty"`
	DeviceTypeCode string `json:"deviceTypeCode,omitempty"`
}

// SearchMaintenanceSOPArgs 是维护 SOP 检索参数。
type SearchMaintenanceSOPArgs struct {
	Query       string  `json:"query"`
	KBID        int64   `json:"kbId,omitempty"`
	SearchScope string  `json:"searchScope,omitempty"`
	DomainID    int64   `json:"domainId,omitempty"`
	DocumentIDs []int64 `json:"documentIds,omitempty"`
	TopK        int64   `json:"topK,omitempty"`
	AnswerMode  string  `json:"answerMode,omitempty"`
}

// QueryAlarmEventsArgs 是告警事件查询参数。
type QueryAlarmEventsArgs struct {
	FarmCode  string `json:"farmCode"`
	TowerCode string `json:"towerCode,omitempty"`
	AlarmCode string `json:"alarmCode,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Status    int64  `json:"status,omitempty"`
	HasStatus bool   `json:"hasStatus,omitempty"`
}

// QuerySensorTimeseriesArgs 是传感器时序查询参数。
// Field 使用单数命名，以保持与现有 WindTimeseriesQueryReq 接口一致。
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
	RadarDistanceM int64    `json:"radarDistanceM,omitempty"` // 雷达测量距离,例如30/50
}

// CompareSensorTrendArgs 是传感器趋势对比参数。
type CompareSensorTrendArgs struct {
	FarmCode       string   `json:"farmCode"`
	TowerCode      string   `json:"towerCode,omitempty"`
	DeviceCode     string   `json:"deviceCode,omitempty"`
	DeviceTypeCode string   `json:"deviceTypeCode"`
	Field          []string `json:"field,omitempty"`
	StartTime      string   `json:"startTime,omitempty"`
	EndTime        string   `json:"endTime,omitempty"`
	IndexID        int64    `json:"indexId,omitempty"`
	RadarDistanceM int64    `json:"radarDistanceM,omitempty"` // 雷达测量距离,例如30/50
}

// GenerateAlarmAnalysisDraftArgs 是告警分析草稿参数。
type GenerateAlarmAnalysisDraftArgs struct {
	FarmCode     string `json:"farmCode"`
	TowerCode    string `json:"towerCode,omitempty"`
	AlarmCode    string `json:"alarmCode,omitempty"`
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
	EvidenceJSON string `json:"evidenceJson,omitempty"`
}

// GenerateHealthReportArgs 是健康报告草稿参数。
type GenerateHealthReportArgs struct {
	ReportType   string `json:"reportType,omitempty"`
	FarmCode     string `json:"farmCode"`
	TowerCode    string `json:"towerCode,omitempty"`
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
	EvidenceJSON string `json:"evidenceJson,omitempty"`
}

// CreateMaintenanceTicketDraftArgs 是维修工单草稿参数。
type CreateMaintenanceTicketDraftArgs struct {
	FarmCode     string `json:"farmCode"`
	TowerCode    string `json:"towerCode,omitempty"`
	AlarmCode    string `json:"alarmCode,omitempty"`
	Priority     string `json:"priority,omitempty"`
	EvidenceJSON string `json:"evidenceJson,omitempty"`
}

// toolResult 是工具执行函数返回给统一执行器的内部结果。
// 执行器会将它转换为 RPC WindToolExecuteResp，并统一记录审计日志。
type toolResult struct {
	ResultJSON   string
	EvidenceJSON string
	Citations    []*pb.Citation
	Message      string
}
