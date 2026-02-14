//go:generate go run go.uber.org/mock/mockgen -source $GOFILE -destination ./dao_mock.go -package $GOPACKAGE
package system

import (
	"context"

	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type UserDAO interface {
	Create(ctx context.Context, m *model.User) error
	Update(ctx context.Context, filter *model.UserFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.User, update map[string]any) error
	Count(ctx context.Context, filter *model.UserFilter) (int64, error)
	Find(ctx context.Context, filter *model.UserFilter, opts ...db.Option) ([]*model.User, error)
	FindOneByID(ctx context.Context, id int64) (*model.User, error)
	FindOne(ctx context.Context, filter *model.UserFilter, opts ...db.Option) (*model.User, error)
	Delete(ctx context.Context, filter *model.UserFilter) (int64, error)
}

type RefreshTokenDAO interface {
	Create(ctx context.Context, m *model.RefreshToken) error
	Update(ctx context.Context, filter *model.RefreshTokenFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.RefreshToken, update map[string]any) error
	Count(ctx context.Context, filter *model.RefreshTokenFilter) (int64, error)
	Find(ctx context.Context, filter *model.RefreshTokenFilter, opts ...db.Option) ([]*model.RefreshToken, error)
	FindOneByID(ctx context.Context, id int64) (*model.RefreshToken, error)
	FindOne(ctx context.Context, filter *model.RefreshTokenFilter, opts ...db.Option) (*model.RefreshToken, error)
	Delete(ctx context.Context, filter *model.RefreshTokenFilter) (int64, error)
}

type RoleDAO interface {
	Create(ctx context.Context, m *model.Role) error
	Update(ctx context.Context, filter *model.RoleFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Role, update map[string]any) error
	Count(ctx context.Context, filter *model.RoleFilter) (int64, error)
	Find(ctx context.Context, filter *model.RoleFilter, opts ...db.Option) ([]*model.Role, error)
	FindOneByID(ctx context.Context, id int64) (*model.Role, error)
	FindOne(ctx context.Context, filter *model.RoleFilter, opts ...db.Option) (*model.Role, error)
	Delete(ctx context.Context, filter *model.RoleFilter) (int64, error)
}

type MenuDAO interface {
	Create(ctx context.Context, m *model.Menu) error
	Update(ctx context.Context, filter *model.MenuFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Menu, update map[string]any) error
	Count(ctx context.Context, filter *model.MenuFilter) (int64, error)
	Find(ctx context.Context, filter *model.MenuFilter, opts ...db.Option) ([]*model.Menu, error)
	FindOneByID(ctx context.Context, id int64) (*model.Menu, error)
	FindOne(ctx context.Context, filter *model.MenuFilter, opts ...db.Option) (*model.Menu, error)
	Delete(ctx context.Context, filter *model.MenuFilter) (int64, error)
}

type PermissionDAO interface {
	Create(ctx context.Context, m *model.Permission) error
	Update(ctx context.Context, filter *model.PermissionFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Permission, update map[string]any) error
	Count(ctx context.Context, filter *model.PermissionFilter) (int64, error)
	Find(ctx context.Context, filter *model.PermissionFilter, opts ...db.Option) ([]*model.Permission, error)
	FindOneByID(ctx context.Context, id int64) (*model.Permission, error)
	FindOne(ctx context.Context, filter *model.PermissionFilter, opts ...db.Option) (*model.Permission, error)
	Delete(ctx context.Context, filter *model.PermissionFilter) (int64, error)
}

type DataMigrationDAO interface {
	AutoMigrate() error
	Create(ctx context.Context, m *model.DataMigration) error
	Update(ctx context.Context, filter *model.DataMigrationFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.DataMigration, update map[string]any) error
	Count(ctx context.Context, filter *model.DataMigrationFilter) (int64, error)
	Find(ctx context.Context, filter *model.DataMigrationFilter, opts ...db.Option) ([]*model.DataMigration, error)
	FindOneByID(ctx context.Context, id int64) (*model.DataMigration, error)
	FindOne(ctx context.Context, filter *model.DataMigrationFilter, opts ...db.Option) (*model.DataMigration, error)
	Delete(ctx context.Context, filter *model.DataMigrationFilter) (int64, error)
}
