package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type ProviderDao struct {
	*db.BaseDAO[model.Provider, model.ProviderFilter]
}

func NewProviderDao(i do.Injector) (relay.ProviderDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &ProviderDao{
		BaseDAO: db.NewBaseDAO[model.Provider, model.ProviderFilter](dbConn),
	}, nil
}
