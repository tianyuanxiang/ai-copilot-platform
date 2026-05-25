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

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	hybridDefaultTopK = 5
	hybridMaxTopK     = 20
	hybridRecallLimit = 20
	rrfK              = 60.0
)

type SearchKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
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
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "query 不能为空")
	}

	topK := normalizeTopK(in.TopK)
	accessibleKbIDs, err := l.resolveAccessibleKbIDs(in)
	if err != nil {
		return nil, err
	}
	if len(accessibleKbIDs) == 0 {
		return &pb.SearchKnowledgeResp{
			Chunks:  []*pb.ChunkItem{},
			Mode:    "hybrid",
			Message: "no accessible knowledge bases",
		}, nil
	}

	embedded, err := l.svcCtx.EngineCallClient.EngineEmbed(l.ctx, []string{query})
	if err != nil {
		return nil, err
	}
	if len(embedded.Vectors) != 1 || len(embedded.Vectors[0]) == 0 {
		return nil, fmt.Errorf("engine embed returned invalid query vector")
	}
	queryVector := model.PgVector(embedded.Vectors[0])

	vectorResults, err := l.searchPgvector(queryVector, accessibleKbIDs, in.DocumentIds, hybridRecallLimit)
	if err != nil {
		return nil, err
	}

	bm25Results, err := l.searchElasticsearch(query, accessibleKbIDs, in.DocumentIds, hybridRecallLimit)
	if err != nil {
		return nil, err
	}

	candidates := fuseByRRF(vectorResults, bm25Results, topK)
	chunks := make([]*pb.ChunkItem, 0, len(candidates))
	for _, item := range candidates {
		chunks = append(chunks, &pb.ChunkItem{
			DocumentId: item.DocumentId,
			ChunkId:    item.ChunkId,
			Title:      item.Title,
			Snippet:    makeSnippet(item.Content, 160),
			Content:    item.Content,
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
		Mode:    "hybrid",
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

func (l *SearchKnowledgeLogic) resolveAccessibleKbIDs(in *pb.SearchKnowledgeReq) ([]int64, error) {
	if in.HasKbId || in.KbId > 0 {
		kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
			}
			return nil, err
		}
		if !l.canReadKnowledgeBase(kb, in.UserId) {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有读取该知识库的权限")
		}
		return []int64{in.KbId}, nil
	}

	return l.svcCtx.AiKnowledgeBaseModel.FindAccessibleKnowledgeBaseIDs(l.ctx, in.UserId)
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
		limit = hybridRecallLimit
	}

	filters := []map[string]any{
		{"terms": map[string]any{"kb_id": int64SliceToStringSlice(kbIDs)}},
	}
	if len(documentIDs) > 0 {
		filters = append(filters, map[string]any{"terms": map[string]any{"document_id": int64SliceToStringSlice(documentIDs)}})
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

func fuseByRRF(vectorResults []retrievedCandidate, bm25Results []retrievedCandidate, topK int) []retrievedCandidate {
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
		if existing.Source != item.Source {
			existing.Source = "hybrid"
		}
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
	if topK > 0 && len(items) > topK {
		return items[:topK]
	}
	return items
}

func makeSnippet(content string, maxRunes int) string {
	text := strings.TrimSpace(content)
	runes := []rune(text)
	if maxRunes <= 0 || len(runes) <= maxRunes {
		return text
	}
	return string(runes[:maxRunes])
}

func int64SliceToStringSlice(items []int64) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, strconv.FormatInt(item, 10))
	}
	return result
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
