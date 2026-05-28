package ai_chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// 修改这里的配置来切换测试环境
const (
	testBaseURL = "http://127.0.0.1:8356"
	testToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJ0b2tlblR5cGUiOiJhY2Nlc3MiLCJleHAiOjE3Nzk5MjIxMDcsIm5iZiI6MTc3OTg1MDEwNywiaWF0IjoxNzc5ODUwMTA3fQ.wRUYRnWtsPbawaQoiZai7btrftT9AXZpa8I4TF3i7-0"
)

type streamCitation struct {
	RefIndex   int     `json:"refIndex"`
	RefText    string  `json:"refText"`
	DocumentId int64   `json:"documentId"`
	ChunkId    int64   `json:"chunkId"`
	Title      string  `json:"title"`
	Snippet    string  `json:"snippet"`
	Score      float64 `json:"score"`
}

// TestAiChatStream 流式问答集成测试
func TestAiChatStream(t *testing.T) {
	reqBody, err := json.Marshal(map[string]any{
		"question":       "你的回答我无法理解",
		"answerMode":     "auto",
		"conversationId": "6",
	})
	if err != nil {
		t.Fatalf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, testBaseURL+"/api/v1/ai/chat/stream", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("发送请求失败（服务是否已启动？）: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("HTTP 状态码: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("期望状态码 200，实际: %d", resp.StatusCode)
	}

	// 逐行读取 SSE 流，实时打印每个 token
	scanner := bufio.NewScanner(resp.Body)
	var fullAnswer strings.Builder
	eventCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}

		var event struct {
			Type      string           `json:"type"`
			Content   string           `json:"content"`
			TraceID   string           `json:"traceId"`
			Citations []streamCitation `json:"citations,omitempty"`
		}

		if err := json.Unmarshal([]byte(data), &event); err != nil {
			t.Logf("解析事件失败，原始数据: %s, err: %v", data, err)
			continue
		}

		eventCount++
		switch event.Type {
		case "token":
			fullAnswer.WriteString(event.Content)
			// fmt.Print(event.Type)
			fmt.Print(event.Content) // 实时打印，看到打字机效果
		case "error":
			fmt.Println()
			t.Fatalf("收到错误事件: %s", event.Content)
		case "done":
			fmt.Println()
			t.Logf("流结束，traceId: %s", event.TraceID)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("读取流失败: %v", err)
	}

	t.Logf("共收到 %d 个事件", eventCount)
	// t.Logf("完整回答:\n%s", fullAnswer.String())

	if fullAnswer.Len() == 0 {
		t.Fatal("回答内容为空")
	}
}
