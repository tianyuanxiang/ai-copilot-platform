package ai_chat

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"ai-copilot-platform/ai-rpc/pb"
	"ai-copilot-platform/gateway/internal/logic/ai_wind_tool"

	_ "github.com/taosdata/driver-go/v3/taosWS"
)

func TestStreamCitationsFromRPCAddsReferenceLabels(t *testing.T) {
	citations := ai_wind_tool.StreamCitationsFromRPC([]*pb.Citation{
		{
			DocumentId: 10,
			ChunkId:    20,
			Title:      "deploy.md",
			Snippet:    "kubectl get pod -n blade",
			Score:      0.87,
		},
		{
			DocumentId: 11,
			ChunkId:    21,
			Title:      "ops.md",
			Snippet:    "查看 Pod 是否启动成功",
			Score:      0.76,
		},
	})

	if len(citations) != 2 {
		t.Fatalf("len(citations) = %d, want 2", len(citations))
	}
	if citations[0].RefIndex != 1 || citations[0].RefText != "[1]" {
		t.Fatalf("first citation ref = (%d, %q), want (1, %q)", citations[0].RefIndex, citations[0].RefText, "[1]")
	}
	if citations[1].RefIndex != 2 || citations[1].RefText != "[2]" {
		t.Fatalf("second citation ref = (%d, %q), want (2, %q)", citations[1].RefIndex, citations[1].RefText, "[2]")
	}
}

func TestStreamReferencesFromAnswerMapsMarkersToCitations(t *testing.T) {
	citations := []ai_wind_tool.StreamCitation{
		{RefIndex: 1, RefText: "[1]", DocumentId: 10, ChunkId: 20, Title: "deploy.md"},
		{RefIndex: 2, RefText: "[2]", DocumentId: 11, ChunkId: 21, Title: "ops.md"},
	}

	references := ai_wind_tool.StreamReferencesFromAnswer("命令用于查看 Pod[1]，也可检查启动状态[2]。", citations)
	if len(references) != 2 {
		t.Fatalf("len(references) = %d, want 2", len(references))
	}
	if references[0].RefIndex != 1 || references[0].Citation.DocumentId != 10 {
		t.Fatalf("first reference = %+v, want refIndex=1 documentId=10", references[0])
	}
	if references[1].RefIndex != 2 || references[1].Citation.DocumentId != 11 {
		t.Fatalf("second reference = %+v, want refIndex=2 documentId=11", references[1])
	}
	if references[0].Start >= references[0].End {
		t.Fatalf("invalid first reference range: start=%d end=%d", references[0].Start, references[0].End)
	}
}

func TestTDConnection(t *testing.T) {
	db, err := sql.Open("taosWS", "root:bk147258.@ws(172.16.90.50:6041)/")
	if err != nil {
		panic("Failed to connect to " + "; ErrMessage: " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var v int
	err = db.QueryRowContext(ctx, "select 1").Scan(&v)
	fmt.Printf("query: v=%d err=%v\n", v, err)
}
