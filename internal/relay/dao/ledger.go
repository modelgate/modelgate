package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type LedgerDao struct {
	*db.BaseDAO[model.Ledger, model.LedgerFilter]
}

func NewLedgerDao(i do.Injector) (relay.LedgerDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &LedgerDao{
		BaseDAO: db.NewBaseDAO[model.Ledger, model.LedgerFilter](dbConn),
	}, nil
}
