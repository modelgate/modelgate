package service

import (
	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/dao"
)

type Service struct {
	userDao         system.UserDAO
	refreshTokenDao system.RefreshTokenDAO
	roleDao         system.RoleDAO
	menuDao         system.MenuDAO
	permissionDao   system.PermissionDAO
}

func New(i do.Injector) (system.Service, error) {
	return &Service{
		userDao:         do.MustInvoke[system.UserDAO](i),
		refreshTokenDao: do.MustInvoke[system.RefreshTokenDAO](i),
		roleDao:         do.MustInvoke[system.RoleDAO](i),
		menuDao:         do.MustInvoke[system.MenuDAO](i),
		permissionDao:   do.MustInvoke[system.PermissionDAO](i),
	}, nil
}

func Init(i do.Injector) {
	dao.Init(i)
	do.Provide(i, New)
}
