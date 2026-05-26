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
	query := m.db.WithContext(ctx).Table("wind_turbine").Where("is_delete = 0")
	if farmCode = strings.TrimSpace(farmCode); farmCode != "" {
		query = query.Where("farm_code = ?", strings.ToUpper(farmCode))
	}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("tower_code like ? or farm_name like ?", like, like)
	}
	err := query.Select("tower_id as turbine_id, tower_id, tower_code, tower_code as tower_name, farm_code, farm_name, risk_level, ai_enabled").
		Order("farm_code asc, tower_code asc").
		Find(&rows).Error
	return rows, err
}
