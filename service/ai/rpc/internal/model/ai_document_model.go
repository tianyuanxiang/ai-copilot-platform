package model

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiDocumentModel = (*customAiDocumentModel)(nil)

type (
	// AiDocumentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiDocumentModel.
	AiDocumentModel interface {
		aiDocumentModel
		withSession(session sqlx.Session) AiDocumentModel
		FindByKbIDContentHash(ctx context.Context, kbID int64, contentHash string) (*AiDocument, error)
		FindByIDKbID(ctx context.Context, id int64, kbID int64) (*AiDocument, error)
		ListByKbID(ctx context.Context, query AiDocumentListQuery) ([]AiDocument, int64, error)
		InsertTrans(ctx context.Context, tx *gorm.DB, data *AiDocument) (int64, error)
		UpdateTrans(ctx context.Context, tx *gorm.DB, id int64, updates map[string]interface{}) error
		UpdateStatus(ctx context.Context, id int64, status string, errorMsg string) error
		DeleteByIDKbIDTrans(ctx context.Context, tx *gorm.DB, id int64, kbID int64) error
	}

	AiDocumentListQuery struct {
		KbID     int64
		Page     int64
		PageSize int64
		Keyword  string
		Status   string
	}

	customAiDocumentModel struct {
		*defaultAiDocumentModel
		db *gorm.DB
	}
)

// NewAiDocumentModel returns a model for the database table.
func NewAiDocumentModel(conn sqlx.SqlConn, db *gorm.DB) AiDocumentModel {
	return &customAiDocumentModel{
		defaultAiDocumentModel: newAiDocumentModel(conn),
		db:                     db,
	}
}

func (m *customAiDocumentModel) withSession(session sqlx.Session) AiDocumentModel {
	return NewAiDocumentModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiDocumentModel) FindByKbIDContentHash(ctx context.Context, kbID int64, contentHash string) (*AiDocument, error) {
	var doc AiDocument
	result := m.db.WithContext(ctx).Table("ai_document").
		Where("kb_id = ? AND content_hash = ?", kbID, contentHash).
		First(&doc)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &doc, nil
}

func (m *customAiDocumentModel) FindByIDKbID(ctx context.Context, id int64, kbID int64) (*AiDocument, error) {
	var doc AiDocument
	result := m.db.WithContext(ctx).Table("ai_document").
		Where("id = ? AND kb_id = ?", id, kbID).
		First(&doc)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &doc, nil
}

func (m *customAiDocumentModel) ListByKbID(ctx context.Context, query AiDocumentListQuery) ([]AiDocument, int64, error) {
	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := m.db.WithContext(ctx).Table("ai_document").Where("kb_id = ?", query.KbID)
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		db = db.Where("file_name ILIKE ?", "%"+keyword+"%")
	}
	if status := strings.TrimSpace(query.Status); status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var docs []AiDocument
	err := db.Order("id DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&docs).Error
	return docs, total, err
}

func (m *customAiDocumentModel) InsertTrans(ctx context.Context, tx *gorm.DB, data *AiDocument) (int64, error) {
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	result := tx.WithContext(ctx).Table("ai_document").Create(data)
	return data.Id, result.Error
}

func (m *customAiDocumentModel) UpdateTrans(ctx context.Context, tx *gorm.DB, id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = time.Now()
	result := tx.WithContext(ctx).Table("ai_document").
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (m *customAiDocumentModel) UpdateStatus(ctx context.Context, id int64, status string, errorMsg string) error {
	return m.db.WithContext(ctx).Table("ai_document").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"error_msg":  errorMsg,
			"updated_at": time.Now(),
		}).Error
}

func (m *customAiDocumentModel) DeleteByIDKbIDTrans(ctx context.Context, tx *gorm.DB, id int64, kbID int64) error {
	result := tx.WithContext(ctx).Table("ai_document").
		Where("id = ? AND kb_id = ?", id, kbID).
		Delete(&AiDocument{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
