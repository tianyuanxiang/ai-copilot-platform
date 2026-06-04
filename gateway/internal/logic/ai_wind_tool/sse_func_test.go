package ai_wind_tool

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteChatSSEKeepsConversationID(t *testing.T) {
	recorder := httptest.NewRecorder()
	if err := WriteChatSSE(recorder, StreamEvent{
		Type:           "done",
		TraceID:        "trace-1",
		ConversationID: "conversation-1",
	}); err != nil {
		t.Fatalf("WriteChatSSE() error = %v", err)
	}

	payload := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(recorder.Body.String()), "data: "))
	var event struct {
		ConversationID string `json:"conversationId"`
	}
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		t.Fatalf("decode SSE payload: %v", err)
	}
	if event.ConversationID != "conversation-1" {
		t.Fatalf("conversationId = %q, want %q", event.ConversationID, "conversation-1")
	}
}
