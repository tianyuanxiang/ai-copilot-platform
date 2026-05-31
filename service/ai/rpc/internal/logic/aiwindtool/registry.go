// Package aiwindtool 的 registry 文件维护 Agent 可调用工具白名单。
// 白名单是工具执行层的边界，不在这里引入动态注册或复杂权限框架。
package aiwindtool

import "strings"

var allowedTools = map[string]struct{}{
	ToolGetTurbineMetadata:           {},
	ToolSearchMaintenanceSOP:         {},
	ToolQueryAlarmEvents:             {},
	ToolQuerySensorTimeseries:        {},
	ToolCompareSensorTrend:           {},
	ToolGenerateAlarmAnalysisDraft:   {},
	ToolGenerateHealthReport:         {},
	ToolCreateMaintenanceTicketDraft: {},
}

// IsAllowedTool 判断工具名是否位于明确允许的白名单中。
func IsAllowedTool(name string) bool {
	_, ok := allowedTools[strings.TrimSpace(name)]
	return ok
}
