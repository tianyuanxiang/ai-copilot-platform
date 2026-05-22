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

type AiDocumentChunkModel interface {
	InsertTrans(ctx context.Context, tx *gorm.DB, data *AiDocumentChunk) (int64, error)
	ListByDocumentID(ctx context.Context, documentID int64) ([]AiDocumentChunk, error)
	UpdateEmbeddingTrans(ctx context.Context, tx *gorm.DB, id int64, embedding PgVector) error
	DeleteByDocumentIDTrans(ctx context.Context, tx *gorm.DB, documentID int64) error
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
