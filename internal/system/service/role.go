package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) CreateRole(ctx context.Context, req *model.CreateRoleRequest) (role *model.Role, err error) {
	if req.Role.Name == "" {
		err = errors.New("name is required")
		return
	}
	if req.Role.Code == "" {
		err = errors.New("code is required")
		return
	}
	_, err = s.roleDao.FindOne(ctx, &model.RoleFilter{Code: db.Eq(req.Role.Code)})
	if err == nil {
		err = errors.New("role already exists")
		return
	} else if db.IsDbError(err) {
		return
	}
	role = &model.Role{
		Name:         req.Role.Name,
		Code:         req.Role.Code,
		IsSuperAdmin: req.Role.IsSuperAdmin,
		Description:  req.Role.Description,
		Status:       model.EnableStatusEnabled,
	}
	if err = s.roleDao.Create(ctx, role); err != nil {
		err = errors.Errorf("failed to create role: %v", err)
		return
	}
	return
}

func (s *Service) UpdateRole(ctx context.Context, req *model.UpdateRoleRequest) (role *model.Role, err error) {
	role, err = s.roleDao.FindOneByID(ctx, req.Role.Id)
	if err != nil {
		return
	}
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "name") {
		update["name"] = req.Role.Name
	}
	if lo.Contains(req.UpdateMask, "code") {
		update["code"] = req.Role.Code
	}
	if lo.Contains(req.UpdateMask, "is_super_admin") {
		update["is_super_admin"] = req.Role.IsSuperAdmin
	}
	if lo.Contains(req.UpdateMask, "description") {
		update["description"] = req.Role.Description
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.Role.Status
	}
	if lo.Contains(req.UpdateMask, "permissions") {
		var data []byte
		data, err = json.Marshal(req.Role.Permission)
		if err != nil {
			err = errors.Errorf("failed to marshal permissions: %v", err)
			return
		}
		update["permissions"] = data
	}
	if len(update) == 0 {
		return
	}
	if err = s.roleDao.UpdateOne(ctx, role, update); err != nil {
		err = errors.Errorf("failed to update role: %v", err)
		return
	}
	return
}

func (s *Service) UpdateRolePermission(ctx context.Context, req *model.UpdateRolePermissionRequest) (role *model.Role, err error) {
	role, err = s.roleDao.FindOneByID(ctx, req.ID)
	if err != nil {
		return
	}
	var permission model.RolePermission
	if err = json.Unmarshal(role.Permission, &permission); err != nil {
		return
	}
	isChanged := false
	if lo.Contains(req.UpdateMask, "home") {
		permission.Home = req.Home
		isChanged = true
	}
	if lo.Contains(req.UpdateMask, "menu_ids") {
		permission.MenuIds = req.MenuIds
		isChanged = true
	}
	if lo.Contains(req.UpdateMask, "buttons") {
		permission.Buttons = req.Buttons
		isChanged = true
	}
	if !isChanged {
		return
	}
	data, _ := json.Marshal(permission)
	if err = s.roleDao.UpdateOne(ctx, role, map[string]any{"permission": data}); err != nil {
		err = errors.Errorf("failed to update role: %v", err)
		return
	}
	return
}

func (s *Service) DeleteRoles(ctx context.Context, req *model.DeleteRolesRequest) (err error) {
	if len(req.IDs) == 0 {
		err = errors.New("ids is empty")
		return
	}
	_, err = s.roleDao.Delete(ctx, &model.RoleFilter{IDs: db.In(req.IDs)})
	return err
}

func (s *Service) GetRoleList(ctx context.Context, req *model.GetRoleListRequest) (total int64, list []*model.Role, err error) {
	filter := &model.RoleFilter{
		Name:        db.Eq(req.Name, db.OmitIfZero[string]()),
		Code:        db.Eq(req.Code, db.OmitIfZero[string]()),
		Description: db.Eq(req.Description, db.OmitIfZero[string]()),
		Status:      db.Eq(req.Status, db.OmitIfZero[model.EnableStatus]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.roleDao.Count(ctx, filter)
		if err != nil {
			return
		}
		if !db.HasRecrods(total, req.PageParam.Page, req.PageParam.PageSize) {
			return
		}
		options = append(options,
			db.WithPaging(req.PageParam.Page, req.PageParam.PageSize),
			db.WithOrder(req.PageParam.OrderBy, nil))
	}
	list, err = s.roleDao.Find(ctx, filter, options...)
	return
}

func (s *Service) GetRole(ctx context.Context, req *model.GetRoleRequest) (role *model.Role, err error) {
	f := &model.RoleFilter{
		ID: db.Eq(req.ID),
	}
	role, err = s.roleDao.FindOne(ctx, f)
	return
}
