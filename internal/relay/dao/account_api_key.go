package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type AccountApiKeyDao struct {
	*db.BaseDAO[model.AccountApiKey, model.AccountApiKeyFilter]
}

func NewAccountApiKeyDao(i do.Injector) (relay.AccountApiKeyDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &AccountApiKeyDao{
		BaseDAO: db.NewBaseDAO[model.AccountApiKey, model.AccountApiKeyFilter](dbConn),
	}, nil
}
