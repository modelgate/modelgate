package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type ProviderApiKeyDao struct {
	*db.BaseDAO[model.ProviderApiKey, model.ProviderApiKeyFilter]
}

func NewProviderApiKeyDao(i do.Injector) (relay.ProviderApiKeyDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &ProviderApiKeyDao{
		BaseDAO: db.NewBaseDAO[model.ProviderApiKey, model.ProviderApiKeyFilter](dbConn),
	}, nil
}
