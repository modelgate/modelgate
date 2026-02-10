package service

import (
	"testing"

	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/system"
)

func TestMain(m *testing.M) {

	i := do.New()
	config.Init(i, "../../../../configs/config.toml")

	do.Provide(nil, func(i do.Injector) (*Service, error) {
		return &Service{
			userDao:         do.MustInvoke[system.UserDAO](nil),
			refreshTokenDao: do.MustInvoke[system.RefreshTokenDAO](nil),
			roleDao:         do.MustInvoke[system.RoleDAO](nil),
			menuDao:         do.MustInvoke[system.MenuDAO](nil),
		}, nil
	})

	do.Provide(nil, func(i do.Injector) (system.UserDAO, error) {
		return nil, nil
	})

	do.Provide(nil, func(i do.Injector) (system.RefreshTokenDAO, error) {
		return nil, nil
	})

	do.Provide(nil, func(i do.Injector) (system.RoleDAO, error) {
		return nil, nil
	})

	do.Provide(nil, func(i do.Injector) (system.MenuDAO, error) {
		return nil, nil
	})

	_ = m.Run()
}
