package aiwinddraft

import (
	"encoding/json"
	"testing"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/model"
)

func TestSelectAlarmSamplesDoesNotUseUnrecoveredReason(t *testing.T) {
	rows := []map[string]string{
		{"ts": "2026-05-25 10:00:00", "tower_id": "03", "alarm_level": "4", "alarm_code": "1001", "status": "1"},
		{"ts": "2026-05-25 10:01:00", "tower_id": "04", "alarm_level": "2", "alarm_code": "1002", "status": "0"},
	}

	samples := selectAlarmSamples(rows, []string{"1001"}, []string{"03"}, "2026-05-25 10:00:00")
	if len(samples) == 0 {
		t.Fatal("expected samples")
	}
	for _, sample := range samples {
		reasons, ok := sample["sample_reason"].([]string)
		if !ok {
			t.Fatalf("unexpected sample_reason type: %#v", sample["sample_reason"])
		}
		for _, reason := range reasons {
			if reason == "unrecovered" {
				t.Fatalf("status must not produce unrecovered reason: %#v", sample)
			}
		}
	}
}

func TestBuildEvidenceSummaryKeepsAlarmAggregateContract(t *testing.T) {
	evidence := []map[string]any{
		{
			"source":            "tdengine.alarm",
			"total":             12,
			"sampled":           2,
			"truncated":         true,
			"granularity":       "1h",
			"by_level":          map[string]int64{"4": 2, "2": 10},
			"by_status":         map[string]int64{"0": 12},
			"by_alarm_code_top": []map[string]any{{"key": "1001", "count": int64(5)}},
			"by_tower_top":      []map[string]any{{"key": "03", "count": int64(7)}},
			"time_bucket_count": 1,
			"first_ts":          "2026-05-25 10:00:00",
			"last_ts":           "2026-05-25 11:00:00",
			"peak_bucket":       map[string]any{"bucket_start": "2026-05-25 10:00:00", "count": int64(6)},
			"samples":           []map[string]any{{"alarm_level": 4}},
			"unrecovered_count": 1,
			"is_Delete_count":   1,
			"isDelete_count":    1,
		},
	}

	summaryJSON := BuildEvidenceSummary(model.WindEvidenceJSON(evidence))
	var summary []map[string]any
	if err := json.Unmarshal([]byte(summaryJSON), &summary); err != nil {
		t.Fatalf("decode summary: %v", err)
	}
	if len(summary) != 1 {
		t.Fatalf("expected one summary item, got %d", len(summary))
	}

	item := summary[0]
	for _, key := range []string{
		"total",
		"sampled",
		"truncated",
		"granularity",
		"by_level",
		"by_status",
		"by_alarm_code_top",
		"by_tower_top",
		"time_bucket_count",
		"first_ts",
		"last_ts",
		"peak_bucket",
	} {
		if _, ok := item[key]; !ok {
			t.Fatalf("expected summary to keep %q: %#v", key, item)
		}
	}
	for _, key := range []string{"samples", "unrecovered_count", "is_Delete_count", "isDelete_count", "by_time_bucket"} {
		if _, ok := item[key]; ok {
			t.Fatalf("expected summary to drop %q: %#v", key, item)
		}
	}
}

func TestDraftContentDoesNotIncludeMetrics(t *testing.T) {
	resp := &engine.WindDraftResponse{
		Status:  "draft",
		Summary: "测试摘要",
		Metrics: map[string]any{
			"alarm_count": 10,
			"risk":        "high",
			"by_time_bucket": []map[string]any{
				{"bucket_start": "2026-05-25 10:00:00", "count": 6},
				{"bucket_start": "2026-05-25 11:00:00", "count": 4},
			},
		},
		Sections:        []map[string]any{{"name": "问题描述", "content": "测试"}},
		Recommendations: []string{"建议1"},
		Todo:            []string{"待办1"},
		Message:         "测试消息",
	}

	contentJSON := DraftContent(resp)
	var content map[string]any
	if err := json.Unmarshal([]byte(contentJSON), &content); err != nil {
		t.Fatalf("decode content: %v", err)
	}

	// 验证 metrics 字段不存在
	if _, ok := content["metrics"]; ok {
		t.Fatalf("expected content to not include metrics key, got: %#v", content)
	}

	// 验证其他字段保留
	for _, key := range []string{"status", "summary", "sections", "recommendations", "todo", "message"} {
		if _, ok := content[key]; !ok {
			t.Fatalf("expected content to keep %q: %#v", key, content)
		}
	}
}
