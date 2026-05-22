package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiKbMemberModel = (*customAiKbMemberModel)(nil)

type (
	// AiKbMemberModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiKbMemberModel.
	AiKbMemberModel interface {
		aiKbMemberModel
		withSession(session sqlx.Session) AiKbMemberModel
		FindByKbIDUserID(ctx context.Context, kbID int64, userID int64) (*AiKbMember, error)
	}

	customAiKbMemberModel struct {
		*defaultAiKbMemberModel
		db *gorm.DB
	}
)

// NewAiKbMemberModel returns a model for the database table.
func NewAiKbMemberModel(conn sqlx.SqlConn, db *gorm.DB) AiKbMemberModel {
	return &customAiKbMemberModel{
		defaultAiKbMemberModel: newAiKbMemberModel(conn),
		db:                     db,
	}
}

func (m *customAiKbMemberModel) withSession(session sqlx.Session) AiKbMemberModel {
	return NewAiKbMemberModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiKbMemberModel) FindByKbIDUserID(ctx context.Context, kbID int64, userID int64) (*AiKbMember, error) {
	var member AiKbMember
	result := m.db.WithContext(ctx).Table("ai_kb_member").
		Where("kb_id = ? AND user_id = ?", kbID, userID).
		First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &member, nil
}
