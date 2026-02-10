package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type PermissionDao struct {
	*db.BaseDAO[model.Permission, model.PermissionFilter]
}

func NewPermissionDao(i do.Injector) (system.PermissionDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &PermissionDao{
		BaseDAO: db.NewBaseDAO[model.Permission, model.PermissionFilter](dbConn),
	}, nil
}
