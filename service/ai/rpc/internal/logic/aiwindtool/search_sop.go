// Package aiwindtool 的 search_sop 文件实现维护 SOP 知识库检索工具。
package aiwindtool

import (
	"context"
	"fmt"
	"strings"

	aiknowledgeservicelogic "ai-copilot-platform/ai-rpc/internal/logic/aiknowledgeservice"
	"ai-copilot-platform/ai-rpc/pb"
)

// executeSearchMaintenanceSOP 将 Agent 的 SOP 查询转换为现有 SearchKnowledge 请求。
// SearchKnowledge 内部已经包含知识库读取权限、向量召回、BM25 和 RRF 融合，
// 工具层不重复实现检索逻辑。
func (e *Executor) executeSearchMaintenanceSOP(ctx context.Context, req *pb.WindToolExecuteReq) (*toolResult, error) {
	var args SearchMaintenanceSOPArgs
	if err := decodeArgs(req.ArgumentsJson, &args); err != nil {
		return &toolResult{Status: statusInvalidArguments}, err
	}
	args.Query = strings.TrimSpace(args.Query)
	if args.Query == "" {
		return &toolResult{Status: statusInvalidArguments}, fmt.Errorf("query 不能为空")
	}
	args.TopK = normalizeTopK(args.TopK)
	if strings.TrimSpace(args.SearchScope) == "" {
		args.SearchScope = "public"
	}

	logic := aiknowledgeservicelogic.NewSearchKnowledgeLogic(ctx, e.svcCtx)
	resp, err := logic.SearchKnowledge(&pb.SearchKnowledgeReq{
		UserId:      req.UserId,
		KbId:        args.KBID,
		HasKbId:     args.KBID > 0,
		Query:       args.Query,
		TopK:        args.TopK,
		AnswerMode:  args.AnswerMode,
		SearchScope: args.SearchScope,
		DomainId:    args.DomainID,
		HasDomainId: args.DomainID > 0,
		DocumentIds: args.DocumentIDs,
	})
	if err != nil {
		return nil, err
	}

	citations := make([]*pb.Citation, 0, len(resp.Chunks))
	for _, chunk := range resp.Chunks {
		citations = append(citations, &pb.Citation{
			DocumentId: chunk.DocumentId,
			ChunkId:    chunk.ChunkId,
			Title:      chunk.Title,
			Snippet:    chunk.Snippet,
			Score:      chunk.Score,
		})
	}

	return &toolResult{
		ResultJSON: marshalJSON(map[string]any{
			"mode":    resp.Mode,
			"message": resp.Message,
			"chunks":  resp.Chunks,
		}),
		Citations: citations,
		Message:   resp.Message,
	}, nil
}
