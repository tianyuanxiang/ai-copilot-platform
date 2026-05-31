// Package aiwindtool 的 audit 文件负责写入 Agent 工具调用审计日志。
// 每次调用只写一条最终日志，避免草稿工具和统一执行器重复记录。
package aiwindtool

import (
	"context"
	"encoding/json"
	"strings"
	"unicode/utf8"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const toolCallLogMaxBytes = 32768

type toolCallLogRow struct {
	ID        int64  `gorm:"column:id;primaryKey"`
	UserID    int64  `gorm:"column:user_id"`
	TraceID   string `gorm:"column:trace_id"`
	ToolName  string `gorm:"column:tool_name"`
	Arguments string `gorm:"column:arguments"`
	Result    string `gorm:"column:result"`
	Status    string `gorm:"column:status"`
	ErrorMsg  string `gorm:"column:error_msg"`
	LatencyMs int64  `gorm:"column:latency_ms"`
}

// writeToolCallLog 写入一条最终工具审计日志，并返回自增 ID。
// 现有 AiToolCallLogModel.Insert 不返回 PostgreSQL 自增 ID，而 ExecuteTool 响应需要 toolCallId，
// 因此这里使用 ServiceContext 已有的 GORM 连接执行 Create，不修改现有 model。
// 审计失败只记录错误，不中断 Agent 的正常推理链路。
func writeToolCallLog(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *pb.WindToolExecuteReq,
	status string,
	resultJSON string,
	message string,
	latencyMs int64,
) int64 {
	if svcCtx == nil || svcCtx.Orm == nil || req == nil {
		return 0
	}

	row := &toolCallLogRow{
		UserID:    req.UserId,
		TraceID:   req.TraceId,
		ToolName:  req.ToolName,
		Arguments: truncateLogJSON(req.ArgumentsJson),
		Result:    truncateLogJSON(resultJSON),
		Status:    status,
		LatencyMs: latencyMs,
	}
	if status != statusSuccess {
		row.ErrorMsg = truncateRunes(message, 1000)
	}
	if err := svcCtx.Orm.WithContext(ctx).Table("ai_tool_call_log").Create(row).Error; err != nil {
		logx.WithContext(ctx).Errorf("write ai_tool_call_log failed: %v", err)
		return 0
	}
	return row.ID
}

func truncateLogJSON(value string) string {
	if strings.TrimSpace(value) == "" {
		return "{}"
	}
	if len(value) <= toolCallLogMaxBytes {
		return value
	}
	data, _ := json.Marshal(map[string]any{
		"truncated":      true,
		"original_bytes": len(value),
	})
	return string(data)
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 || utf8.RuneCountInString(value) <= limit {
		return value
	}
	return string([]rune(value)[:limit])
}
