package aiwinddraft

import (
	"encoding/json"
	"testing"

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
			"by_time_bucket":    []map[string]any{{"bucket_start": "2026-05-25 10:00:00", "count": int64(6)}},
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
		"by_time_bucket",
		"first_ts",
		"last_ts",
		"peak_bucket",
	} {
		if _, ok := item[key]; !ok {
			t.Fatalf("expected summary to keep %q: %#v", key, item)
		}
	}
	for _, key := range []string{"samples", "unrecovered_count", "is_Delete_count", "isDelete_count"} {
		if _, ok := item[key]; ok {
			t.Fatalf("expected summary to drop %q: %#v", key, item)
		}
	}
}
