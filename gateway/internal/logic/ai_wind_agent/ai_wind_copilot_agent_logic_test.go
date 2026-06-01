package ai_wind_agent

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"ai-copilot-platform/ai-rpc/pb"
)

type testWindAgentStreamReceiver struct {
	events []*pb.WindAgentStreamEvent
}

func (r *testWindAgentStreamReceiver) Recv() (*pb.WindAgentStreamEvent, error) {
	if len(r.events) == 0 {
		return nil, io.EOF
	}
	event := r.events[0]
	r.events = r.events[1:]
	return event, nil
}

func TestForwardWindAgentStreamKeepsDoneContent(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := forwardWindAgentStream(recorder, &testWindAgentStreamReceiver{
		events: []*pb.WindAgentStreamEvent{{
			Type:           "done",
			TraceId:        "trace-1",
			AgentSessionId: "session-1",
			Content:        "完整处理建议",
		}},
	}, "")
	if err != nil {
		t.Fatalf("forwardWindAgentStream() error = %v", err)
	}

	payload := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(recorder.Body.String()), "data: "))
	var event typesWindAgentStreamEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		t.Fatalf("decode SSE payload: %v", err)
	}
	if event.Content != "完整处理建议" {
		t.Fatalf("content = %q, want %q", event.Content, "完整处理建议")
	}
}

type typesWindAgentStreamEvent struct {
	Content string `json:"content"`
}
