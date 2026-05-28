package engine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWindDraftClientCallsEndpoints(t *testing.T) {
	paths := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths[r.URL.Path] = true
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var payload WindDraftRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload.TraceID != "wind-test" || payload.FarmCode != "FY" {
			t.Fatalf("unexpected payload: %+v", payload)
		}
		_ = json.NewEncoder(w).Encode(WindDraftResponse{
			Title:         "draft",
			Status:        "draft",
			Summary:       "ok",
			EvidenceCount: 1,
			Message:       "done",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	payload := WindDraftRequest{
		UserID:   "1001",
		TraceID:  "wind-test",
		FarmCode: "FY",
		Evidence: []map[string]any{{"source": "test"}},
	}

	if _, err := client.WindAlarmSummary(context.Background(), payload); err != nil {
		t.Fatalf("WindAlarmSummary: %v", err)
	}
	if _, err := client.WindHealthReportDraft(context.Background(), payload); err != nil {
		t.Fatalf("WindHealthReportDraft: %v", err)
	}
	if _, err := client.WindTicketDraft(context.Background(), payload); err != nil {
		t.Fatalf("WindTicketDraft: %v", err)
	}

	for _, path := range []string{"/v1/wind/summary/alarm", "/v1/wind/reports/health/draft", "/v1/wind/tickets/draft"} {
		if !paths[path] {
			t.Fatalf("endpoint was not called: %s", path)
		}
	}
}
