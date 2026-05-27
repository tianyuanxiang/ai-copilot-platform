package aichatservicelogic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
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
	maxTitleRunes        = 80
	maxSummaryRunes      = 1800
	insufficientEvidence = "当前知识库没有足够依据回答该问题。"

	answerModeRag  = "rag"
	answerModeChat = "chat"
	answerModeAuto = "auto"
	answerModeDeep = "deep"

	answerPolicyRag    = "rag"
	answerPolicyChat   = "chat"
	answerPolicyAuto   = "auto"
	retrievalModeQuick = "quick"
	retrievalModeDeep  = "deep"
	minRerankScore     = 0.35
)

type RagChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type chatModeDecision struct {
	AnswerPolicy     string
	RetrievalMode    string
	SearchAnswerMode string
}

type ragChatRun struct {
	Question          string
	TraceID           string
	ConversationID    int64
	Conversation      *model.AiConversation
	RecentMessages    []model.AiMessage
	EffectiveKbID     int64
	HasEffectiveKbID  bool
	AnswerPolicy      string
	RetrievalMode     string
	SearchMode        string
	ResponseMode      string
	CandidateChunks   []*pb.ChunkItem
	EffectiveChunks   []*pb.ChunkItem
	Citations         []*pb.Citation
	CitationsJSON     string
	Prompt            string
	ShouldUseRag      bool
	ShouldCallLLM     bool
	TopCandidateScore float64
}

func NewRagChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RagChatLogic {
	return &RagChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RagChatLogic) RagChat(in *pb.RagChatReq) (*pb.RagChatResp, error) {
	run, err := l.prepareRagChat(in)
	if err != nil {
		return nil, err
	}

	answer := insufficientEvidence
	status := "success"
	errorMsg := ""
	startedAt := time.Now()
	if run.ShouldCallLLM {
		chatResp, callErr := l.svcCtx.EngineCallClient.EngineChatStreamAggregate(l.ctx, engine.ChatStreamRequest{
			UserID:         strconv.FormatInt(in.UserId, 10),
			KbID:           optionalInt64String(run.EffectiveKbID, run.HasEffectiveKbID),
			ConversationID: strconv.FormatInt(run.ConversationID, 10),
			Question:       run.Prompt,
			History:        chatHistoryToEngine(run.RecentMessages),
		})
		if callErr != nil {
			status = "failed"
			errorMsg = callErr.Error()
			l.Logger.Errorf("LLM chat call err: %v", callErr)
			l.writeLlmCallLog(run.TraceID, in.UserId, "chat-stream", run.Prompt, "", startedAt, status, errorMsg)
			return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "调用LLM失败")
		}
		answer = strings.TrimSpace(chatResp.Answer)
		if answer == "" {
			if run.ShouldUseRag {
				answer = insufficientEvidence
			} else {
				answer = "暂时无法生成回答，请稍后重试。"
			}
		}
	}

	if err := l.finishRagChat(in.UserId, run, answer, startedAt, status, errorMsg); err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "对话完成后更新对话后的信息失败")
	}

	return &pb.RagChatResp{
		Answer:         answer,
		Citations:      run.Citations,
		TraceId:        run.TraceID,
		Mode:           run.ResponseMode,
		ConversationId: strconv.FormatInt(run.ConversationID, 10),
	}, nil
}

func (l *RagChatLogic) resolveConversation(in *pb.RagChatReq, question string) (*model.AiConversation, error) {
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
			l.Logger.Errorf("find conversation id failed: %v", err)
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
		l.Logger.Errorf("insert conversation id failed: %v", err)
		return nil, err
	}
	return &model.AiConversation{
		Id:     conversationID,
		UserId: in.UserId,
		KbId: sql.NullInt64{
			Int64: in.KbId,
			Valid: hasKbID,
		},
		Title:               truncateRunes(question, maxTitleRunes),
		ConversationSummary: "",
	}, nil
}

func effectiveConversationKb(in *pb.RagChatReq, conversation *model.AiConversation) (int64, bool) {
	if (in.HasKbId || in.KbId > 0) && in.KbId > 0 {
		return in.KbId, true
	}
	if conversation != nil && conversation.KbId.Valid && conversation.KbId.Int64 > 0 {
		return conversation.KbId.Int64, true
	}
	return 0, false
}

func (l *RagChatLogic) prepareRagChat(in *pb.RagChatReq) (*ragChatRun, error) {
	question := strings.TrimSpace(in.Question)
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if question == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "question 不能为空")
	}

	decision := resolveChatMode(in.AnswerMode)
	traceID := ragTraceID()
	conversation, err := l.resolveConversation(in, question)
	if err != nil {
		return nil, err
	}
	conversationID := conversation.Id
	effectiveKbID, hasEffectiveKbID := effectiveConversationKb(in, conversation)

	// 拿到历史消息
	recentMessages, err := l.svcCtx.AiMessageModel.ListRecentByConversation(l.ctx, conversationID, in.UserId, recentMessageLimit)
	if err != nil {
		return nil, err
	}

	// 插入用户问题
	if _, err := l.svcCtx.AiMessageModel.Insert(l.ctx, &model.AiMessage{
		ConversationId: conversationID,
		Role:           "user",
		Content:        question,
		Citations:      "[]",
		TraceId:        traceID,
	}); err != nil {
		return nil, err
	}

	searchMode := ""
	candidateChunks := []*pb.ChunkItem{}
	if decision.AnswerPolicy != answerPolicyChat {
		searchResp, err := aiknowledgeservicelogic.NewSearchKnowledgeLogic(l.ctx, l.svcCtx).SearchKnowledge(&pb.SearchKnowledgeReq{
			UserId:      in.UserId,
			KbId:        effectiveKbID,
			HasKbId:     hasEffectiveKbID,
			Query:       question,
			TopK:        defaultRagTopK,
			AnswerMode:  decision.SearchAnswerMode,
			SearchScope: in.SearchScope,
			DomainId:    in.DomainId,
			HasDomainId: in.HasDomainId && in.DomainId > 0,
			DocumentIds: in.DocumentIds,
		})
		if err != nil {
			return nil, err
		}
		if searchResp != nil {
			searchMode = strings.TrimSpace(searchResp.Mode)
			candidateChunks = searchResp.Chunks
		}
	}

	effectiveChunks := selectEffectiveChunks(question, candidateChunks)

	shouldUseRag := decision.AnswerPolicy != answerPolicyChat && len(effectiveChunks) > 0
	shouldCallLLM := shouldUseRag || decision.AnswerPolicy == answerPolicyChat || decision.AnswerPolicy == answerPolicyAuto
	citations := []*pb.Citation(nil)
	if shouldUseRag {
		citations = citationsFromChunks(effectiveChunks)
	}

	prompt := buildChatPrompt(question, conversation.ConversationSummary, recentMessages)
	if shouldUseRag || decision.AnswerPolicy == answerPolicyRag {
		prompt = buildRagPrompt(question, conversation.ConversationSummary, recentMessages, effectiveChunks)
	}

	run := &ragChatRun{
		Question:          question,
		TraceID:           traceID,
		ConversationID:    conversationID,
		Conversation:      conversation,
		RecentMessages:    recentMessages,
		EffectiveKbID:     effectiveKbID,
		HasEffectiveKbID:  hasEffectiveKbID,
		AnswerPolicy:      decision.AnswerPolicy,
		RetrievalMode:     decision.RetrievalMode,
		SearchMode:        searchMode,
		ResponseMode:      responseMode(searchMode, decision, shouldUseRag),
		CandidateChunks:   candidateChunks,
		EffectiveChunks:   effectiveChunks,
		Citations:         citations,
		CitationsJSON:     marshalCitations(citations),
		Prompt:            prompt,
		ShouldUseRag:      shouldUseRag,
		ShouldCallLLM:     shouldCallLLM,
		TopCandidateScore: topChunkScore(candidateChunks),
	}
	l.Logger.Infof(
		"rag_chat.mode trace_id=%s answerMode=%s answerPolicy=%s retrievalMode=%s candidate_chunks=%d effective_chunks=%d top_score=%.4f response_mode=%s should_call_llm=%t should_use_rag=%t",
		traceID,
		strings.TrimSpace(in.AnswerMode),
		run.AnswerPolicy,
		run.RetrievalMode,
		len(run.CandidateChunks),
		len(run.EffectiveChunks),
		run.TopCandidateScore,
		run.ResponseMode,
		run.ShouldCallLLM,
		run.ShouldUseRag,
	)
	return run, nil
}

func (l *RagChatLogic) finishRagChat(userID int64, run *ragChatRun, answer string, startedAt time.Time, status string, errorMsg string) error {
	if _, err := l.svcCtx.AiMessageModel.Insert(l.ctx, &model.AiMessage{
		ConversationId: run.ConversationID,
		Role:           "assistant",
		Content:        answer,
		Citations:      run.CitationsJSON,
		TraceId:        run.TraceID,
	}); err != nil {
		l.Logger.Errorf("对话完成后插入AI助手返回的数据失败: %v", err)
		return err
	}
	l.writeLlmCallLog(run.TraceID, userID, "chat-stream", run.Prompt, answer, startedAt, status, errorMsg)

	if err := l.svcCtx.AiConversationModel.TouchUpdatedAt(l.ctx, run.ConversationID); err != nil {
		l.Logger.Errorf("touch ai_conversation updated_at failed: %v", err)
	}
	if err := l.refreshConversationSummary(run.ConversationID, userID); err != nil {
		l.Errorf("refresh conversation summary failed: %v", err)
	}
	return nil
}

func resolveChatMode(mode string) chatModeDecision {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case answerModeChat:
		return chatModeDecision{
			AnswerPolicy:     answerPolicyChat,
			RetrievalMode:    retrievalModeQuick, // 不检索知识库
			SearchAnswerMode: "",
		}
	case answerModeAuto:
		return chatModeDecision{
			AnswerPolicy:     answerPolicyAuto,
			RetrievalMode:    retrievalModeQuick, // quick 检索
			SearchAnswerMode: retrievalModeQuick, // 如果存在有效 chunks 走 RAG，无有效 chunks 走普通聊天
		}
	case answerModeDeep:
		return chatModeDecision{
			AnswerPolicy:     answerPolicyRag,
			RetrievalMode:    retrievalModeDeep, // deep 检索 必走rerank
			SearchAnswerMode: answerModeDeep,
		}
	default:
		return chatModeDecision{
			AnswerPolicy:     answerPolicyRag,
			RetrievalMode:    retrievalModeQuick,
			SearchAnswerMode: retrievalModeQuick,
		}
	}
}

func responseMode(searchMode string, decision chatModeDecision, usedRag bool) string {
	if decision.AnswerPolicy == answerPolicyChat {
		return answerModeChat
	}
	if decision.AnswerPolicy == answerPolicyAuto && !usedRag {
		return answerModeChat
	}
	if strings.TrimSpace(searchMode) != "" {
		return searchMode
	}
	if decision.RetrievalMode == retrievalModeDeep {
		return answerModeDeep
	}
	return answerModeRag
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
		l.Logger.Errorf("engine chat_stream_aggregate err: %v", err)
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

func buildChatPrompt(question string, summary string, recentMessages []model.AiMessage) string {
	var builder strings.Builder
	builder.WriteString("系统约束：\n")
	builder.WriteString("- 你是企业 AI 助手，请直接、准确、结构化地回答用户问题。\n")
	builder.WriteString("- 不要编造不存在的事实、数据、系统权限或来源。\n")
	builder.WriteString("- 如果问题需要业务文档依据但当前没有提供，请明确说明需要补充资料。\n\n")

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

func selectEffectiveChunks(question string, chunks []*pb.ChunkItem) []*pb.ChunkItem {
	if len(chunks) == 0 {
		return []*pb.ChunkItem{}
	}
	terms := evidenceTerms(question)
	selected := make([]*pb.ChunkItem, 0, len(chunks))
	for _, chunk := range chunks {
		if chunk == nil {
			continue
		}
		if chunk.Score >= minRerankScore {
			selected = append(selected, chunk)
			continue
		}
		if len(terms) == 0 {
			continue
		}
		corpus := strings.ToLower(chunk.Title + "\n" + chunk.Snippet + "\n" + chunk.Content)
		for _, term := range terms {
			if strings.Contains(corpus, term) {
				selected = append(selected, chunk)
				break
			}
		}
	}
	return selected
}

func evidenceTerms(question string) []string {
	seen := make(map[string]struct{})
	terms := make([]string, 0)
	add := func(term string) {
		term = strings.ToLower(strings.TrimSpace(term))
		if utf8.RuneCountInString(term) <= 1 {
			return
		}
		if _, ok := seen[term]; ok {
			return
		}
		seen[term] = struct{}{}
		terms = append(terms, term)
	}

	tokens := strings.FieldsFunc(question, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for _, token := range tokens {
		add(token)
		runes := []rune(strings.TrimSpace(token))
		if len(runes) <= 2 {
			continue
		}
		for size := minInt(4, len(runes)); size >= 2; size-- {
			for i := 0; i+size <= len(runes); i++ {
				add(string(runes[i : i+size]))
			}
		}
	}
	return terms
}

func topChunkScore(chunks []*pb.ChunkItem) float64 {
	top := 0.0
	for _, chunk := range chunks {
		if chunk != nil && chunk.Score > top {
			top = chunk.Score
		}
	}
	return top
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

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func ragTraceID() string {
	return fmt.Sprintf("rag-%d", time.Now().UnixNano())
}
