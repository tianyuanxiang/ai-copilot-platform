package model

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ WindFarmModel = (*customWindFarmModel)(nil)

type (
	// WindFarmModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWindFarmModel.
	WindFarmModel interface {
		windFarmModel
		ListFarms(ctx context.Context, keyword string) ([]*WindFarm, error)
		FarmDatabase(ctx context.Context, farmCode string) string
		withSession(session sqlx.Session) WindFarmModel
	}

	customWindFarmModel struct {
		*defaultWindFarmModel
		db *gorm.DB
	}
)

// NewWindFarmModel returns a model for the database table.
func NewWindFarmModel(conn sqlx.SqlConn, db *gorm.DB) WindFarmModel {
	return &customWindFarmModel{
		defaultWindFarmModel: newWindFarmModel(conn),
		db:                   db,
	}
}

func (m *customWindFarmModel) withSession(session sqlx.Session) WindFarmModel {
	return NewWindFarmModel(sqlx.NewSqlConnFromSession(session), m.db)
}

func (m *customWindFarmModel) ListFarms(ctx context.Context, keyword string) ([]*WindFarm, error) {
	var rows []*WindFarm
	query := m.db.WithContext(ctx).Table("wind_farm").Where("is_delete = 0")
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("farm_code like ? or farm_name like ?", like, like)
	}
	err := query.Order("farm_id asc").Find(&rows).Error
	return rows, err
}

func (m *customWindFarmModel) FarmDatabase(ctx context.Context, farmCode string) string {
	code := strings.ToUpper(strings.TrimSpace(farmCode))
	if m.db != nil && code != "" {
		var row struct {
			TDDatabase string `gorm:"column:td_database"`
		}
		err := m.db.WithContext(ctx).Raw(`select td_database from wind_farm where farm_code = ? and is_delete = 0 limit 1`, code).Scan(&row).Error
		if err == nil && row.TDDatabase != "" {
			return row.TDDatabase
		}
	}
	if v := farmDatabaseFallback[code]; v != "" {
		return v
	}
	return strings.ToLower(code)
}
