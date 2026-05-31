package aiwindtool

import (
	"context"
	"testing"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
)

func TestExecutorRejectsNilRequest(t *testing.T) {
	resp, err := NewExecutor(nil).Execute(context.Background(), nil)
	if err != nil {
		t.Fatalf("Execute() RPC error = %v", err)
	}
	if resp.Status != statusFailed {
		t.Fatalf("status = %q, want %q", resp.Status, statusFailed)
	}
}

func TestExecutorRejectsInvalidUserID(t *testing.T) {
	resp, err := NewExecutor(&svc.ServiceContext{}).Execute(context.Background(), &pb.WindToolExecuteReq{
		UserId:   0,
		TraceId:  "trace-1",
		ToolName: ToolGetTurbineMetadata,
	})
	if err != nil {
		t.Fatalf("Execute() RPC error = %v", err)
	}
	if resp.Status != statusFailed {
		t.Fatalf("status = %q, want %q", resp.Status, statusFailed)
	}
}

func TestExecutorRejectsEmptyTraceID(t *testing.T) {
	resp, err := NewExecutor(&svc.ServiceContext{}).Execute(context.Background(), &pb.WindToolExecuteReq{
		UserId:   1,
		ToolName: ToolGetTurbineMetadata,
	})
	if err != nil {
		t.Fatalf("Execute() RPC error = %v", err)
	}
	if resp.Status != statusFailed {
		t.Fatalf("status = %q, want %q", resp.Status, statusFailed)
	}
}

func TestExecutorDeniesUnknownToolWithoutAuditDatabase(t *testing.T) {
	resp, err := NewExecutor(&svc.ServiceContext{}).Execute(context.Background(), &pb.WindToolExecuteReq{
		UserId:   1,
		TraceId:  "trace-1",
		ToolName: "delete_everything",
	})
	if err != nil {
		t.Fatalf("Execute() RPC error = %v", err)
	}
	if resp.Status != statusDenied {
		t.Fatalf("status = %q, want %q", resp.Status, statusDenied)
	}
	if resp.ToolCallId != 0 {
		t.Fatalf("toolCallId = %d, want 0 when audit database is unavailable", resp.ToolCallId)
	}
}

func TestAllowedTools(t *testing.T) {
	tools := []string{
		ToolGetTurbineMetadata,
		ToolSearchMaintenanceSOP,
		ToolQueryAlarmEvents,
		ToolQuerySensorTimeseries,
		ToolCompareSensorTrend,
		ToolGenerateAlarmAnalysisDraft,
		ToolGenerateHealthReport,
		ToolCreateMaintenanceTicketDraft,
	}
	if len(allowedTools) != len(tools) {
		t.Fatalf("allowedTools size = %d, want %d", len(allowedTools), len(tools))
	}
	for _, tool := range tools {
		if !IsAllowedTool(tool) {
			t.Fatalf("tool %q should be allowed", tool)
		}
	}
}
