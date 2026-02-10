package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
	modelpb "github.com/modelgate/modelgate/pkg/proto/model/system"
)

func (s *Service) CreatePermission(ctx context.Context, req *model.CreatePermissionRequest) (permission *model.Permission, err error) {
	apiPerms := lo.Map(req.Permission.Data, func(item *modelpb.ApiPerm, index int) *model.ApiPerm {
		return &model.ApiPerm{
			Path:   item.Path,
			Method: item.Method,
		}
	})
	data, err := json.Marshal(apiPerms)
	if err != nil {
		return
	}
	permission = &model.Permission{
		Name: req.Permission.Name,
		Code: req.Permission.Code,
		Data: data,
		Desc: req.Permission.Desc,
	}
	if err = s.permissionDao.Create(ctx, permission); err != nil {
		err = errors.Errorf("failed to create permission: %v", err)
		return
	}
	return
}

func (s *Service) UpdatePermission(ctx context.Context, req *model.UpdatePermissionRequest) (permission *model.Permission, err error) {
	permission, err = s.permissionDao.FindOneByID(ctx, req.Permission.Id)
	if err != nil {
		return
	}
	update := map[string]any{}
	if lo.Contains(req.UpdateMask, "name") {
		update["name"] = req.Permission.Name
	}
	if lo.Contains(req.UpdateMask, "code") {
		update["code"] = req.Permission.Code
	}
	if lo.Contains(req.UpdateMask, "data") {
		apiPerms := lo.Map(req.Permission.Data, func(item *modelpb.ApiPerm, index int) *model.ApiPerm {
			return &model.ApiPerm{
				Path:   item.Path,
				Method: item.Method,
			}
		})
		var data []byte
		data, err = json.Marshal(apiPerms)
		if err != nil {
			return
		}
		update["data"] = data
	}
	if lo.Contains(req.UpdateMask, "desc") {
		update["desc"] = req.Permission.Desc
	}
	if len(update) == 0 {
		return
	}
	err = s.permissionDao.UpdateOne(ctx, permission, update)
	if err != nil {
		err = errors.Errorf("failed to update permission: %v", err)
		return
	}
	return
}

func (s *Service) DeletePermissions(ctx context.Context, req *model.DeletePermissionsRequest) (err error) {
	if len(req.IDs) == 0 {
		err = errors.New("ids is empty")
		return
	}
	if _, err = s.permissionDao.Delete(ctx, &model.PermissionFilter{IDs: db.In(req.IDs)}); err != nil {
		err = errors.Errorf("failed to delete permissions: %v", err)
		return
	}
	return
}

func (s *Service) GetPermissionList(ctx context.Context, req *model.GetPermissionListRequest) (total int64, list []*model.Permission, err error) {
	f := &model.PermissionFilter{
		Name: db.Like(req.Name+"%", db.OmitIfZero[string]()),
		Code: db.Like(req.Code+"%", db.OmitIfZero[string]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.permissionDao.Count(ctx, f)
		if err != nil {
			return
		}
		if !db.HasRecrods(total, req.PageParam.Page, req.PageParam.PageSize) {
			return
		}
		options = append(options,
			db.WithPaging(req.PageParam.Page, req.PageParam.PageSize),
			db.WithOrder(req.PageParam.OrderBy, nil),
		)
	}
	list, err = s.permissionDao.Find(ctx, f, options...)
	return
}
