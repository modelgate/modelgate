package v1

import (
	"context"
	"encoding/json"

	"connectrpc.com/connect/v2"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/system"
	"github.com/modelgate/modelgate/internal/system/model"
	v1pb "github.com/modelgate/modelgate/pkg/proto/admin/v1"
	systempb "github.com/modelgate/modelgate/pkg/proto/model/system"
	"github.com/modelgate/modelgate/pkg/types"
	"github.com/modelgate/modelgate/pkg/utils"
)

type SystemService struct {
	v1pb.UnimplementedSystemServiceHandler
	systemService system.Service
}

func NewSystemService(i do.Injector) (*SystemService, error) {
	return &SystemService{
		systemService: do.MustInvoke[system.Service](i),
	}, nil
}

func (s *SystemService) CreateUser(ctx context.Context, req *v1pb.CreateUserRequest) (userInfo *systempb.User, err error) {
	user, err := s.systemService.CreateUser(ctx, &model.CreateUserRequest{
		User: req.User,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	userInfo = user.ToProto()
	return
}

func (s *SystemService) UpdateUser(ctx context.Context, req *v1pb.UpdateUserRequest) (userInfo *systempb.User, err error) {
	user, err := s.systemService.UpdateUser(ctx, &model.UpdateUserRequest{
		User:       req.User,
		UpdateMask: req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	userInfo = user.ToProto()
	return
}

func (s *SystemService) GetUserList(ctx context.Context, req *v1pb.GetUserListRequest) (resp *v1pb.GetUserListResponse, err error) {
	total, users, err := s.systemService.GetUserList(ctx, &model.GetUserListRequest{
		PageParam: types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		Username:  req.Username,
		Phone:     req.Phone,
		Nickname:  req.Nickname,
		Gender:    req.Gender,
		Email:     req.Email,
		Status:    model.EnableStatus(req.Status),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetUserListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(users, func(user *model.User, _ int) *systempb.User {
			return user.ToProto()
		}),
	}
	return resp, nil
}

func (s *SystemService) DeleteUsers(ctx context.Context, req *v1pb.DeleteUsersRequest) (_ *emptypb.Empty, err error) {
	err = s.systemService.DeleteUsers(ctx, &model.DeleteUsersRequest{IDs: req.Ids})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	return
}

func (s *SystemService) GetRoleList(ctx context.Context, req *v1pb.GetRoleListRequest) (resp *v1pb.GetRoleListResponse, err error) {
	total, roles, err := s.systemService.GetRoleList(ctx, &model.GetRoleListRequest{
		PageParam:   types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      model.EnableStatus(req.Status),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetRoleListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(roles, func(role *model.Role, _ int) *systempb.Role {
			return role.ToProto()
		}),
	}
	return resp, nil
}

func (s *SystemService) GetRoleInfo(ctx context.Context, req *v1pb.GetRoleInfoRequest) (resp *systempb.Role, err error) {
	role, err := s.systemService.GetRole(ctx, &model.GetRoleRequest{ID: req.Id})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}

	perm, _ := role.GetPermission()
	info := role.ToProto()
	info.Permission = &systempb.RolePermission{
		Home:    perm.Home,
		MenuIds: perm.MenuIds,
		Buttons: perm.Buttons,
	}
	resp = info
	return resp, nil
}

func (s *SystemService) CreateRole(ctx context.Context, req *v1pb.CreateRoleRequest) (resp *systempb.Role, err error) {
	role, err := s.systemService.CreateRole(ctx, &model.CreateRoleRequest{
		Role: req.Role})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = role.ToProto()
	return resp, nil
}

func (s *SystemService) UpdateRole(ctx context.Context, req *v1pb.UpdateRoleRequest) (resp *systempb.Role, err error) {
	role, err := s.systemService.UpdateRole(ctx, &model.UpdateRoleRequest{
		Role:       req.Role,
		UpdateMask: req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = role.ToProto()
	return resp, nil
}

func (s *SystemService) UpdateRolePermission(ctx context.Context, req *v1pb.UpdateRolePermissionRequest) (resp *systempb.Role, err error) {
	role, err := s.systemService.UpdateRolePermission(ctx, &model.UpdateRolePermissionRequest{
		ID:         req.Id,
		UpdateMask: req.UpdateMask.Paths,
		Home:       req.Home,
		MenuIds:    req.MenuIds,
		Buttons:    req.Buttons,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = role.ToProto()
	return resp, nil
}

func (s *SystemService) DeleteRoles(ctx context.Context, req *v1pb.DeleteRolesRequest) (_ *emptypb.Empty, err error) {
	err = s.systemService.DeleteRoles(ctx, &model.DeleteRolesRequest{IDs: req.Ids})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	return
}

func (s *SystemService) GetMenuList(ctx context.Context, req *v1pb.GetMenuListRequest) (resp *v1pb.GetMenuListResponse, err error) {
	total, menus, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{
		PageParam: types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		IsRoot:    true,
		Name:      req.Name,
		Status:    model.EnableStatus(req.Status),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}

	_, childMenus, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{IsChild: true})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	groupedMenus := lo.GroupBy(childMenus, func(menu *model.Menu) int64 {
		return menu.Pid
	})
	var getChildren func(pid int64) []*systempb.Menu
	getChildren = func(pid int64) []*systempb.Menu {
		children := groupedMenus[pid]
		if len(children) == 0 {
			return nil
		}
		return lo.Map(children, func(menu *model.Menu, _ int) *systempb.Menu {
			item := menu.ToProto()
			item.Children = getChildren(menu.ID)
			return item
		})
	}

	resp = &v1pb.GetMenuListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(menus, func(menu *model.Menu, _ int) *systempb.Menu {
			item := menu.ToProto()
			item.Children = getChildren(menu.ID)
			return item
		}),
	}
	return resp, nil
}

func (s *SystemService) CreateMenu(ctx context.Context, req *v1pb.CreateMenuRequest) (resp *systempb.Menu, err error) {
	menu, err := s.systemService.CreateMenu(ctx, &model.CreateMenuRequest{
		Menu: req.Menu})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = menu.ToProto()
	return resp, nil
}

func (s *SystemService) UpdateMenu(ctx context.Context, req *v1pb.UpdateMenuRequest) (resp *systempb.Menu, err error) {
	menu, err := s.systemService.UpdateMenu(ctx, &model.UpdateMenuRequest{
		Menu:       req.Menu,
		UpdateMask: req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = menu.ToProto()
	return resp, nil
}

func (s *SystemService) DeleteMenus(ctx context.Context, req *v1pb.DeleteMenusRequest) (resp *emptypb.Empty, err error) {
	err = s.systemService.DeleteMenus(ctx, &model.DeleteMenusRequest{IDs: req.Ids})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return
}

func (s *SystemService) GetConstantRoutes(ctx context.Context, req *v1pb.GetConstantRoutesRequest) (resp *v1pb.GetConstantRoutesResponse, err error) {
	_, menus, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{
		IsConstant: true,
		IsRoot:     true,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	_, childMenus, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{
		IsChild: true,
		Pids: lo.Map(menus, func(menu *model.Menu, _ int) int64 {
			return menu.ID
		}),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	groupedMenus := lo.GroupBy(childMenus, func(menu *model.Menu) int64 {
		return menu.Pid
	})

	var getChildren func(pid int64) []*systempb.MenuRoute
	getChildren = func(pid int64) []*systempb.MenuRoute {
		children := groupedMenus[pid]
		if len(children) == 0 {
			return nil
		}
		return lo.Map(children, func(menu *model.Menu, _ int) *systempb.MenuRoute {
			item := menu.ToMenuRouteProto()
			item.Children = getChildren(menu.ID)
			return item
		})
	}

	resp = &v1pb.GetConstantRoutesResponse{
		Routes: lo.Map(menus, func(menu *model.Menu, _ int) *systempb.MenuRoute {
			item := menu.ToMenuRouteProto()
			item.Children = getChildren(menu.ID)
			return item
		}),
	}
	return
}

func (s *SystemService) GetUserRoutes(ctx context.Context, req *v1pb.GetUserRoutesRequest) (resp *v1pb.GetUserRoutesResponse, err error) {
	userId, ok := system.UserIDFromContext(ctx)
	if !ok {
		err = connect.NewError(connect.CodeUnauthenticated, "invalid access token")
		return
	}

	isSuperAdmin, perms, err := s.systemService.GetUserPermissions(ctx, userId)
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	var home string
	var authMenuIds []int64
	for _, perm := range perms {
		if home == "" {
			home = perm.Home
		}
		authMenuIds = append(authMenuIds, perm.MenuIds...)
	}

	logrus.Infof("auth menu ids: %v", authMenuIds)
	_, menuList, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{
		Status:      model.EnableStatusEnabled,
		IsUserRoute: true,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	groupMenus := lo.GroupBy(menuList, func(menu *model.Menu) int64 { return menu.Pid })

	var getChildMenu func(pid int64) []*systempb.MenuRoute
	getChildMenu = func(pid int64) []*systempb.MenuRoute {
		children := groupMenus[pid]
		if len(children) == 0 {
			return nil
		}
		return lo.FlatMap(children, func(menu *model.Menu, _ int) []*systempb.MenuRoute {
			children := getChildMenu(menu.ID)
			if !isSuperAdmin && len(children) == 0 && !lo.Contains(authMenuIds, menu.ID) {
				return nil
			}
			item := menu.ToMenuRouteProto()
			item.Children = children
			return []*systempb.MenuRoute{item}
		})
	}

	resp = &v1pb.GetUserRoutesResponse{
		Home:   home,
		Routes: getChildMenu(0),
	}
	return
}

func (s *SystemService) GetPageList(ctx context.Context, req *v1pb.GetPageListRequest) (resp *v1pb.GetPageListResponse, err error) {
	total, menuList, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{
		Status:    model.EnableStatusEnabled,
		PageParam: types.NewPageParam(1, 1000, ""),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetPageListResponse{
		Total: uint32(total),
		Records: lo.Map(menuList, func(item *model.Menu, _ int) string {
			return item.RouteName
		}),
	}
	return
}

// GetButtonList 按钮列表
func (s *SystemService) GetButtonList(ctx context.Context, req *v1pb.GetButtonListRequest) (resp *v1pb.GetButtonListResponse, err error) {
	_, menuList, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{
		Status:    model.EnableStatusEnabled,
		HasButton: true,
		PageParam: types.NewPageParam(1, 1000, ""),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}

	resp = &v1pb.GetButtonListResponse{
		Total: uint32(len(menuList)),
		Records: lo.Map(menuList, func(item *model.Menu, i int) *systempb.ButtonNode {
			var buttonsList []*model.MenuButton
			_ = json.Unmarshal([]byte(item.Buttons), &buttonsList)
			id := (i + 1) * 10000

			return &systempb.ButtonNode{
				Id:    int32(id),
				Label: item.Name,
				Code:  item.Name,
				Children: lo.Map(buttonsList, func(button *model.MenuButton, j int) *systempb.ButtonNode {
					return &systempb.ButtonNode{
						Id:    int32(id + j + 1),
						Label: button.Desc,
						Code:  button.Code,
					}
				}),
			}
		}),
	}
	return
}

// GetMenuTree 菜单树
func (s *SystemService) GetMenuTree(ctx context.Context, req *v1pb.GetMenuTreeRequest) (resp *v1pb.GetMenuTreeResponse, err error) {
	_, menuList, err := s.systemService.GetMenuList(ctx, &model.GetMenuListRequest{
		PageParam: types.NewPageParam(1, 1000, ""),
		Status:    model.EnableStatusEnabled,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	groupMenus := lo.GroupBy(menuList, func(menu *model.Menu) int64 { return menu.Pid })

	var getChildMenu func(pid int64) []*systempb.MenuNode
	getChildMenu = func(pid int64) []*systempb.MenuNode {
		children := groupMenus[pid]
		if len(children) == 0 {
			return nil
		}
		var list []*systempb.MenuNode
		for _, item := range children {
			info := &systempb.MenuNode{
				Id:       item.ID,
				Pid:      item.Pid,
				Label:    item.Name,
				Children: getChildMenu(item.ID),
			}
			list = append(list, info)
		}
		return list
	}
	resp = &v1pb.GetMenuTreeResponse{
		Records: getChildMenu(0),
	}
	return
}

func (s *SystemService) GetApiList(ctx context.Context, req *v1pb.GetApiListRequest) (resp *v1pb.GetApiListResponse, err error) {
	fileDescriptors := []protoreflect.FileDescriptor{
		v1pb.File_admin_v1_system_proto,
		v1pb.File_admin_v1_relay_proto,
	}
	var apiList []string
	for _, fileDescriptor := range fileDescriptors {
		apiList = append(apiList, utils.GetProtoMethods(fileDescriptor)...)
	}
	resp = &v1pb.GetApiListResponse{
		Records: apiList,
	}
	return
}

func (s *SystemService) CreatePermission(ctx context.Context, req *v1pb.CreatePermissionRequest) (resp *systempb.Permission, err error) {
	permission, err := s.systemService.CreatePermission(ctx, &model.CreatePermissionRequest{
		Permission: req.Permission,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = permission.ToProto()
	return
}

func (s *SystemService) UpdatePermission(ctx context.Context, req *v1pb.UpdatePermissionRequest) (resp *systempb.Permission, err error) {
	permission, err := s.systemService.UpdatePermission(ctx, &model.UpdatePermissionRequest{
		Permission: req.Permission,
		UpdateMask: req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = permission.ToProto()
	return
}

func (s *SystemService) DeletePermissions(ctx context.Context, req *v1pb.DeletePermissionsRequest) (resp *emptypb.Empty, err error) {
	if err = s.systemService.DeletePermissions(ctx, &model.DeletePermissionsRequest{
		IDs: req.Ids,
	}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return
}

func (s *SystemService) GetPermissionList(ctx context.Context, req *v1pb.GetPermissionListRequest) (resp *v1pb.GetPermissionListResponse, err error) {
	total, list, err := s.systemService.GetPermissionList(ctx, &model.GetPermissionListRequest{
		PageParam: types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		Name:      req.Name,
		Code:      req.Code,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetPermissionListResponse{
		Total: uint32(total),
		Records: lo.Map(list, func(item *model.Permission, _ int) *systempb.Permission {
			var apiPerms []*model.ApiPerm
			_ = json.Unmarshal([]byte(item.Data), &apiPerms)

			var apiPermList []*systempb.ApiPerm
			for _, apiPerm := range apiPerms {
				apiPermList = append(apiPermList, apiPerm.ToProto())
			}
			info := item.ToProto()
			info.Data = apiPermList
			return info
		}),
	}
	return
}

func (s *SystemService) GetVersion(ctx context.Context, req *v1pb.GetVersionRequest) (resp *v1pb.GetVersionResponse, err error) {
	resp = &v1pb.GetVersionResponse{
		Version:   config.Version,
		BuildTime: config.BuildTime,
		GitCommit: config.GitCommit,
	}
	return
}
