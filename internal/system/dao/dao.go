package dao

import (
	"github.com/samber/do/v2"
)

// Init 注册Dao
func Init(i do.Injector) {
	do.Provide(i, NewUserDao)
	do.Provide(i, NewRefreshTokenDao)
	do.Provide(i, NewRoleDao)
	do.Provide(i, NewMenuDao)
	do.Provide(i, NewPermissionDao)
	do.Provide(i, NewDataMigrationDao)
}
