package model

import (
	"context"
	"strings"

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
		ListByKbID(ctx context.Context, kbID int64, page, pageSize int64, keyword, role string, hasRole bool) ([]AiKbMember, int64, error)
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

func (m *customAiKbMemberModel) ListByKbID(ctx context.Context, kbID int64, page, pageSize int64, keyword, role string, hasRole bool) ([]AiKbMember, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := m.db.WithContext(ctx).Table("ai_kb_member").Where("kb_id = ?", kbID)
	if kw := strings.TrimSpace(keyword); kw != "" {
		db = db.Where("user_id IN (SELECT id FROM sys_user WHERE username ILIKE ? OR nickname ILIKE ?)", "%"+kw+"%", "%"+kw+"%")
	}
	if hasRole && strings.TrimSpace(role) != "" {
		db = db.Where("role = ?", strings.TrimSpace(role))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var members []AiKbMember
	err := db.Order("id ASC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&members).Error
	return members, total, err
}
