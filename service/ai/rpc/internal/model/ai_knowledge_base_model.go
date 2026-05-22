package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiKnowledgeBaseModel = (*customAiKnowledgeBaseModel)(nil)

type (
	// AiKnowledgeBaseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiKnowledgeBaseModel.
	AiKnowledgeBaseModel interface {
		aiKnowledgeBaseModel
		withSession(session sqlx.Session) AiKnowledgeBaseModel
		FindByID(ctx context.Context, id int64) (*AiKnowledgeBase, error)
	}

	customAiKnowledgeBaseModel struct {
		*defaultAiKnowledgeBaseModel
		db *gorm.DB
	}
)

// NewAiKnowledgeBaseModel returns a model for the database table.
func NewAiKnowledgeBaseModel(conn sqlx.SqlConn, db *gorm.DB) AiKnowledgeBaseModel {
	return &customAiKnowledgeBaseModel{
		defaultAiKnowledgeBaseModel: newAiKnowledgeBaseModel(conn),
		db:                          db,
	}
}

func (m *customAiKnowledgeBaseModel) withSession(session sqlx.Session) AiKnowledgeBaseModel {
	return NewAiKnowledgeBaseModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiKnowledgeBaseModel) FindByID(ctx context.Context, id int64) (*AiKnowledgeBase, error) {
	var kb AiKnowledgeBase
	result := m.db.WithContext(ctx).Table("ai_knowledge_base").
		Where("id = ?", id).
		First(&kb)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &kb, nil
}
