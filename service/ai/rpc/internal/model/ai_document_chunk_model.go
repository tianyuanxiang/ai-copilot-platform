package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PgVector []float64

func (v PgVector) Literal() string {
	parts := make([]string, 0, len(v))
	for _, item := range v {
		parts = append(parts, strconv.FormatFloat(item, 'f', -1, 64))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

type AiDocumentChunk struct {
	Id            int64     `db:"id" gorm:"column:id"`
	DocumentId    int64     `db:"document_id" gorm:"column:document_id"`
	ParentChunkId int64     `db:"parent_chunk_id" gorm:"column:parent_chunk_id"`
	ChunkIndex    int64     `db:"chunk_index" gorm:"column:chunk_index"`
	Content       string    `db:"content" gorm:"column:content"`
	Embedding     PgVector  `db:"embedding" gorm:"-"`
	TokenCount    int64     `db:"token_count" gorm:"column:token_count"`
	CreatedAt     time.Time `db:"created_at" gorm:"column:created_at"`
}

type VectorChunkResult struct {
	ChunkId       int64   `gorm:"column:chunk_id"`
	DocumentId    int64   `gorm:"column:document_id"`
	ParentChunkId int64   `gorm:"column:parent_chunk_id"`
	Title         string  `gorm:"column:title"`
	Content       string  `gorm:"column:content"`
	Score         float64 `gorm:"column:score"`
}

type AiDocumentChunkModel interface {
	InsertTrans(ctx context.Context, tx *gorm.DB, data *AiDocumentChunk) (int64, error)
	ListByDocumentID(ctx context.Context, documentID int64) ([]AiDocumentChunk, error)
	UpdateEmbeddingTrans(ctx context.Context, tx *gorm.DB, id int64, embedding PgVector) error
	DeleteByDocumentIDTrans(ctx context.Context, tx *gorm.DB, documentID int64) error
	SearchDocumentChunksByVector(ctx context.Context, db *gorm.DB, queryVector PgVector, kbIDs []int64, documentIDs []int64, limit int) ([]VectorChunkResult, error)
}

type customAiDocumentChunkModel struct {
	db *gorm.DB
}

func NewAiDocumentChunkModel(db *gorm.DB) AiDocumentChunkModel {
	return &customAiDocumentChunkModel{db: db}
}

func (m *customAiDocumentChunkModel) InsertTrans(ctx context.Context, tx *gorm.DB, data *AiDocumentChunk) (int64, error) {
	if len(data.Embedding) == 0 {
		return 0, fmt.Errorf("embedding is empty")
	}

	data.CreatedAt = time.Now()

	var id int64
	err := tx.WithContext(ctx).Raw(`
		INSERT INTO ai_document_chunk
			(document_id, parent_chunk_id, chunk_index, content, embedding, token_count, created_at)
		VALUES
			(?, ?, ?, ?, ?::vector, ?, ?)
		RETURNING id
	`, data.DocumentId, data.ParentChunkId, data.ChunkIndex, data.Content, data.Embedding.Literal(), data.TokenCount, data.CreatedAt).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	data.Id = id
	return id, nil
}

func (m *customAiDocumentChunkModel) ListByDocumentID(ctx context.Context, documentID int64) ([]AiDocumentChunk, error) {
	var chunks []AiDocumentChunk
	err := m.db.WithContext(ctx).Table("ai_document_chunk").
		Where("document_id = ?", documentID).
		Order("parent_chunk_id ASC, chunk_index ASC, id ASC").
		Find(&chunks).Error
	return chunks, err
}

func (m *customAiDocumentChunkModel) UpdateEmbeddingTrans(ctx context.Context, tx *gorm.DB, id int64, embedding PgVector) error {
	if len(embedding) == 0 {
		return fmt.Errorf("embedding is empty")
	}
	result := tx.WithContext(ctx).Exec(
		`UPDATE ai_document_chunk SET embedding = ?::vector WHERE id = ?`,
		embedding.Literal(),
		id,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (m *customAiDocumentChunkModel) DeleteByDocumentIDTrans(ctx context.Context, tx *gorm.DB, documentID int64) error {
	return tx.WithContext(ctx).Table("ai_document_chunk").
		Where("document_id = ?", documentID).
		Delete(&AiDocumentChunk{}).Error
}

func (m *customAiDocumentChunkModel) SearchDocumentChunksByVector(ctx context.Context, db *gorm.DB, queryVector PgVector, kbIDs []int64, documentIDs []int64, limit int) ([]VectorChunkResult, error) {
	if len(queryVector) == 0 {
		return nil, fmt.Errorf("query vector is empty")
	}
	if len(kbIDs) == 0 {
		return []VectorChunkResult{}, nil
	}
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT
			c.id AS chunk_id,
			c.document_id,
			c.parent_chunk_id,
			d.file_name AS title,
			c.content,
			1 - (c.embedding <=> ?::vector) AS score
		FROM ai_document_chunk c
		JOIN ai_document d ON d.id = c.document_id
		JOIN ai_knowledge_base kb ON kb.id = d.kb_id
		WHERE d.kb_id IN ?
		  AND d.status = 'ready'
		  AND kb.status = 1
	`
	args := []interface{}{queryVector.Literal(), kbIDs}
	if len(documentIDs) > 0 {
		query += " AND d.id IN ?"
		args = append(args, documentIDs)
	}
	query += `
		ORDER BY c.embedding <=> ?::vector
		LIMIT ?
	`
	args = append(args, queryVector.Literal(), limit)

	var results []VectorChunkResult
	if err := db.WithContext(ctx).Raw(query, args...).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
