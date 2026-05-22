package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiDocumentParentChunkModel = (*customAiDocumentParentChunkModel)(nil)

type (
	// AiDocumentParentChunkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiDocumentParentChunkModel.
	AiDocumentParentChunkModel interface {
		aiDocumentParentChunkModel
		withSession(session sqlx.Session) AiDocumentParentChunkModel
		InsertTrans(ctx context.Context, tx *gorm.DB, data *AiDocumentParentChunk) (int64, error)
		DeleteByDocumentIDTrans(ctx context.Context, tx *gorm.DB, documentID int64) error
	}

	customAiDocumentParentChunkModel struct {
		*defaultAiDocumentParentChunkModel
		db *gorm.DB
	}
)

// NewAiDocumentParentChunkModel returns a model for the database table.
func NewAiDocumentParentChunkModel(conn sqlx.SqlConn, db *gorm.DB) AiDocumentParentChunkModel {
	return &customAiDocumentParentChunkModel{
		defaultAiDocumentParentChunkModel: newAiDocumentParentChunkModel(conn),
		db:                                db,
	}
}

func (m *customAiDocumentParentChunkModel) withSession(session sqlx.Session) AiDocumentParentChunkModel {
	return NewAiDocumentParentChunkModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiDocumentParentChunkModel) InsertTrans(ctx context.Context, tx *gorm.DB, data *AiDocumentParentChunk) (int64, error) {
	data.CreatedAt = time.Now()
	result := tx.WithContext(ctx).Table("ai_document_parent_chunk").Create(data)
	return data.Id, result.Error
}

func (m *customAiDocumentParentChunkModel) DeleteByDocumentIDTrans(ctx context.Context, tx *gorm.DB, documentID int64) error {
	return tx.WithContext(ctx).Table("ai_document_parent_chunk").
		Where("document_id = ?", documentID).
		Delete(&AiDocumentParentChunk{}).Error
}
