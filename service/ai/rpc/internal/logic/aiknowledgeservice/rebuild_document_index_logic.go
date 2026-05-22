package aiknowledgeservicelogic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RebuildDocumentIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRebuildDocumentIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RebuildDocumentIndexLogic {
	return &RebuildDocumentIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RebuildDocumentIndex re-embeds existing child chunks and rewrites Elasticsearch.
func (l *RebuildDocumentIndexLogic) RebuildDocumentIndex(in *pb.RebuildDocumentIndexReq) (*pb.TaskResp, error) {
	if in.UserId <= 0 || in.KbId <= 0 || in.DocumentId <= 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "user_id、kb_id、document_id 不能为空")
	}

	kb, err := l.svcCtx.AiKnowledgeBaseModel.FindByID(l.ctx, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "知识库不存在")
		}
		return nil, err
	}
	if !canMaintainKnowledgeBase(l.ctx, l.svcCtx, kb, in.UserId) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "没有维护该知识库的权限")
	}

	doc, err := l.svcCtx.AiDocumentModel.FindByIDKbID(l.ctx, in.DocumentId, in.KbId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "文档不存在")
		}
		return nil, err
	}
	if doc.Status == documentStatusParsing || doc.Status == documentStatusIndexing {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "文档正在处理，不能重建")
	}

	chunks, err := l.svcCtx.AiDocumentChunkModel.ListByDocumentID(l.ctx, doc.Id)
	if err != nil {
		return nil, err
	}
	if len(chunks) == 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "该文档未解析，请重新解析")
	}

	if err := l.svcCtx.AiDocumentModel.UpdateStatus(l.ctx, doc.Id, documentStatusIndexing, ""); err != nil {
		return nil, err
	}

	docCopy := *doc
	kbCopy := *kb
	chunksCopy := append([]model.AiDocumentChunk(nil), chunks...)
	go l.runRebuildDocumentIndexJob(docCopy, kbCopy, chunksCopy)

	return &pb.TaskResp{
		TaskId:  strconv.FormatInt(doc.Id, 10),
		Status:  documentStatusIndexing,
		Message: "document rebuild task accepted",
	}, nil
}

func (l *RebuildDocumentIndexLogic) runRebuildDocumentIndexJob(doc model.AiDocument, kb model.AiKnowledgeBase, chunks []model.AiDocumentChunk) {
	ctx, cancel := context.WithTimeout(context.Background(), backgroundIngestTimeout)
	defer cancel()

	ingestLogic := NewIngestDocumentLogic(ctx, l.svcCtx)

	rebuildChunks := make([]model.AiDocumentChunk, 0, len(chunks))
	texts := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		content := strings.TrimSpace(chunk.Content)
		if content == "" {
			continue
		}
		chunk.Content = content
		rebuildChunks = append(rebuildChunks, chunk)
		texts = append(texts, content)
	}
	if len(rebuildChunks) == 0 {
		ingestLogic.markDocumentFailed(doc.Id, fmt.Errorf("该文档没有可重建的子 chunk"))
		return
	}

	embedded, err := ingestLogic.callEmbed(ctx, texts)
	if err != nil {
		ingestLogic.markDocumentFailed(doc.Id, err)
		return
	}
	if len(embedded.Vectors) != len(rebuildChunks) {
		ingestLogic.markDocumentFailed(doc.Id, fmt.Errorf("embedding vector count mismatch: got %d, want %d", len(embedded.Vectors), len(rebuildChunks)))
		return
	}

	err = l.svcCtx.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, chunk := range rebuildChunks {
			if err := l.svcCtx.AiDocumentChunkModel.UpdateEmbeddingTrans(ctx, tx, chunk.Id, model.PgVector(embedded.Vectors[i])); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		ingestLogic.markDocumentFailed(doc.Id, err)
		return
	}

	if err := deleteDocumentFromElasticsearch(ctx, l.svcCtx, doc.Id); err != nil {
		ingestLogic.markDocumentFailed(doc.Id, err)
		return
	}

	indexItems := make([]flatChunk, 0, len(rebuildChunks))
	for _, chunk := range rebuildChunks {
		indexItems = append(indexItems, flatChunk{
			ParentDBID: chunk.ParentChunkId,
			ChunkDBID:  chunk.Id,
			Content:    chunk.Content,
		})
	}
	if err := ingestLogic.writeElasticsearch(ctx, &kb, ingestDocumentJob{
		DocumentID: doc.Id,
		UserID:     doc.UploadedBy,
		FileName:   doc.FileName,
		FileType:   doc.FileType,
		Knowledge:  kb,
	}, indexItems); err != nil {
		ingestLogic.markDocumentFailed(doc.Id, err)
		return
	}

	if err := l.svcCtx.AiDocumentModel.UpdateStatus(ctx, doc.Id, documentStatusReady, ""); err != nil {
		l.Logger.Errorf("mark rebuilt document ready error. document_id=%d err=%v", doc.Id, err)
	}
}
