package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) CreateProvider(ctx context.Context, req *model.CreateProviderRequest) (info *model.Provider, err error) {
	info, err = s.providerDao.FindOne(ctx, &model.ProviderFilter{Code: db.Eq(req.Provider.Code)})
	if db.IsDbError(err) {
		return
	}
	if info != nil {
		err = fmt.Errorf("provider already exists, code: %s", req.Provider.Code)
		return
	}
	info = &model.Provider{
		Name:    req.Provider.Name,
		Code:    req.Provider.Code,
		BaseUrl: req.Provider.BaseUrl,
		Status:  model.EnableStatus(req.Provider.Status),
	}
	err = s.providerDao.Create(ctx, info)
	return
}

func (s *Service) UpdateProvider(ctx context.Context, req *model.UpdateProviderRequest) (info *model.Provider, err error) {
	info, err = s.providerDao.FindOneByID(ctx, req.Provider.Id)
	if db.IsDbError(err) {
		return
	}
	if info == nil {
		err = fmt.Errorf("provider not found, id: %d", req.Provider.Id)
		return
	}
	isChanged := lo.Ternary(info.Code != req.Provider.Code, true, false)
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "name") {
		update["name"] = req.Provider.Name
	}
	if lo.Contains(req.UpdateMask, "code") {
		update["code"] = req.Provider.Code
	}
	if lo.Contains(req.UpdateMask, "base_url") {
		update["base_url"] = req.Provider.BaseUrl
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.Provider.Status
	}
	if len(update) == 0 {
		err = fmt.Errorf("no fields to update")
		return
	}
	err = s.providerDao.UpdateOne(ctx, info, update)
	if err != nil {
		return
	}
	if isChanged {
		_, err = s.modelDao.Update(ctx, &model.ModelFilter{ProviderId: db.Eq(info.ID)}, map[string]any{"provider_code": info.Code})
		if db.IsDbError(err) {
			return
		}
		_, err = s.providerApiKeyDao.Update(ctx, &model.ProviderApiKeyFilter{ProviderId: db.Eq(info.ID)}, map[string]any{"provider_code": info.Code})
		if db.IsDbError(err) {
			return
		}
	}
	return
}

func (s *Service) DeleteProviders(ctx context.Context, req *model.DeleteProvidersRequest) (err error) {
	_, err = s.providerDao.Delete(ctx, &model.ProviderFilter{IDs: db.In(req.Ids)})
	return
}

func (s *Service) GetProviderList(ctx context.Context, req *model.GetProviderListRequest) (total int64, list []*model.Provider, err error) {
	f := &model.ProviderFilter{
		Name:   db.Like(req.Name+"%", db.OmitIf(func(s string) bool { return s == "%" })),
		Code:   db.Eq(req.Code, db.OmitIfZero[string]()),
		Status: db.Eq(req.Status, db.OmitIfZero[model.EnableStatus]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.providerDao.Count(ctx, f)
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
	list, err = s.providerDao.Find(ctx, f, options...)
	return
}
