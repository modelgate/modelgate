package service

import (
	"context"

	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/dao"
)

type Service struct {
	userDao          system.UserDAO
	refreshTokenDao  system.RefreshTokenDAO
	roleDao          system.RoleDAO
	menuDao          system.MenuDAO
	permissionDao    system.PermissionDAO
	dataMigrationDao system.DataMigrationDAO
}

func New(i do.Injector) (system.Service, error) {
	return &Service{
		userDao:          do.MustInvoke[system.UserDAO](i),
		refreshTokenDao:  do.MustInvoke[system.RefreshTokenDAO](i),
		roleDao:          do.MustInvoke[system.RoleDAO](i),
		menuDao:          do.MustInvoke[system.MenuDAO](i),
		permissionDao:    do.MustInvoke[system.PermissionDAO](i),
		dataMigrationDao: do.MustInvoke[system.DataMigrationDAO](i),
	}, nil
}

func Init(i do.Injector) {
	dao.Init(i)
	do.Provide(i, New)

	s := do.MustInvoke[system.Service](i)
	if err := s.DataMigrate(context.Background()); err != nil {
		log.Errorf("data migrate failed: %v", err)
	}
}
