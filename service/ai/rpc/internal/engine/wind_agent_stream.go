package engine

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const windAgentSSEMaxEventBytes = 1024 * 1024

type WindAgentRunRequest struct {
	UserID         int64  `json:"user_id"`
	ConversationID string `json:"conversation_id,omitempty"`
	Input          string `json:"input"`
}

type WindAgentResumeRequest struct {
	UserID         int64  `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Action         string `json:"action"`
	Content        string `json:"content,omitempty"`
}

type WindAgentToolCall struct {
	ToolCallID    int64  `json:"tool_call_id"`
	ToolName      string `json:"tool_name"`
	Status        string `json:"status"`
	ArgumentsJSON string `json:"arguments_json"`
	ResultJSON    string `json:"result_json"`
	Message       string `json:"message"`
	LatencyMS     int64  `json:"latency_ms"`
}

type WindAgentCitation struct {
	DocumentID int64   `json:"document_id"`
	ChunkID    int64   `json:"chunk_id"`
	Title      string  `json:"title"`
	Snippet    string  `json:"snippet"`
	Score      float64 `json:"score"`
}

type WindAgentDraftRef struct {
	DraftType string `json:"draft_type"`
	DraftID   int64  `json:"draft_id"`
	Title     string `json:"title"`
}

type WindAgentStreamEvent struct {
	Type           string              `json:"type"`
	TraceID        string              `json:"trace_id"`
	ConversationID string              `json:"conversation_id"`
	Content        string              `json:"content"`
	ToolCall       *WindAgentToolCall  `json:"tool_call,omitempty"`
	ToolCalls      []WindAgentToolCall `json:"tool_calls,omitempty"`
	Citations      []WindAgentCitation `json:"citations,omitempty"`
	Draft          *WindAgentDraftRef  `json:"draft,omitempty"`
	ErrorMsg       string              `json:"error_msg,omitempty"`
}

// EngineWindAgentStream 启动一次 Agent 会话。Go 只负责转发事件，编排状态仍由 Python Agent 控制。
func (c *Client) EngineWindAgentStream(ctx context.Context, payload WindAgentRunRequest, onEvent func(WindAgentStreamEvent) error) error {
	return c.engineWindAgentStream(ctx, "/v1/agent/stream", payload, onEvent)
}

// EngineWindAgentResumeStream 恢复被澄清或审批中断的 Agent 会话。
func (c *Client) EngineWindAgentResumeStream(ctx context.Context, payload WindAgentResumeRequest, onEvent func(WindAgentStreamEvent) error) error {
	return c.engineWindAgentStream(ctx, "/v1/agent/resume/stream", payload, onEvent)
}

// engineWindAgentStream 统一读取 Python 返回的 SSE，保留请求 context 以便客户端断开时及时取消。
func (c *Client) engineWindAgentStream(ctx context.Context, path string, payload any, onEvent func(WindAgentStreamEvent) error) error {
	if c == nil || c.baseURL == "" {
		return fmt.Errorf("engine base url is empty")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal engine wind agent request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create engine wind agent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request engine wind agent stream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		if readErr != nil {
			return fmt.Errorf("read engine wind agent error response: %w", readErr)
		}
		return fmt.Errorf("engine wind agent stream returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024), windAgentSSEMaxEventBytes)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" || data == "[DONE]" {
			continue
		}

		var event WindAgentStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return fmt.Errorf("decode engine wind agent event: %w", err)
		}
		if onEvent != nil {
			if err := onEvent(event); err != nil {
				return fmt.Errorf("handle engine wind agent event: %w", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read engine wind agent stream: %w", err)
	}
	return nil
}
