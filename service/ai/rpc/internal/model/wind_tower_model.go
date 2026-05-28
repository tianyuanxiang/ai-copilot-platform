package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ WindTowerModel = (*customWindTowerModel)(nil)

type (
	// WindTowerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWindTowerModel.
	WindTowerModel interface {
		windTowerModel
		ListTurbines(ctx context.Context, farmCode string, keyword string) ([]*WindTower, error)
		withSession(session sqlx.Session) WindTowerModel
	}

	customWindTowerModel struct {
		*defaultWindTowerModel
		db *gorm.DB
	}
)

// NewWindTowerModel returns a model for the database table.
func NewWindTowerModel(conn sqlx.SqlConn, db *gorm.DB) WindTowerModel {
	return &customWindTowerModel{
		defaultWindTowerModel: newWindTowerModel(conn),
		db:                    db,
	}
}

func (m *customWindTowerModel) withSession(session sqlx.Session) WindTowerModel {
	return NewWindTowerModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customWindTowerModel) ListTurbines(ctx context.Context, farmCode string, keyword string) ([]*WindTower, error) {
	var rows []*WindTower
	query := m.db.WithContext(ctx).Table("wind_tower").Where("is_delete = 0")
	if farmCode = strings.TrimSpace(farmCode); farmCode != "" {
		query = query.Where("farm_code = ?", strings.ToUpper(farmCode))
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("tower_code like ? or farm_name like ?", like, like)
	}
	err := query.Order("farm_code asc, tower_code asc").Find(&rows).Error
	return rows, err
}
