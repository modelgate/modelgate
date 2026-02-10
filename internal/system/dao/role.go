package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type RoleDao struct {
	*db.BaseDAO[model.Role, model.RoleFilter]
}

func NewRoleDao(i do.Injector) (system.RoleDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &RoleDao{
		BaseDAO: db.NewBaseDAO[model.Role, model.RoleFilter](dbConn),
	}, nil
}
