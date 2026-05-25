package engine

import (
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
	return c.EngineEmbed(ctx, texts)
}

func (c *Client) EngineEmbed(ctx context.Context, texts []string) (*EmbedResponse, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("engine base url is empty")
	}
	body, err := json.Marshal(map[string][]string{"texts": texts})
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
		return nil, err
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
