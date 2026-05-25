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
		FindByOwnerUserId(ctx context.Context, userId int64) ([]*AiKnowledgeBase, error)
		FindAllKbByUserId(ctx context.Context, userId int64) ([]*AiKnowledgeBase, error)
		FindAccessibleKnowledgeBaseIDs(ctx context.Context, userId int64) ([]int64, error)
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

func (m *customAiKnowledgeBaseModel) FindByOwnerUserId(ctx context.Context, userId int64) ([]*AiKnowledgeBase, error) {
	var kbs []*AiKnowledgeBase
	result := m.db.WithContext(ctx).Table("ai_knowledge_base").
		Where("owner_user_id = ?", userId).
		Find(&kbs)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return kbs, nil
}

func (m *customAiKnowledgeBaseModel) FindAllKbByUserId(ctx context.Context, userId int64) ([]*AiKnowledgeBase, error) {
	var kbs []*AiKnowledgeBase
	result := m.db.WithContext(ctx).Table("ai_knowledge_base AS kb").
		Select("DISTINCT kb.*").
		Joins("left join ai_kb_member AS m "+
			"on m.kb.id = kb.id and m.user_id = ?", userId).
		Where("kb.status = ?", 1).
		Where(" kb.owner_user_id = ? OR kb.visibility = ? OR m.role IN ? ",
			userId, "public", []string{"viewer", "editor", "manager"}).
		Find(&kbs)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return kbs, nil
}

func (m *customAiKnowledgeBaseModel) FindAccessibleKnowledgeBaseIDs(ctx context.Context, userID int64) ([]int64, error) {
	type row struct {
		ID int64 `gorm:"column:id"`
	}
	var rows []row
	err := m.db.WithContext(ctx).Raw(`
		SELECT DISTINCT kb.id
		FROM ai_knowledge_base kb
		LEFT JOIN ai_kb_member m
		  ON m.kb_id = kb.id AND m.user_id = ?
		WHERE kb.status = 1
		  AND (
		    kb.owner_user_id = ?
		    OR kb.visibility = 'public'
		    OR m.role IN ('viewer', 'editor', 'manager')
		  )
	`, userID, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, item := range rows {
		ids = append(ids, item.ID)
	}
	return ids, nil
}
