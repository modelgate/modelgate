package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type RelayUsageDao struct {
	*db.BaseDAO[model.RelayUsage, model.RelayUsageFilter]
}

func NewRelayUsageDao(i do.Injector) (relay.RelayUsageDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &RelayUsageDao{
		BaseDAO: db.NewBaseDAO[model.RelayUsage, model.RelayUsageFilter](dbConn),
	}, nil
}
