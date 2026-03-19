package dao

import (
	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

// Init 注册Dao
func Init(i do.Injector) {
	do.Provide(i, NewUserDao)
	do.Provide(i, NewRefreshTokenDao)
	do.Provide(i, NewRoleDao)
	do.Provide(i, NewMenuDao)
	do.Provide(i, NewPermissionDao)
	do.Provide(i, NewDataMigrationDao)

	dbConn := do.MustInvoke[*gorm.DB](i)
	dbConn.AutoMigrate(&model.DataMigration{})
}
