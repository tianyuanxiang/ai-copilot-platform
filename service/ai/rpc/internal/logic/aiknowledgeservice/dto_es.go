package aiknowledgeservicelogic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/svc"
	"go-zero-rpc/common/xerr"
)

func deleteDocumentFromElasticsearch(ctx context.Context, svcCtx *svc.ServiceContext, documentID int64) error {
	indexName := strings.TrimSpace(svcCtx.Config.Elasticsearch.KbChunksIndex)
	if indexName == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrInternal, "Elasticsearch.KbChunksIndex 必须配置")
	}
	if svcCtx.ES == nil {
		return xerr.NewCodeErrorMsg(xerr.ErrInternal, "Elasticsearch client 未初始化")
	}

	body, err := json.Marshal(map[string]any{
		"query": map[string]any{
			"term": map[string]any{
				"document_id": strconv.FormatInt(documentID, 10),
			},
		},
	})
	if err != nil {
		return err
	}

	resp, err := svcCtx.ES.DeleteByQuery(
		[]string{indexName},
		bytes.NewReader(body),
		svcCtx.ES.DeleteByQuery.WithContext(ctx),
		svcCtx.ES.DeleteByQuery.WithConflicts("proceed"),
		svcCtx.ES.DeleteByQuery.WithRefresh(true),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("elasticsearch delete_by_query returned %s: %s", resp.Status(), string(respBody))
	}
	return nil
}
