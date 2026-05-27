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

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type EmbedResponse struct {
	Vectors     [][]float64 `json:"vectors"`
	TokenCounts []int64     `json:"token_counts"`
	TotalTokens int64       `json:"total_tokens"`
	Model       string      `json:"model"`
	Dimension   int64       `json:"dimension"`
	Mode        string      `json:"mode"`
}

type RerankChunk struct {
	ChunkID    string  `json:"chunk_id"`
	DocumentID string  `json:"document_id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
	Source     string  `json:"source"`
}

type RerankResponse struct {
	Chunks []RerankChunk `json:"chunks"`
	Mode   string        `json:"mode"`
}

type ChildChunk struct {
	ChunkIndex int64  `json:"chunk_index"`
	Content    string `json:"content"`
	TokenCount int64  `json:"token_count"`
}

type ParentChunk struct {
	ParentIndex int64        `json:"parent_index"`
	Content     string       `json:"content"`
	TokenCount  int64        `json:"token_count"`
	Children    []ChildChunk `json:"children"`
}

type ParseResponse struct {
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata"`
	Parents  []ParentChunk  `json:"parents"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatStreamRequest struct {
	UserID         string        `json:"user_id"`
	KbID           string        `json:"kb_id"`
	ConversationID string        `json:"conversation_id"`
	Question       string        `json:"question"`
	History        []ChatMessage `json:"history"`
}

type ChatStreamEvent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	TraceID string `json:"trace_id"`
}

type ChatAggregateResponse struct {
	Answer  string `json:"answer"`
	TraceID string `json:"trace_id"`
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: httpClient,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) Embed(ctx context.Context, texts []string) (*EmbedResponse, error) {
	return c.EngineEmbed(ctx, texts, "document")
}

func (c *Client) EngineEmbed(ctx context.Context, texts []string, inputType ...string) (*EmbedResponse, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("engine base url is empty")
	}
	embedType := "document"
	if len(inputType) > 0 && strings.TrimSpace(inputType[0]) != "" {
		embedType = strings.TrimSpace(inputType[0])
	}
	body, err := json.Marshal(map[string]any{
		"texts":      texts,
		"input_type": embedType,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/embed", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行Embedding 请求 %v 出现错误, err: %v", req, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("engine embed returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var embedded EmbedResponse
	if err := json.Unmarshal(respBody, &embedded); err != nil {
		return nil, err
	}
	return &embedded, nil
}

func (c *Client) EngineRerank(ctx context.Context, query string, chunks []RerankChunk, topK int64) (*RerankResponse, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("engine base url is empty")
	}
	body, err := json.Marshal(map[string]any{
		"query":  query,
		"chunks": chunks,
		"top_k":  topK,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/rerank", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("engine rerank returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var reranked RerankResponse
	if err := json.Unmarshal(respBody, &reranked); err != nil {
		return nil, err
	}
	return &reranked, nil
}

func (c *Client) EngineParse(ctx context.Context, fileName string, fileType string, content string) (*ParseResponse, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("engine base url is empty")
	}
	body, err := json.Marshal(map[string]string{
		"file_name": fileName,
		"file_type": fileType,
		"content":   content,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/parse", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("engine parse returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed ParseResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (c *Client) EngineChatStreamAggregate(ctx context.Context, payload ChatStreamRequest) (*ChatAggregateResponse, error) {
	var providerErr error
	resp, err := c.EngineChatStream(ctx, payload, func(event ChatStreamEvent) error {
		if event.Type == "error" && strings.TrimSpace(event.Content) != "" {
			providerErr = fmt.Errorf("%s", event.Content)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if providerErr != nil {
		return nil, providerErr
	}
	return resp, nil
}

func (c *Client) EngineChatStream(ctx context.Context, payload ChatStreamRequest, onEvent func(ChatStreamEvent) error) (*ChatAggregateResponse, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("engine base url is empty")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/stream", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		return nil, fmt.Errorf("engine chat stream returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var answer strings.Builder
	var traceID string
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" || data == "[DONE]" {
			continue
		}

		var event ChatStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return nil, err
		}
		fmt.Println("event:", event)
		if event.TraceID != "" {
			traceID = event.TraceID
		}
		if event.Type == "token" {
			answer.WriteString(event.Content)
		}
		if onEvent != nil {
			if err := onEvent(event); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ChatAggregateResponse{
		Answer:  answer.String(),
		TraceID: traceID,
	}, nil
}
