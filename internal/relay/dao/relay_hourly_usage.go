package dao

import (
	"context"

	"github.com/samber/do/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type RelayHourlyUsageDao struct {
	*db.BaseDAO[model.RelayHourlyUsage, model.RelayHourlyUsageFilter]
}

func NewRelayHourlyUsageDao(i do.Injector) (relay.RelayHourlyUsageDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &RelayHourlyUsageDao{
		BaseDAO: db.NewBaseDAO[model.RelayHourlyUsage, model.RelayHourlyUsageFilter](dbConn),
	}, nil
}

func (d *RelayHourlyUsageDao) Create(ctx context.Context, m *model.RelayHourlyUsage) (err error) {
	err = d.GetDB().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "time"}, {Name: "provider_code"}},
		DoUpdates: clause.Assignments(map[string]any{
			"total_request": gorm.Expr("total_request + ?", m.TotalRequest),
			"total_success": gorm.Expr("total_success + ?", m.TotalSuccess),
			"total_failed":  gorm.Expr("total_failed + ?", m.TotalFailed),
			"total_point":   gorm.Expr("total_point + ?", m.TotalPoint),
		}),
	}).Create(m).Error
	return
}
