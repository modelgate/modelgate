package dao

import (
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type RefreshTokenDao struct {
	*db.BaseDAO[model.RefreshToken, model.RefreshTokenFilter]
}

func NewRefreshTokenDao(i do.Injector) (system.RefreshTokenDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &RefreshTokenDao{
		BaseDAO: db.NewBaseDAO[model.RefreshToken, model.RefreshTokenFilter](dbConn),
	}, nil
}
