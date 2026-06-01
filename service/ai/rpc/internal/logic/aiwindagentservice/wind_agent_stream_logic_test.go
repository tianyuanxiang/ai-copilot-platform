package aiwindagentservicelogic

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"

	"google.golang.org/grpc/metadata"
)

type testWindAgentStream struct {
	ctx     context.Context
	events  []*pb.WindAgentStreamEvent
	sendErr error
}

func (s *testWindAgentStream) Send(event *pb.WindAgentStreamEvent) error {
	s.events = append(s.events, event)
	return s.sendErr
}

func (s *testWindAgentStream) SetHeader(metadata.MD) error  { return nil }
func (s *testWindAgentStream) SendHeader(metadata.MD) error { return nil }
func (s *testWindAgentStream) SetTrailer(metadata.MD)       {}
func (s *testWindAgentStream) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return context.Background()
}
func (s *testWindAgentStream) SendMsg(any) error { return nil }
func (s *testWindAgentStream) RecvMsg(any) error { return nil }

func TestWindAgentStreamEventToPB(t *testing.T) {
	got := windAgentStreamEventToPB(engine.WindAgentStreamEvent{
		Type:           "confirmation_required",
		TraceID:        "trace-1",
		AgentSessionId: "conv-1",
		Content:        "approve?",
		ToolCall: &engine.WindAgentToolCall{
			ToolCallID:    1,
			ToolName:      "query_alarm",
			Status:        "success",
			ArgumentsJSON: `{"farm":"A"}`,
			ResultJSON:    `{"count":1}`,
			Message:       "ok",
			LatencyMS:     5,
		},
		ToolCalls: []engine.WindAgentToolCall{{ToolCallID: 2, ToolName: "search_sop", LatencyMS: 6}},
		Citations: []engine.WindAgentCitation{{DocumentID: 3, ChunkID: 4, Title: "SOP", Snippet: "step", Score: 0.9}},
		Draft:     &engine.WindAgentDraftRef{DraftType: "ticket", DraftID: 7, Title: "draft"},
		ErrorMsg:  "error detail",
	})

	if got.Type != "confirmation_required" || got.TraceId != "trace-1" || got.AgentSessionId != "conv-1" || got.Content != "approve?" {
		t.Fatalf("event = %+v", got)
	}
	if got.ToolCall == nil || got.ToolCall.LatencyMs != 5 || len(got.ToolCalls) != 1 || got.ToolCalls[0].LatencyMs != 6 {
		t.Fatalf("tool calls = %+v / %+v", got.ToolCall, got.ToolCalls)
	}
	if len(got.Citations) != 1 || got.Citations[0].Score != 0.9 {
		t.Fatalf("citations = %+v", got.Citations)
	}
	if got.Draft == nil || got.Draft.DraftId != 7 || got.ErrorMsg != "error detail" {
		t.Fatalf("draft/error = %+v / %q", got.Draft, got.ErrorMsg)
	}
}

func TestRunAgentStreamValidation(t *testing.T) {
	logic := NewRunAgentStreamLogic(context.Background(), &svc.ServiceContext{})
	stream := &testWindAgentStream{}
	for _, request := range []*pb.WindAgentRunReq{
		nil,
		{UserId: 0, Input: "question"},
		{UserId: 1, Input: "  "},
	} {
		if err := logic.RunAgentStream(request, stream); err == nil {
			t.Fatalf("RunAgentStream(%+v) error = nil", request)
		}
	}
}

func TestResumeAgentStreamValidation(t *testing.T) {
	logic := NewResumeAgentStreamLogic(context.Background(), &svc.ServiceContext{})
	stream := &testWindAgentStream{}
	for _, request := range []*pb.WindAgentResumeReq{
		nil,
		{UserId: 0, AgentSessionId: "conv-1", Action: "approve"},
		{UserId: 1, AgentSessionId: "  ", Action: "approve"},
		{UserId: 1, AgentSessionId: "conv-1", Action: "unknown"},
		{UserId: 1, AgentSessionId: "conv-1", Action: "clarify", Content: "  "},
	} {
		if err := logic.ResumeAgentStream(request, stream); err == nil {
			t.Fatalf("ResumeAgentStream(%+v) error = nil", request)
		}
	}
}

func TestRunAgentStreamForwardsPythonErrorEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `data: {"type":"error","trace_id":"trace-1","agent_session_id":"conv-1","error_msg":"agent failed"}`+"\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	stream := &testWindAgentStream{}
	logic := NewRunAgentStreamLogic(context.Background(), &svc.ServiceContext{
		EngineCallClient: engine.NewClient(server.URL, server.Client()),
	})
	if err := logic.RunAgentStream(&pb.WindAgentRunReq{UserId: 1, Input: "question"}, stream); err != nil {
		t.Fatalf("RunAgentStream() error = %v", err)
	}
	if len(stream.events) != 1 || stream.events[0].Type != "error" || stream.events[0].ErrorMsg != "agent failed" {
		t.Fatalf("events = %+v", stream.events)
	}
}

func TestResumeAgentStreamForwardsEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/agent/resume/stream" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `data: {"type":"done","agent_session_id":"conv-1","content":"approved"}`+"\n\n")
	}))
	defer server.Close()

	stream := &testWindAgentStream{}
	logic := NewResumeAgentStreamLogic(context.Background(), &svc.ServiceContext{
		EngineCallClient: engine.NewClient(server.URL, server.Client()),
	})
	if err := logic.ResumeAgentStream(&pb.WindAgentResumeReq{
		UserId:         1,
		AgentSessionId: "conv-1",
		Action:         " approve ",
	}, stream); err != nil {
		t.Fatalf("ResumeAgentStream() error = %v", err)
	}
	if len(stream.events) != 1 || stream.events[0].Type != "done" || stream.events[0].Content != "approved" {
		t.Fatalf("events = %+v", stream.events)
	}
}

func TestRunAgentStreamReturnsSendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `data: {"type":"token","content":"x"}`+"\n\n")
	}))
	defer server.Close()

	sendErr := errors.New("grpc send failed")
	logic := NewRunAgentStreamLogic(context.Background(), &svc.ServiceContext{
		EngineCallClient: engine.NewClient(server.URL, server.Client()),
	})
	err := logic.RunAgentStream(&pb.WindAgentRunReq{UserId: 1, Input: "question"}, &testWindAgentStream{sendErr: sendErr})
	if !errors.Is(err, sendErr) {
		t.Fatalf("error = %v, want %v", err, sendErr)
	}
}
