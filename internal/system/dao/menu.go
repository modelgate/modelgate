package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type MenuDao struct {
	*db.BaseDAO[model.Menu, model.MenuFilter]
}

func NewMenuDao(i do.Injector) (system.MenuDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &MenuDao{
		BaseDAO: db.NewBaseDAO[model.Menu, model.MenuFilter](dbConn),
	}, nil
}
