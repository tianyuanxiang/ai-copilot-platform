package aichatservicelogic

import (
	"testing"

	"ai-copilot-platform/ai-rpc/pb"
)

func TestResolveChatModeKeepsDeepAsStrictRagWithRerank(t *testing.T) {
	decision := resolveChatMode("deep")

	if decision.AnswerPolicy != answerPolicyRag {
		t.Fatalf("AnswerPolicy = %q, want %q", decision.AnswerPolicy, answerPolicyRag)
	}
	if decision.RetrievalMode != retrievalModeDeep {
		t.Fatalf("RetrievalMode = %q, want %q", decision.RetrievalMode, retrievalModeDeep)
	}
	if decision.SearchAnswerMode != answerModeDeep {
		t.Fatalf("SearchAnswerMode = %q, want %q", decision.SearchAnswerMode, answerModeDeep)
	}
}

func TestResolveChatModeDefaultsInvalidToStrictRag(t *testing.T) {
	decision := resolveChatMode("surprise")

	if decision.AnswerPolicy != answerPolicyRag {
		t.Fatalf("AnswerPolicy = %q, want %q", decision.AnswerPolicy, answerPolicyRag)
	}
	if decision.RetrievalMode != retrievalModeQuick {
		t.Fatalf("RetrievalMode = %q, want %q", decision.RetrievalMode, retrievalModeQuick)
	}
}

func TestSelectEffectiveChunksFiltersIrrelevantWindChunk(t *testing.T) {
	chunks := []*pb.ChunkItem{
		{
			Title:   "风机混塔智能运维 SOP",
			Snippet: "风机混塔巡检、知识库检索、SOP 处置流程",
			Content: "风机混塔智能运维 AI Copilot 的实施、知识库检索和 SOP。",
			Score:   0.016,
		},
	}

	selected := selectEffectiveChunks("减脂餐怎么做", chunks)
	if len(selected) != 0 {
		t.Fatalf("selected %d chunks, want 0", len(selected))
	}
}

func TestSelectEffectiveChunksKeepsTermMatchedChunk(t *testing.T) {
	chunks := []*pb.ChunkItem{
		{
			Title:   "风机混塔智能运维 SOP",
			Snippet: "风机混塔巡检、知识库检索、SOP 处置流程",
			Content: "风机混塔智能运维 AI Copilot 的实施、知识库检索和 SOP。",
			Score:   0.016,
		},
	}

	selected := selectEffectiveChunks("风机混塔 SOP 怎么处理", chunks)
	if len(selected) != 1 {
		t.Fatalf("selected %d chunks, want 1", len(selected))
	}
}
