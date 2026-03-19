package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type RequestDao struct {
	*db.BaseDAO[model.Request, model.RequestFilter]
}

func NewRequestDao(i do.Injector) (relay.RequestDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &RequestDao{
		BaseDAO: db.NewBaseDAO[model.Request, model.RequestFilter](dbConn),
	}, nil
}
