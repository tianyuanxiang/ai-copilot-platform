package aichatservicelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"ai-copilot-platform/ai-rpc/internal/engine"
	aiknowledgeservicelogic "ai-copilot-platform/ai-rpc/internal/logic/aiknowledgeservice"
	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	defaultRagTopK       = 5
	recentMessageLimit   = 8
	maxPromptRunes       = 6000
	maxPromptChunkRunes  = 800
	maxTitleRunes        = 40
	maxSummaryRunes      = 1800
	insufficientEvidence = "当前知识库没有足够依据回答该问题。"
)

type RagChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRagChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RagChatLogic {
	return &RagChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RagChatLogic) RagChat(in *pb.RagChatReq) (*pb.RagChatResp, error) {
	question := strings.TrimSpace(in.Question)
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if question == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "question 不能为空")
	}

	traceID := ragTraceID()
	conversation, err := l.resolveConversation(in, question)
	if err != nil {
		return nil, err
	}
	conversationID := conversation.Id
	effectiveKbID, hasEffectiveKbID := effectiveConversationKb(in, conversation)

	recentMessages, err := l.svcCtx.AiMessageModel.ListRecentByConversation(l.ctx, conversationID, in.UserId, recentMessageLimit)
	if err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.AiMessageModel.Insert(l.ctx, &model.AiMessage{
		ConversationId: conversationID,
		Role:           "user",
		Content:        question,
		Citations:      "[]",
		TraceId:        traceID,
	}); err != nil {
		return nil, err
	}

	searchResp, err := aiknowledgeservicelogic.NewSearchKnowledgeLogic(l.ctx, l.svcCtx).SearchKnowledge(&pb.SearchKnowledgeReq{
		UserId:      in.UserId,
		KbId:        effectiveKbID,
		HasKbId:     hasEffectiveKbID,
		Query:       question,
		TopK:        defaultRagTopK,
		AnswerMode:  in.AnswerMode,
		SearchScope: in.SearchScope,
		DomainId:    in.DomainId,
		HasDomainId: in.HasDomainId && in.DomainId > 0,
		DocumentIds: in.DocumentIds,
	})
	if err != nil {
		return nil, err
	}

	mode := strings.TrimSpace(searchResp.Mode)
	if mode == "" {
		mode = "rag"
	}
	citations := citationsFromChunks(searchResp.Chunks)
	citationsJSON := marshalCitations(citations)
	prompt := buildRagPrompt(question, conversation.ConversationSummary, recentMessages, searchResp.Chunks)

	answer := insufficientEvidence
	status := "success"
	errorMsg := ""
	startedAt := time.Now()
	if len(searchResp.Chunks) > 0 {
		chatResp, callErr := l.svcCtx.EngineCallClient.EngineChatStreamAggregate(l.ctx, engine.ChatStreamRequest{
			UserID:         strconv.FormatInt(in.UserId, 10),
			KbID:           optionalInt64String(effectiveKbID, hasEffectiveKbID),
			ConversationID: strconv.FormatInt(conversationID, 10),
			Question:       prompt,
			History:        chatHistoryToEngine(recentMessages),
		})
		if callErr != nil {
			status = "failed"
			errorMsg = callErr.Error()
			l.writeLlmCallLog(traceID, in.UserId, "chat-stream", prompt, "", startedAt, status, errorMsg)
			return nil, callErr
		}
		answer = strings.TrimSpace(chatResp.Answer)
		if answer == "" {
			answer = insufficientEvidence
		}
	}

	if _, err := l.svcCtx.AiMessageModel.Insert(l.ctx, &model.AiMessage{
		ConversationId: conversationID,
		Role:           "assistant",
		Content:        answer,
		Citations:      citationsJSON,
		TraceId:        traceID,
	}); err != nil {
		return nil, err
	}
	l.writeLlmCallLog(traceID, in.UserId, "chat-stream", prompt, answer, startedAt, status, errorMsg)

	if err := l.svcCtx.AiConversationModel.TouchUpdatedAt(l.ctx, conversationID); err != nil {
		l.Errorf("touch ai_conversation updated_at failed: %v", err)
	}
	if err := l.refreshConversationSummary(conversationID, in.UserId); err != nil {
		l.Errorf("refresh conversation summary failed: %v", err)
	}

	return &pb.RagChatResp{
		Answer:         answer,
		Citations:      citations,
		TraceId:        traceID,
		Mode:           mode,
		ConversationId: strconv.FormatInt(conversationID, 10),
	}, nil
}

func (l *RagChatLogic) resolveConversation(in *pb.RagChatReq, question string) (*model.AiConversationContext, error) {
	rawConversationID := strings.TrimSpace(in.ConversationId)
	if rawConversationID != "" {
		conversationID, err := strconv.ParseInt(rawConversationID, 10, 64)
		if err != nil || conversationID <= 0 {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "conversationId 必须是有效数字")
		}
		conversation, err := l.svcCtx.AiConversationModel.FindContextByIDUserID(l.ctx, conversationID, in.UserId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "会话不存在")
			}
			return nil, err
		}
		return conversation, nil
	}

	hasKbID := (in.HasKbId || in.KbId > 0) && in.KbId > 0
	conversationID, err := l.svcCtx.AiConversationModel.InsertReturningID(l.ctx, &model.AiConversation{
		UserId: in.UserId,
		KbId: sql.NullInt64{
			Int64: in.KbId,
			Valid: hasKbID,
		},
		Title: truncateRunes(question, maxTitleRunes),
	})
	if err != nil {
		return nil, err
	}
	return &model.AiConversationContext{
		Id:                  conversationID,
		UserId:              in.UserId,
		KbId:                in.KbId,
		HasKbId:             hasKbID,
		Title:               truncateRunes(question, maxTitleRunes),
		ConversationSummary: "",
	}, nil
}

func effectiveConversationKb(in *pb.RagChatReq, conversation *model.AiConversationContext) (int64, bool) {
	if (in.HasKbId || in.KbId > 0) && in.KbId > 0 {
		return in.KbId, true
	}
	if conversation != nil && conversation.HasKbId && conversation.KbId > 0 {
		return conversation.KbId, true
	}
	return 0, false
}

func (l *RagChatLogic) refreshConversationSummary(conversationID int64, userID int64) error {
	messages, err := l.svcCtx.AiMessageModel.ListActiveByConversation(l.ctx, conversationID, userID)
	if err != nil {
		return err
	}
	olderCount := len(messages) - recentMessageLimit
	if olderCount <= 0 {
		return l.svcCtx.AiConversationModel.UpdateSummaryByIDUserID(l.ctx, conversationID, userID, "")
	}

	prompt := buildSummaryPrompt(messages[:olderCount])
	startedAt := time.Now()
	traceID := fmt.Sprintf("summary-%d", time.Now().UnixNano())
	summaryResp, err := l.svcCtx.EngineCallClient.EngineChatStreamAggregate(l.ctx, engine.ChatStreamRequest{
		UserID:         strconv.FormatInt(userID, 10),
		ConversationID: strconv.FormatInt(conversationID, 10),
		Question:       prompt,
	})
	if err != nil {
		l.writeLlmCallLog(traceID, userID, "chat-summary", prompt, "", startedAt, "failed", err.Error())
		return err
	}

	summary := truncateRunes(strings.TrimSpace(summaryResp.Answer), maxSummaryRunes)
	if summary == "" {
		summary = ""
	}
	l.writeLlmCallLog(traceID, userID, "chat-summary", prompt, summary, startedAt, "success", "")
	return l.svcCtx.AiConversationModel.UpdateSummaryByIDUserID(l.ctx, conversationID, userID, summary)
}

func (l *RagChatLogic) writeLlmCallLog(traceID string, userID int64, modelName string, prompt string, answer string, startedAt time.Time, status string, errorMsg string) {
	_, err := l.svcCtx.AiLlmCallLogModel.Insert(l.ctx, &model.AiLlmCallLog{
		UserId: sql.NullInt64{
			Int64: userID,
			Valid: userID > 0,
		},
		TraceId:          traceID,
		Provider:         "python-engine",
		Model:            modelName,
		Prompt:           prompt,
		PromptTokens:     estimateWhitespaceTokens(prompt),
		CompletionTokens: estimateWhitespaceTokens(answer),
		LatencyMs:        time.Since(startedAt).Milliseconds(),
		Status:           status,
		ErrorMsg:         truncateRunes(errorMsg, 1000),
	})
	if err != nil {
		l.Errorf("write ai_llm_call_log failed: %v", err)
	}
}

func buildRagPrompt(question string, summary string, recentMessages []model.AiMessage, chunks []*pb.ChunkItem) string {
	var builder strings.Builder
	builder.WriteString("系统约束：\n")
	builder.WriteString("- 只能基于提供的文档片段回答。\n")
	builder.WriteString("- 如果文档中没有足够依据，必须明确说明“当前知识库没有足够依据”。\n")
	builder.WriteString("- 不得编造来源、文档、数据或结论。\n")
	builder.WriteString("- 回答尽量结构化，并在相关结论后保留引用编号。\n\n")

	if strings.TrimSpace(summary) != "" {
		builder.WriteString("长期对话摘要：\n")
		builder.WriteString(truncateRunes(strings.TrimSpace(summary), maxSummaryRunes))
		builder.WriteString("\n\n")
	}

	if len(recentMessages) > 0 {
		builder.WriteString("近期对话：\n")
		for _, item := range recentMessages {
			role := strings.TrimSpace(item.Role)
			content := strings.TrimSpace(item.Content)
			if role == "" || content == "" {
				continue
			}
			builder.WriteString("- ")
			builder.WriteString(role)
			builder.WriteString(": ")
			builder.WriteString(truncateRunes(content, 260))
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("命中文档片段：\n")
	if len(chunks) == 0 {
		builder.WriteString("无\n\n")
	} else {
		for index, chunk := range chunks {
			if chunk == nil {
				continue
			}
			builder.WriteString(fmt.Sprintf("[%d] title=%s document_id=%d chunk_id=%d score=%.4f\ncontent=%s\n\n",
				index+1,
				truncateRunes(chunk.Title, 80),
				chunk.DocumentId,
				chunk.ChunkId,
				chunk.Score,
				truncateRunes(firstNonEmpty(chunk.Content, chunk.Snippet), maxPromptChunkRunes),
			))
		}
	}

	builder.WriteString("当前用户问题：\n")
	builder.WriteString(question)
	builder.WriteString("\n")

	return truncateRunes(builder.String(), maxPromptRunes)
}

func buildSummaryPrompt(messages []model.AiMessage) string {
	var builder strings.Builder
	builder.WriteString("请为以下已归档的对话生成企业知识问答场景下的长期上下文摘要。\n")
	builder.WriteString("要求：只保留对后续问答有帮助的事实、用户偏好、已确认结论和待解决问题；删除寒暄和重复内容；不要引入原文没有的信息；只输出摘要正文。\n\n")
	builder.WriteString("对话内容：\n")
	for _, item := range messages {
		role := strings.TrimSpace(item.Role)
		content := strings.TrimSpace(item.Content)
		if role == "" || content == "" {
			continue
		}
		builder.WriteString(role)
		builder.WriteString(": ")
		builder.WriteString(truncateRunes(content, 500))
		builder.WriteString("\n")
	}
	return truncateRunes(builder.String(), maxPromptRunes)
}

func chatHistoryToEngine(messages []model.AiMessage) []engine.ChatMessage {
	items := make([]engine.ChatMessage, 0, len(messages))
	for _, item := range messages {
		role := strings.TrimSpace(item.Role)
		content := strings.TrimSpace(item.Content)
		if role == "" || content == "" {
			continue
		}
		items = append(items, engine.ChatMessage{
			Role:    role,
			Content: content,
		})
	}
	return items
}

func citationsFromChunks(chunks []*pb.ChunkItem) []*pb.Citation {
	citations := make([]*pb.Citation, 0, len(chunks))
	for _, chunk := range chunks {
		if chunk == nil {
			continue
		}
		citations = append(citations, &pb.Citation{
			DocumentId: chunk.DocumentId,
			ChunkId:    chunk.ChunkId,
			Title:      chunk.Title,
			Snippet:    firstNonEmpty(chunk.Snippet, chunk.Content),
			Score:      chunk.Score,
		})
	}
	return citations
}

func marshalCitations(citations []*pb.Citation) string {
	type citationJSON struct {
		DocumentID int64   `json:"document_id"`
		ChunkID    int64   `json:"chunk_id"`
		Title      string  `json:"title"`
		Snippet    string  `json:"snippet"`
		Score      float64 `json:"score"`
	}
	items := make([]citationJSON, 0, len(citations))
	for _, item := range citations {
		items = append(items, citationJSON{
			DocumentID: item.DocumentId,
			ChunkID:    item.ChunkId,
			Title:      item.Title,
			Snippet:    item.Snippet,
			Score:      item.Score,
		})
	}
	data, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func parseCitationsJSON(raw string) []*pb.Citation {
	type citationJSON struct {
		DocumentID int64   `json:"document_id"`
		ChunkID    int64   `json:"chunk_id"`
		Title      string  `json:"title"`
		Snippet    string  `json:"snippet"`
		Score      float64 `json:"score"`
	}
	var items []citationJSON
	if err := json.Unmarshal([]byte(firstNonEmpty(raw, "[]")), &items); err != nil {
		return nil
	}
	citations := make([]*pb.Citation, 0, len(items))
	for _, item := range items {
		citations = append(citations, &pb.Citation{
			DocumentId: item.DocumentID,
			ChunkId:    item.ChunkID,
			Title:      item.Title,
			Snippet:    item.Snippet,
			Score:      item.Score,
		})
	}
	return citations
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func optionalInt64String(value int64, ok bool) string {
	if !ok || value <= 0 {
		return ""
	}
	return strconv.FormatInt(value, 10)
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 || utf8.RuneCountInString(value) <= limit {
		return value
	}
	runes := []rune(value)
	return string(runes[:limit])
}

func estimateWhitespaceTokens(value string) int64 {
	count := len(strings.Fields(value))
	if count == 0 && strings.TrimSpace(value) != "" {
		return 1
	}
	return int64(count)
}

func ragTraceID() string {
	return fmt.Sprintf("rag-%d", time.Now().UnixNano())
}
