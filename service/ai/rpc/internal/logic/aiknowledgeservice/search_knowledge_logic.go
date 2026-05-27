package aiknowledgeservicelogic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"ai-copilot-platform/ai-rpc/internal/engine"
	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	hybridDefaultTopK = 5
	hybridMaxTopK     = 20
	quickRecallLimit  = 20
	deepMinRecall     = 40
	rrfK              = 60.0
)

type SearchKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type searchStrategy struct {
	Mode      string
	Recall    int
	UseRerank bool
}

type retrievedCandidate struct {
	ChunkId       int64
	DocumentId    int64
	ParentChunkId int64
	Title         string
	Content       string
	Score         float64
	Source        string
	Rank          int
}

func NewSearchKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchKnowledgeLogic {
	return &SearchKnowledgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchKnowledge executes permission-scoped hybrid retrieval over pgvector and Elasticsearch.
func (l *SearchKnowledgeLogic) SearchKnowledge(in *pb.SearchKnowledgeReq) (*pb.SearchKnowledgeResp, error) {
	query := strings.TrimSpace(in.Query)
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	if query == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "query cannot be empty")
	}

	topK := normalizeTopK(in.TopK)

	strategy := resolveSearchStrategy(in.AnswerMode, topK)

	// 获取可操作的知识库id
	accessibleKbIDs, err := l.resolveAccessibleKbIDs(in)
	if err != nil {
		return nil, err
	}
	if len(accessibleKbIDs) == 0 {
		return &pb.SearchKnowledgeResp{
			Chunks:  []*pb.ChunkItem{},
			Mode:    strategy.Mode,
			Message: "no accessible knowledge bases",
		}, nil
	}

	embedded, err := l.svcCtx.EngineCallClient.EngineEmbed(l.ctx, []string{query}, "query")
	if err != nil {
		return nil, err
	}
	if len(embedded.Vectors) != 1 || len(embedded.Vectors[0]) == 0 {
		return nil, fmt.Errorf("engine embed returned invalid query vector")
	}

	vectorResults, err := l.searchPgvector(embedded.Vectors[0], accessibleKbIDs, in.DocumentIds, strategy.Recall)
	if err != nil {
		return nil, err
	}
	bm25Results, err := l.searchElasticsearch(query, accessibleKbIDs, in.DocumentIds, strategy.Recall)
	if err != nil {
		return nil, err
	}

	candidates := fuseByRRF(vectorResults, bm25Results, strategy.Recall)
	if strategy.UseRerank && len(candidates) > 0 {
		candidates, strategy.Mode, err = l.rerankCandidates(query, candidates, topK)
		if err != nil {
			return nil, err
		}
	} else if len(candidates) > topK {
		candidates = candidates[:topK]
	}

	chunks := make([]*pb.ChunkItem, 0, len(candidates))
	for _, item := range candidates {
		chunks = append(chunks, &pb.ChunkItem{
			DocumentId: item.DocumentId,
			ChunkId:    item.ChunkId,
			Title:      item.Title,
			Snippet:    makeQueryAwareSnippet(query, item.Content, 180),
			Content:    cleanDisplayText(item.Content),
			Score:      item.Score,
			Source:     item.Source,
		})
	}

	message := "ok"
	if len(chunks) == 0 {
		message = "no matched chunks"
	}
	return &pb.SearchKnowledgeResp{
		Chunks:  chunks,
		Mode:    strategy.Mode,
		Message: message,
	}, nil
}

func normalizeTopK(value int64) int {
	if value <= 0 {
		return hybridDefaultTopK
	}
	if value > hybridMaxTopK {
		return hybridMaxTopK
	}
	return int(value)
}

func resolveSearchStrategy(answerMode string, topK int) searchStrategy {
	if strings.EqualFold(strings.TrimSpace(answerMode), "deep") {
		recall := topK * 4
		if recall < deepMinRecall {
			recall = deepMinRecall
		}
		return searchStrategy{
			Mode:      "hybrid-deep-rerank",
			Recall:    recall,
			UseRerank: true,
		}
	}
	return searchStrategy{
		Mode:      "hybrid-quick",
		Recall:    quickRecallLimit,
		UseRerank: false,
	}
}

func (l *SearchKnowledgeLogic) resolveAccessibleKbIDs(in *pb.SearchKnowledgeReq) ([]int64, error) {
	if in.HasKbId || in.KbId > 0 {
		kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "knowledge base not found")
			}
			return nil, err
		}
		if !l.canReadKnowledgeBase(kb, in.UserId) {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "no permission to read knowledge base")
		}
		return []int64{in.KbId}, nil
	}
	// 查出当前用户有权限访问的所有知识库 ID 列表。
	KbIds, err := l.svcCtx.AiKnowledgeBaseModel.FindAccessibleKnowledgeBaseIDsByScope(
		l.ctx,
		in.UserId,
		normalizeSearchScope(in.SearchScope),
		in.DomainId,
		in.HasDomainId && in.DomainId > 0,
	)
	if err != nil {
		l.Logger.Error("")
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "knowledge base not found")
	}

	return KbIds, nil
}

func normalizeSearchScope(scope string) string {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "personal":
		return "personal"
	case "public":
		return "public"
	default:
		return "all"
	}
}

func (l *SearchKnowledgeLogic) canReadKnowledgeBase(kb *model.AiKnowledgeBase, userID int64) bool {
	if kb == nil || kb.Status != 1 {
		return false
	}
	if kb.OwnerUserId.Valid && kb.OwnerUserId.Int64 == userID {
		return true
	}
	if strings.EqualFold(kb.Visibility, "public") {
		return true
	}
	member, err := l.svcCtx.AiKbMemberModel.FindByKbIDUserID(l.ctx, kb.Id, userID)
	if err != nil {
		return false
	}
	return member.Role == "viewer" || member.Role == "editor" || member.Role == "manager"
}

func (l *SearchKnowledgeLogic) searchPgvector(queryVector model.PgVector, kbIDs []int64, documentIDs []int64, limit int) ([]retrievedCandidate, error) {
	rows, err := l.svcCtx.AiDocumentChunkModel.SearchDocumentChunksByVector(l.ctx, l.svcCtx.Orm, queryVector, kbIDs, documentIDs, limit)
	if err != nil {
		return nil, err
	}

	results := make([]retrievedCandidate, 0, len(rows))
	for index, row := range rows {
		results = append(results, retrievedCandidate{
			ChunkId:       row.ChunkId,
			DocumentId:    row.DocumentId,
			ParentChunkId: row.ParentChunkId,
			Title:         row.Title,
			Content:       row.Content,
			Score:         row.Score,
			Source:        "vector",
			Rank:          index + 1,
		})
	}
	return results, nil
}

func (l *SearchKnowledgeLogic) searchElasticsearch(query string, kbIDs []int64, documentIDs []int64, limit int) ([]retrievedCandidate, error) {
	indexName := strings.TrimSpace(l.svcCtx.Config.Elasticsearch.KbChunksIndex)
	if indexName == "" || l.svcCtx.ES == nil {
		return []retrievedCandidate{}, nil
	}
	if limit <= 0 {
		limit = quickRecallLimit
	}

	filters := []map[string]any{
		{"terms": map[string]any{"kb_id": kbIDs}},
	}
	if len(documentIDs) > 0 {
		filters = append(filters, map[string]any{"terms": map[string]any{"document_id": documentIDs}})
	}

	payload := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{"match": map[string]any{"content": query}},
				},
				"filter": filters,
			},
		},
		"size": limit,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := l.svcCtx.ES.Search(
		l.svcCtx.ES.Search.WithContext(l.ctx),
		l.svcCtx.ES.Search.WithIndex(indexName),
		l.svcCtx.ES.Search.WithBody(bytes.NewReader(body)),
		l.svcCtx.ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("elasticsearch search returned %s: %s", resp.Status(), string(respBody))
	}

	var parsed struct {
		Hits struct {
			Hits []struct {
				Score  float64        `json:"_score"`
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	results := make([]retrievedCandidate, 0, len(parsed.Hits.Hits))
	for index, hit := range parsed.Hits.Hits {
		chunkID := anyToInt64(hit.Source["chunk_id"])
		documentID := anyToInt64(hit.Source["document_id"])
		if chunkID <= 0 || documentID <= 0 {
			continue
		}
		results = append(results, retrievedCandidate{
			ChunkId:       chunkID,
			DocumentId:    documentID,
			ParentChunkId: anyToInt64(hit.Source["parent_chunk_id"]),
			Title:         anyToString(hit.Source["title"]),
			Content:       anyToString(hit.Source["content"]),
			Score:         hit.Score,
			Source:        "bm25",
			Rank:          index + 1,
		})
	}
	return results, nil
}

func (l *SearchKnowledgeLogic) rerankCandidates(query string, candidates []retrievedCandidate, topK int) ([]retrievedCandidate, string, error) {
	rerankChunks := make([]engine.RerankChunk, 0, len(candidates))
	candidateByChunkID := make(map[string]retrievedCandidate, len(candidates))
	for _, item := range candidates {
		chunkID := strconv.FormatInt(item.ChunkId, 10)
		candidateByChunkID[chunkID] = item
		rerankChunks = append(rerankChunks, engine.RerankChunk{
			ChunkID:    chunkID,
			DocumentID: strconv.FormatInt(item.DocumentId, 10),
			Title:      item.Title,
			Content:    item.Content,
			Score:      item.Score,
			Source:     item.Source,
		})
	}

	resp, err := l.svcCtx.EngineCallClient.EngineRerank(l.ctx, query, rerankChunks, int64(topK))
	if err != nil {
		return nil, "", err
	}

	reranked := make([]retrievedCandidate, 0, len(resp.Chunks))
	for rank, item := range resp.Chunks {
		candidate, ok := candidateByChunkID[item.ChunkID]
		if !ok {
			continue
		}
		candidate.Score = item.Score
		candidate.Source = mergeSource(candidate.Source, "rerank")
		candidate.Rank = rank + 1
		reranked = append(reranked, candidate)
	}

	mode := "hybrid-deep-rerank"
	if strings.Contains(strings.ToLower(resp.Mode), "mock") {
		mode = "hybrid-deep-rerank-mock"
	}
	if len(reranked) == 0 {
		return candidates[:minInt(len(candidates), topK)], mode, nil
	}
	if len(reranked) > topK {
		reranked = reranked[:topK]
	}
	return reranked, mode, nil
}

// RRF 倒数排名融合
func fuseByRRF(vectorResults []retrievedCandidate, bm25Results []retrievedCandidate, limit int) []retrievedCandidate {
	merged := make(map[int64]retrievedCandidate)
	add := func(item retrievedCandidate) {
		if item.ChunkId <= 0 || item.Rank <= 0 {
			return
		}
		score := 1 / (rrfK + float64(item.Rank))
		existing, ok := merged[item.ChunkId]
		if !ok {
			item.Score = score
			merged[item.ChunkId] = item
			return
		}
		existing.Score += score
		existing.Source = mergeSource(existing.Source, item.Source)
		// Title 和 Content 的兜底
		if existing.Title == "" {
			existing.Title = item.Title
		}
		if existing.Content == "" {
			existing.Content = item.Content
		}
		merged[item.ChunkId] = existing
	}

	for _, item := range vectorResults {
		add(item)
	}
	for _, item := range bm25Results {
		add(item)
	}

	items := make([]retrievedCandidate, 0, len(merged))
	for _, item := range merged {
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		if math.Abs(items[i].Score-items[j].Score) > 1e-12 {
			return items[i].Score > items[j].Score
		}
		return items[i].ChunkId < items[j].ChunkId
	})
	if limit > 0 && len(items) > limit {
		return items[:limit]
	}
	return items
}

func makeQueryAwareSnippet(query string, content string, maxRunes int) string {
	text := cleanDisplayText(content)
	if text == "" {
		return ""
	}

	bestIndex := -1
	lowerText := strings.ToLower(text)
	for _, term := range queryTerms(query) {
		if index := strings.Index(lowerText, strings.ToLower(term)); index >= 0 && (bestIndex == -1 || index < bestIndex) {
			bestIndex = index
		}
	}

	runes := []rune(text)
	if maxRunes <= 0 || len(runes) <= maxRunes {
		return text
	}
	if bestIndex < 0 {
		return firstSentenceSnippet(text, maxRunes)
	}

	prefixRunes := utf8.RuneCountInString(text[:bestIndex])
	start := prefixRunes - maxRunes/3
	if start < 0 {
		start = 0
	}
	end := start + maxRunes
	if end > len(runes) {
		end = len(runes)
		start = maxInt(0, end-maxRunes)
	}

	snippet := strings.TrimSpace(string(runes[start:end]))
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(runes) {
		snippet += "..."
	}
	return snippet
}

func firstSentenceSnippet(text string, maxRunes int) string {
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == 0x3002 || r == 0xff01 || r == 0xff1f || r == '\n'
	})
	if len(sentences) > 0 && strings.TrimSpace(sentences[0]) != "" {
		text = strings.TrimSpace(sentences[0])
	}
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}
	return strings.TrimSpace(string(runes[:maxRunes])) + "..."
}

func cleanDisplayText(content string) string {
	text := strings.TrimSpace(content)
	text = strings.ReplaceAll(text, "```", "")
	text = strings.ReplaceAll(text, "~~~", "")

	lines := make([]string, 0)
	for _, line := range strings.Split(text, "\n") {
		clean := strings.TrimSpace(line)
		if clean == "" || isMarkdownTableSeparator(clean) {
			continue
		}
		lines = append(lines, clean)
	}
	return strings.TrimSpace(strings.Join(strings.Fields(strings.Join(lines, "\n")), " "))
}

func queryTerms(query string) []string {
	raw := strings.FieldsFunc(query, func(r rune) bool {
		return r == ' ' ||
			r == '\t' ||
			r == '\n' ||
			r == ',' ||
			r == '?' ||
			r == ':' ||
			r == 0xff0c ||
			r == 0xff1f ||
			r == 0x3002 ||
			r == 0xff1a
	})

	seen := make(map[string]struct{})
	terms := make([]string, 0)
	add := func(term string) {
		term = strings.TrimSpace(term)
		if utf8.RuneCountInString(term) <= 1 {
			return
		}
		if _, ok := seen[term]; ok {
			return
		}
		seen[term] = struct{}{}
		terms = append(terms, term)
	}

	for _, term := range raw {
		add(term)
		runes := []rune(term)
		if len(runes) > 4 {
			for size := 4; size >= 2; size-- {
				for i := 0; i+size <= len(runes); i++ {
					add(string(runes[i : i+size]))
				}
			}
		}
	}
	return terms
}

func isMarkdownTableSeparator(line string) bool {
	if !strings.Contains(line, "|") {
		return false
	}
	cells := strings.Split(strings.Trim(line, "|"), "|")
	for _, cell := range cells {
		clean := strings.TrimSpace(cell)
		if clean == "" {
			continue
		}
		clean = strings.Trim(clean, ":")
		for _, r := range clean {
			if r != '-' {
				return false
			}
		}
	}
	return true
}

func mergeSource(left string, right string) string {
	if left == "" {
		return right
	}
	if right == "" || left == right || strings.Contains(left, right) {
		return left
	}
	return left + "+" + right
}

func anyToString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}

func anyToInt64(value any) int64 {
	switch typed := value.(type) {
	case string:
		item, _ := strconv.ParseInt(typed, 10, 64)
		return item
	case float64:
		return int64(typed)
	case json.Number:
		item, _ := typed.Int64()
		return item
	case int64:
		return typed
	case int:
		return int64(typed)
	default:
		return 0
	}
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
