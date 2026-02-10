package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type UserDao struct {
	*db.BaseDAO[model.User, model.UserFilter]
}

func NewUserDao(i do.Injector) (system.UserDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &UserDao{
		BaseDAO: db.NewBaseDAO[model.User, model.UserFilter](dbConn),
	}, nil
}
