package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type ModelDao struct {
	*db.BaseDAO[model.Model, model.ModelFilter]
}

func NewModelDao(i do.Injector) (relay.ModelDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &ModelDao{
		BaseDAO: db.NewBaseDAO[model.Model, model.ModelFilter](dbConn),
	}, nil
}
