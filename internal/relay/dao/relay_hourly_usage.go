package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

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
