package model

import (
	"context"
	"database/sql"
	"strings"

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
		FindAccessibleKnowledgeBaseIDsByScope(ctx context.Context, userId int64, scope string, domainID int64, hasDomainID bool) ([]int64, error)
		CountDocumentsByKbID(ctx context.Context, kbID int64) (int64, error)
		ListByCondition(ctx context.Context, query KbListQuery) ([]AiKnowledgeBase, int64, error)
	}

	KbListQuery struct {
		KbType       string
		HasKbType    bool
		DomainId     int64
		HasDomainId  bool
		Keyword      string
		Visibility   string
		HasVisibility bool
		Status       int64
		HasStatus    bool
		Page         int64
		PageSize     int64
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
	return m.FindAccessibleKnowledgeBaseIDsByScope(ctx, userID, "all", 0, false)
}

func (m *customAiKnowledgeBaseModel) FindAccessibleKnowledgeBaseIDsByScope(ctx context.Context, userID int64, scope string, domainID int64, hasDomainID bool) ([]int64, error) {
	type row struct {
		ID int64 `gorm:"column:id"`
	}
	var rows []row
	query := `
		SELECT DISTINCT kb.id
		FROM ai_knowledge_base kb
		LEFT JOIN ai_kb_member m
		  ON m.kb_id = kb.id AND m.user_id = ?
		WHERE kb.status = 1
	`
	args := []interface{}{userID}
	switch scope {
	case "personal":
		query += " AND kb.kb_type = 'personal' AND kb.owner_user_id = ?"
		args = append(args, userID)
	case "public":
		query += `
		  AND kb.kb_type = 'public'
		  AND (kb.visibility = 'public' OR m.role IN ('viewer', 'editor', 'manager'))
		`
		if hasDomainID && domainID > 0 {
			query += " AND kb.domain_id = ?"
			args = append(args, domainID)
		}
	default:
		query += `
		  AND (
		    kb.owner_user_id = ?
		    OR kb.visibility = 'public'
		    OR m.role IN ('viewer', 'editor', 'manager')
		  )
		`
		args = append(args, userID)
		if hasDomainID && domainID > 0 {
			query += " AND (kb.kb_type <> 'public' OR kb.domain_id = ?)"
			args = append(args, domainID)
		}
	}
	err := m.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, item := range rows {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func (m *customAiKnowledgeBaseModel) CountDocumentsByKbID(ctx context.Context, kbID int64) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Table("ai_document").
		Where("kb_id = ?", kbID).
		Count(&count).Error
	return count, err
}

func (m *customAiKnowledgeBaseModel) ListByCondition(ctx context.Context, query KbListQuery) ([]AiKnowledgeBase, int64, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	db := m.db.WithContext(ctx).Table("ai_knowledge_base")

	if query.HasKbType && strings.TrimSpace(query.KbType) != "" {
		db = db.Where("kb_type = ?", strings.TrimSpace(query.KbType))
	}
	if query.HasDomainId && query.DomainId > 0 {
		db = db.Where("domain_id = ?", query.DomainId)
	}
	if kw := strings.TrimSpace(query.Keyword); kw != "" {
		db = db.Where("name ILIKE ?", "%"+kw+"%")
	}
	if query.HasVisibility && strings.TrimSpace(query.Visibility) != "" {
		db = db.Where("visibility = ?", strings.TrimSpace(query.Visibility))
	}
	if query.HasStatus {
		db = db.Where("status = ?", query.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var kbs []AiKnowledgeBase
	err := db.Order("id DESC").
		Limit(int(query.PageSize)).
		Offset(int((query.Page - 1) * query.PageSize)).
		Find(&kbs).Error
	return kbs, total, err
}

// GetDomainNameByID returns the domain name for a given domain ID.
func GetDomainNameByID(ctx context.Context, db *gorm.DB, domainID int64) string {
	if domainID <= 0 {
		return ""
	}
	var name string
	err := db.WithContext(ctx).Table("ai_kb_domain").
		Select("name").
		Where("id = ?", domainID).
		Scan(&name).Error
	if err != nil {
		return ""
	}
	return name
}

// NullInt64Value returns the int64 value of a sql.NullInt64, or 0 if not valid.
func NullInt64Value(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}
