//go:generate go run go.uber.org/mock/mockgen -source $GOFILE -destination ./service_mock.go -package $GOPACKAGE
package system

import (
	"context"

	"github.com/modelgate/modelgate/internal/system/model"
)

type Service interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error)
	Authenticate(ctx context.Context, tokenStr string) (int64, error)

	GetUser(ctx context.Context, req *model.GetUserRequest) (*model.User, error)
	GetUserList(ctx context.Context, req *model.GetUserListRequest) (int64, []*model.User, error)
	CreateUser(ctx context.Context, user *model.CreateUserRequest) (*model.User, error)
	UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUsers(ctx context.Context, req *model.DeleteUsersRequest) error
	GetUserPermissions(ctx context.Context, userId int64) (bool, []*model.RolePermission, error)
	GetUserApiPermissions(ctx context.Context, userId int64) (bool, []*model.ApiPerm, error)

	GetRole(ctx context.Context, req *model.GetRoleRequest) (*model.Role, error)
	GetRoleList(ctx context.Context, req *model.GetRoleListRequest) (int64, []*model.Role, error)
	CreateRole(ctx context.Context, role *model.CreateRoleRequest) (*model.Role, error)
	UpdateRole(ctx context.Context, req *model.UpdateRoleRequest) (*model.Role, error)
	UpdateRolePermission(ctx context.Context, req *model.UpdateRolePermissionRequest) (*model.Role, error)
	DeleteRoles(ctx context.Context, req *model.DeleteRolesRequest) error

	GetMenu(ctx context.Context, req *model.GetMenuRequest) (*model.Menu, error)
	GetMenuList(ctx context.Context, req *model.GetMenuListRequest) (int64, []*model.Menu, error)
	CreateMenu(ctx context.Context, menu *model.CreateMenuRequest) (*model.Menu, error)
	UpdateMenu(ctx context.Context, req *model.UpdateMenuRequest) (*model.Menu, error)
	DeleteMenus(ctx context.Context, req *model.DeleteMenusRequest) error

	GetPermissionList(ctx context.Context, req *model.GetPermissionListRequest) (int64, []*model.Permission, error)
	CreatePermission(ctx context.Context, permission *model.CreatePermissionRequest) (*model.Permission, error)
	UpdatePermission(ctx context.Context, req *model.UpdatePermissionRequest) (*model.Permission, error)
	DeletePermissions(ctx context.Context, req *model.DeletePermissionsRequest) error
}
