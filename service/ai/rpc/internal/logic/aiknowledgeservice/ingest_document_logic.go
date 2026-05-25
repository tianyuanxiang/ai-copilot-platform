package aiknowledgeservicelogic

import (
	"ai-copilot-platform/ai-rpc/internal/engine"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ai-copilot-platform/ai-rpc/internal/model"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
	"go-zero-rpc/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

const backgroundIngestTimeout = 30 * time.Minute

type IngestDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type ingestDocumentJob struct {
	DocumentID   int64
	UserID       int64
	FileName     string
	FileType     string
	ParseContent string
	Knowledge    model.AiKnowledgeBase
}

func NewIngestDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IngestDocumentLogic {
	return &IngestDocumentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// IngestDocument creates a document task and starts parse/embed/index work in background.
func (l *IngestDocumentLogic) IngestDocument(in *pb.IngestDocumentReq) (*pb.TaskResp, error) {
	if err := validateIngestReq(in); err != nil {
		return nil, err
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

	fileType := normalizeFileType(in.FileType, in.FileName)
	parseContent, hashBytes, err := l.resolveDocumentContent(in, fileType)
	if err != nil {
		return nil, err
	}
	contentHash := sha256Hex(hashBytes)

	existing, err := l.svcCtx.AiDocumentModel.FindByKbIDContentHash(l.ctx, in.KbId, contentHash)
	if err != nil && err != model.ErrNotFound {
		l.Logger.Errorf("Query document is exist failed. %v", err)
		return nil, err
	}
	if existing != nil {
		switch existing.Status {
		case documentStatusParsing, documentStatusIndexing:
			return &pb.TaskResp{
				TaskId:  strconv.FormatInt(existing.Id, 10),
				Status:  existing.Status,
				Message: "document ingest is already running",
			}, nil
		case documentStatusFailed:
			return &pb.TaskResp{
				TaskId:  strconv.FormatInt(existing.Id, 10),
				Status:  documentStatusFailed,
				Message: "document ingest failed before; delete or rebuild it before retrying",
			}, nil
		default:
			return &pb.TaskResp{
				TaskId:  strconv.FormatInt(existing.Id, 10),
				Status:  "existing",
				Message: "document already exists in this knowledge base",
			}, nil
		}
	}

	documentID, err := l.createDocument(in, fileType, contentHash)
	if err != nil {
		l.Logger.Errorf("create document err: %v", err)
		return nil, err
	}

	go l.runIngestDocumentJob(ingestDocumentJob{
		DocumentID:   documentID,
		UserID:       in.UserId,
		FileName:     strings.TrimSpace(in.FileName),
		FileType:     fileType,
		ParseContent: parseContent,
		Knowledge:    *kb,
	})

	return &pb.TaskResp{
		TaskId:  strconv.FormatInt(documentID, 10),
		Status:  documentStatusParsing,
		Message: "document ingest task accepted",
	}, nil
}

func (l *IngestDocumentLogic) runIngestDocumentJob(job ingestDocumentJob) {
	ctx, cancel := context.WithTimeout(context.Background(), backgroundIngestTimeout)
	defer cancel()

	parsed, err := l.svcCtx.EngineCallClient.EngineParse(ctx, job.FileName, job.FileType, job.ParseContent)
	if err != nil {
		l.Logger.Errorf("callParse err: %v", err)
		l.markDocumentFailed(job.DocumentID, err)
		return
	}

	childTexts := collectChildTexts(parsed.Parents)
	if len(childTexts) == 0 {
		l.markDocumentFailed(job.DocumentID, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "文档解析后没有可入库的子 chunk"))
		return
	}

	embedded, err := l.svcCtx.EngineCallClient.EngineEmbed(ctx, childTexts)
	if err != nil {
		l.Logger.Errorf("callEmbed err: %v", err)
		l.markDocumentFailed(job.DocumentID, err)
		return
	}
	if len(embedded.Vectors) != len(childTexts) {
		err = fmt.Errorf("embedding vector count mismatch: got %d, want %d", len(embedded.Vectors), len(childTexts))
		l.markDocumentFailed(job.DocumentID, err)
		return
	}

	indexItems, err := l.writeChunks(ctx, job.DocumentID, parsed.Parents, embedded.Vectors)
	if err != nil {
		l.markDocumentFailed(job.DocumentID, err)
		return
	}

	if err := l.writeElasticsearch(ctx, &job.Knowledge, job, indexItems); err != nil {
		l.markDocumentFailed(job.DocumentID, err)
		return
	}

	if err := l.svcCtx.AiDocumentModel.UpdateStatus(ctx, job.DocumentID, documentStatusReady, ""); err != nil {
		l.Logger.Errorf("mark document ready error. document_id=%d err=%v", job.DocumentID, err)
	}
}

func validateIngestReq(in *pb.IngestDocumentReq) error {
	if in.UserId <= 0 || in.KbId <= 0 {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "user_id 和 kb_id 不能为空")
	}
	if strings.TrimSpace(in.FileName) == "" || strings.TrimSpace(in.FileType) == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "file_name 和 file_type 不能为空")
	}
	sourceType := strings.ToLower(strings.TrimSpace(in.SourceType))
	if sourceType == "" {
		if strings.TrimSpace(in.StoredPath) != "" {
			sourceType = sourceTypeUpload
		} else {
			sourceType = sourceTypeText
		}
	}
	switch sourceType {
	case sourceTypeText:
		if strings.TrimSpace(in.Content) == "" {
			return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "content 不能为空")
		}
	case sourceTypeUpload:
		if strings.TrimSpace(in.StoredPath) == "" {
			return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "stored_path 不能为空")
		}
	default:
		return xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "source_type 只支持 text 或 upload")
	}
	return nil
}

func (l *IngestDocumentLogic) resolveDocumentContent(in *pb.IngestDocumentReq, fileType string) (string, []byte, error) {
	sourceType := strings.ToLower(strings.TrimSpace(in.SourceType))
	if sourceType == "" {
		if strings.TrimSpace(in.StoredPath) != "" {
			sourceType = sourceTypeUpload
		} else {
			sourceType = sourceTypeText
		}
	}

	if sourceType == sourceTypeText {
		if !isTextFileType(fileType) {
			return "", nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "sourceType=text 仅支持 md/markdown/txt")
		}
		content := strings.TrimSpace(in.Content)
		return content, []byte(content), nil
	}

	data, err := l.readUploadedFile(in.StoredPath)
	if err != nil {
		return "", nil, err
	}

	switch fileType {
	case "md", "markdown", "txt":
		return string(data), data, nil
	case "docx":
		return base64.StdEncoding.EncodeToString(data), data, nil
	case "pdf", "xlsx", "xlsm", "xls":
		return "", nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "当前文档入库暂不支持 "+fileType)
	default:
		return "", nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "不支持的文件类型: "+fileType)
	}
}

func (l *IngestDocumentLogic) readUploadedFile(storedPath string) ([]byte, error) {
	uploadRoot := strings.TrimSpace(l.svcCtx.Config.Storage.UploadPath)
	if uploadRoot == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "Storage.UploadPath 未配置")
	}

	rootAbs, err := filepath.Abs(uploadRoot)
	if err != nil {
		return nil, err
	}
	cleanStored := filepath.Clean(strings.TrimLeft(strings.TrimSpace(storedPath), `/\`))
	if cleanStored == "." || strings.HasPrefix(cleanStored, "..") || filepath.IsAbs(cleanStored) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "stored_path 非法")
	}
	rootName := filepath.Base(rootAbs)
	if cleanStored == rootName {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "stored_path 非法")
	}
	if strings.HasPrefix(cleanStored, rootName+string(os.PathSeparator)) {
		cleanStored = strings.TrimPrefix(cleanStored, rootName+string(os.PathSeparator))
	}

	fullPath := filepath.Join(rootAbs, cleanStored)
	fullAbs, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, err
	}
	if fullAbs != rootAbs && !strings.HasPrefix(fullAbs, rootAbs+string(os.PathSeparator)) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "stored_path 越界")
	}

	data, err := os.ReadFile(fullAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "上传文件不存在")
		}
		return nil, err
	}
	return data, nil
}

func (l *IngestDocumentLogic) createDocument(in *pb.IngestDocumentReq, fileType string, contentHash string) (int64, error) {
	doc := &model.AiDocument{
		KbId:        in.KbId,
		UploadedBy:  in.UserId,
		FileName:    strings.TrimSpace(in.FileName),
		FileType:    fileType,
		ContentHash: contentHash,
		Status:      documentStatusParsing,
		ErrorMsg:    "",
	}

	var documentID int64
	err := l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		id, err := l.svcCtx.AiDocumentModel.InsertTrans(l.ctx, tx, doc)
		if err != nil {
			return err
		}
		documentID = id
		return nil
	})
	if err != nil {
		l.Logger.Errorf("Insert aiDocument failed. document_id=%d err=%v", documentID, err)
		return documentID, xerr.NewCodeErrorMsg(xerr.ErrInternal, "文档写入失败")
	}
	return documentID, err
}

func (l *IngestDocumentLogic) writeChunks(ctx context.Context, documentID int64, parents []engine.ParentChunk, vectors [][]float64) ([]flatChunk, error) {
	indexItems := make([]flatChunk, 0, len(vectors))
	vectorIndex := 0

	err := l.svcCtx.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := l.svcCtx.AiDocumentModel.UpdateTrans(ctx, tx, documentID, map[string]interface{}{"status": documentStatusIndexing}); err != nil {
			l.Logger.Errorf("UpdateTrans document status %s error: %v", documentStatusIndexing, err)
			return err
		}

		for _, parent := range parents {
			parentID, err := l.svcCtx.AiDocumentParentChunkModel.InsertTrans(ctx, tx, &model.AiDocumentParentChunk{
				DocumentId:  documentID,
				ParentIndex: parent.ParentIndex,
				Content:     strings.TrimSpace(parent.Content),
				TokenCount:  parent.TokenCount,
			})
			if err != nil {
				l.Logger.Errorf("AiDocumentParentChunk insert failed %v", err)
				return err
			}

			for _, child := range parent.Children {
				if strings.TrimSpace(child.Content) == "" {
					continue
				}
				if vectorIndex >= len(vectors) {
					return fmt.Errorf("embedding vector is missing for chunk %d", vectorIndex)
				}
				chunkID, err := l.svcCtx.AiDocumentChunkModel.InsertTrans(ctx, tx, &model.AiDocumentChunk{
					DocumentId:    documentID,
					ParentChunkId: parentID,
					ChunkIndex:    child.ChunkIndex,
					Content:       strings.TrimSpace(child.Content),
					Embedding:     model.PgVector(vectors[vectorIndex]),
					TokenCount:    child.TokenCount,
				})
				if err != nil {
					l.Logger.Errorf("AiDocumentChunk insert failed %v", err)
					return err
				}
				indexItems = append(indexItems, flatChunk{
					ParentDBID: parentID,
					ChunkDBID:  chunkID,
					Content:    strings.TrimSpace(child.Content),
				})
				vectorIndex++
			}
		}
		if vectorIndex != len(vectors) {
			return fmt.Errorf("unused embedding vectors: %d", len(vectors)-vectorIndex)
		}
		return nil
	})
	return indexItems, err
}

func (l *IngestDocumentLogic) writeElasticsearch(ctx context.Context, kb *model.AiKnowledgeBase, job ingestDocumentJob, items []flatChunk) error {
	indexName := strings.TrimSpace(l.svcCtx.Config.Elasticsearch.KbChunksIndex)
	if indexName == "" {
		return xerr.NewCodeErrorMsg(xerr.ErrInternal, "Elasticsearch.KbChunksIndex 必须配置")
	}
	if l.svcCtx.ES == nil {
		return xerr.NewCodeErrorMsg(xerr.ErrInternal, "Elasticsearch client 未初始化")
	}

	var body bytes.Buffer
	for _, item := range items {
		action := map[string]any{
			"index": map[string]any{
				"_index": indexName,
				"_id":    strconv.FormatInt(item.ChunkDBID, 10),
			},
		}
		doc := map[string]any{
			"kb_id":           strconv.FormatInt(kb.Id, 10),
			"kb_type":         kb.KbType,
			"owner_user_id":   nullInt64String(kb.OwnerUserId),
			"domain_id":       nullInt64String(kb.DomainId),
			"visibility":      kb.Visibility,
			"uploaded_by":     strconv.FormatInt(job.UserID, 10),
			"document_id":     strconv.FormatInt(job.DocumentID, 10),
			"parent_chunk_id": strconv.FormatInt(item.ParentDBID, 10),
			"chunk_id":        strconv.FormatInt(item.ChunkDBID, 10),
			"title":           job.FileName,
			"content":         item.Content,
			"file_name":       job.FileName,
			"file_type":       job.FileType,
			"created_at":      time.Now().Format(time.RFC3339),
		}
		actionBytes, _ := json.Marshal(action)
		docBytes, _ := json.Marshal(doc)
		body.Write(actionBytes)
		body.WriteByte('\n')
		body.Write(docBytes)
		body.WriteByte('\n')
	}

	resp, err := l.svcCtx.ES.Bulk(
		bytes.NewReader(body.Bytes()),
		l.svcCtx.ES.Bulk.WithContext(ctx),
	)
	if err != nil {
		l.Logger.Errorf("Bulk ES document failed %v", err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("elasticsearch bulk returned %s: %s", resp.Status(), string(respBody))
	}

	var result struct {
		Errors bool `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}
	if result.Errors {
		return fmt.Errorf("elasticsearch bulk response contains item errors: %s", string(respBody))
	}
	return nil
}

func (l *IngestDocumentLogic) markDocumentFailed(documentID int64, cause error) {
	if documentID <= 0 {
		return
	}
	msg := cause.Error()
	if len(msg) > 2000 {
		msg = msg[:2000]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := l.svcCtx.AiDocumentModel.UpdateStatus(ctx, documentID, documentStatusFailed, msg); err != nil {
		l.Logger.Errorf("mark document failed error. document_id=%d err=%v", documentID, err)
	}
}

func collectChildTexts(parents []engine.ParentChunk) []string {
	texts := make([]string, 0)
	for _, parent := range parents {
		for _, child := range parent.Children {
			content := strings.TrimSpace(child.Content)
			if content != "" {
				texts = append(texts, content)
			}
		}
	}
	return texts
}

func normalizeFileType(fileType string, fileName string) string {
	clean := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(fileType), "."))
	if clean != "" {
		return clean
	}
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))
}

func isTextFileType(fileType string) bool {
	return fileType == "md" || fileType == "markdown" || fileType == "txt"
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func nullInt64String(value sql.NullInt64) string {
	if value.Valid {
		return strconv.FormatInt(value.Int64, 10)
	}
	return ""
}
