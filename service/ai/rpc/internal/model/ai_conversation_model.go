package model

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ AiConversationModel = (*customAiConversationModel)(nil)

type (
	// AiConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiConversationModel.
	AiConversationModel interface {
		aiConversationModel
		withSession(session sqlx.Session) AiConversationModel
		InsertReturningID(ctx context.Context, data *AiConversation) (int64, error)
		FindByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversation, error)
		FindContextByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversation, error)
		ListByUser(ctx context.Context, query ConversationListQuery) ([]ConversationListItem, int64, error)
		UpdateTitleByIDUserID(ctx context.Context, id int64, userID int64, title string) error
		SoftDeleteByIDUserID(ctx context.Context, id int64, userID int64, deletedBy int64) error
		UpdateSummaryByIDUserID(ctx context.Context, id int64, userID int64, summary string) error
		TouchUpdatedAt(ctx context.Context, id int64) error
	}

	customAiConversationModel struct {
		*defaultAiConversationModel
		db *gorm.DB
	}

	ConversationListQuery struct {
		Page     int64
		PageSize int64
		UserID   int64
		KbID     int64
		HasKbID  bool
		Keyword  string
	}

	ConversationListItem struct {
		Id            int64     `db:"id"`
		KbId          int64     `db:"kb_id"`
		HasKbId       bool      `db:"has_kb_id"`
		Title         string    `db:"title"`
		LatestMessage string    `db:"latest_message"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}
)

// NewAiConversationModel returns a model for the database table.
func NewAiConversationModel(conn sqlx.SqlConn, db *gorm.DB) AiConversationModel {
	return &customAiConversationModel{
		defaultAiConversationModel: newAiConversationModel(conn),
		db:                         db,
	}
}

func (m *customAiConversationModel) withSession(session sqlx.Session) AiConversationModel {
	return NewAiConversationModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customAiConversationModel) InsertReturningID(ctx context.Context, data *AiConversation) (int64, error) {
	type conversationRecord struct {
		ID                  int64         `gorm:"column:id;primaryKey"`
		UserID              int64         `gorm:"column:user_id"`
		KbID                sql.NullInt64 `gorm:"column:kb_id"`
		Title               string        `gorm:"column:title"`
		ConversationSummary string        `gorm:"column:conversation_summary"`
	}
	record := conversationRecord{
		UserID:              data.UserId,
		KbID:                data.KbId,
		Title:               data.Title,
		ConversationSummary: data.ConversationSummary,
	}
	result := m.db.WithContext(ctx).Table("ai_conversation").Create(&record)
	if result.Error != nil {
		return 0, result.Error
	}
	data.Id = record.ID
	return record.ID, nil
}

func (m *customAiConversationModel) FindByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversation, error) {
	var resp AiConversation
	result := m.db.WithContext(ctx).Table("ai_conversation").
		Where("id = ? and user_id = ? and deleted_at is null", id, userID).
		First(&resp)
	switch result.Error {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, result.Error
	}
}

func (m *customAiConversationModel) FindContextByIDUserID(ctx context.Context, id int64, userID int64) (*AiConversation, error) {
	var resp AiConversation
	result := m.db.WithContext(ctx).Table("ai_conversation").
		Where("id = ? and user_id = ? and deleted_at is null", id, userID).
		First(&resp)

	switch result.Error {
	case nil:
		return &resp, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, result.Error
	}
}

func (m *customAiConversationModel) ListByUser(ctx context.Context, query ConversationListQuery) ([]ConversationListItem, int64, error) {
	db := m.db.WithContext(ctx).Table("ai_conversation as c").
		Where("c.user_id = ? and c.deleted_at is null", query.UserID)
	if query.HasKbID && query.KbID > 0 {
		db = db.Where("c.kb_id = ?", query.KbID)
	}
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		db = db.Where("c.title ILIKE ?", "%"+keyword+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, pageSize := normalizePage(query.Page, query.PageSize)
	var list []ConversationListItem
	err := db.Select(`
		  c.id,
		  coalesce(c.kb_id, 0) as kb_id,
		  c.kb_id is not null as has_kb_id,
		  c.title,
		  coalesce((
			select m.content
			from "public"."ai_message" m
			where m.conversation_id = c.id and m.deleted_at is null
			order by m.created_at desc, m.id desc
			limit 1
		  ), '') as latest_message,
		  c.created_at,
		  c.updated_at`).
		Order("c.updated_at desc, c.id desc").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *customAiConversationModel) UpdateTitleByIDUserID(ctx context.Context, id int64, userID int64, title string) error {
	result := m.db.WithContext(ctx).Table("ai_conversation").
		Where("id = ? and user_id = ? and deleted_at is null", id, userID).
		Updates(map[string]interface{}{
			"title":      title,
			"updated_at": time.Now(),
		})
	return gormNotFoundIfNoRows(result)
}

func (m *customAiConversationModel) SoftDeleteByIDUserID(ctx context.Context, id int64, userID int64, deletedBy int64) error {
	now := time.Now()
	result := m.db.WithContext(ctx).Table("ai_conversation").
		Where("id = ? and user_id = ? and deleted_at is null", id, userID).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
		})
	return gormNotFoundIfNoRows(result)
}

func (m *customAiConversationModel) UpdateSummaryByIDUserID(ctx context.Context, id int64, userID int64, summary string) error {
	result := m.db.WithContext(ctx).Table("ai_conversation").
		Where("id = ? and user_id = ? and deleted_at is null", id, userID).
		Updates(map[string]interface{}{
			"conversation_summary": summary,
			"updated_at":           time.Now(),
		})
	return gormNotFoundIfNoRows(result)
}

func (m *customAiConversationModel) TouchUpdatedAt(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Table("ai_conversation").
		Where("id = ? and deleted_at is null", id).
		Update("updated_at", time.Now()).
		Error
}

func normalizePage(page int64, pageSize int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func gormNotFoundIfNoRows(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
