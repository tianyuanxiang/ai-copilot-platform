package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiKbDomainModel = (*customAiKbDomainModel)(nil)

type (
	// AiKbDomainModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiKbDomainModel.
	AiKbDomainModel interface {
		aiKbDomainModel
		withSession(session sqlx.Session) AiKbDomainModel
		List(ctx context.Context, page, pageSize int64, keyword string, status int64, hasStatus bool) ([]AiKbDomain, int64, error)
	}

	customAiKbDomainModel struct {
		*defaultAiKbDomainModel
		db *gorm.DB
	}
)

// NewAiKbDomainModel returns a model for the database table.
func NewAiKbDomainModel(conn sqlx.SqlConn, db *gorm.DB) AiKbDomainModel {
	return &customAiKbDomainModel{
		defaultAiKbDomainModel: newAiKbDomainModel(conn),
		db:                     db,
	}
}

func (m *customAiKbDomainModel) withSession(session sqlx.Session) AiKbDomainModel {
	return NewAiKbDomainModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiKbDomainModel) List(ctx context.Context, page, pageSize int64, keyword string, status int64, hasStatus bool) ([]AiKbDomain, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := m.db.WithContext(ctx).Table("ai_kb_domain")
	if kw := strings.TrimSpace(keyword); kw != "" {
		db = db.Where("name ILIKE ? OR code ILIKE ?", "%"+kw+"%", "%"+kw+"%")
	}
	if hasStatus {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var domains []AiKbDomain
	err := db.Order("sort ASC, id ASC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&domains).Error
	return domains, total, err
}
