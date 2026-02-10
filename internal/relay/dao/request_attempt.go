package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type RequestAttemptDao struct {
	*db.BaseDAO[model.RequestAttempt, model.RequestAttemptFilter]
}

func NewRequestAttemptDao(i do.Injector) (relay.RequestAttemptDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &RequestAttemptDao{
		BaseDAO: db.NewBaseDAO[model.RequestAttempt, model.RequestAttemptFilter](dbConn),
	}, nil
}
