package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/system/model"
	"github.com/modelgate/modelgate/pkg/db"
	systempb "github.com/modelgate/modelgate/pkg/proto/model/system"
)

func (s *Service) CreateMenu(ctx context.Context, req *model.CreateMenuRequest) (menu *model.Menu, err error) {
	if req.Menu.Name == "" {
		err = errors.New("name is required")
		return
	}

	query := lo.Map(req.Menu.Query, func(query *systempb.MenuQuery, _ int) model.MenuQuery {
		return model.MenuQuery{
			Key:   query.Key,
			Value: query.Value,
		}
	})
	queryData, _ := json.Marshal(query)
	buttons := lo.Map(req.Menu.Buttons, func(button *systempb.MenuButton, _ int) model.MenuButton {
		return model.MenuButton{
			Code: button.Code,
			Desc: button.Desc,
		}
	})
	buttonsData, _ := json.Marshal(buttons)

	menu = &model.Menu{
		Pid:             req.Menu.Pid,
		Name:            req.Menu.Name,
		Type:            model.MenuType(req.Menu.Type),
		RouteName:       req.Menu.RouteName,
		RoutePath:       req.Menu.RoutePath,
		Component:       req.Menu.Component,
		Order:           int64(req.Menu.Order),
		I18nKey:         req.Menu.I18NKey,
		IconType:        model.IconType(req.Menu.IconType),
		Icon:            req.Menu.Icon,
		Status:          model.EnableStatus(req.Menu.Status),
		KeepAlive:       req.Menu.KeepAlive,
		Constant:        req.Menu.Constant,
		Href:            req.Menu.Href,
		HideInMenu:      req.Menu.HideInMenu,
		ActiveMenu:      req.Menu.ActiveMenu,
		MultiTab:        req.Menu.MultiTab,
		FixedIndexInTab: int64(req.Menu.FixedIndexInTab),
		Query:           queryData,
		Buttons:         buttonsData,
	}
	if err = s.menuDao.Create(ctx, menu); err != nil {
		err = errors.Errorf("failed to create menu: %v", err)
		return
	}
	return
}

func (s *Service) UpdateMenu(ctx context.Context, req *model.UpdateMenuRequest) (menu *model.Menu, err error) {
	menu, err = s.menuDao.FindOneByID(ctx, req.Menu.Id)
	if err != nil {
		return
	}
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "name") {
		update["name"] = req.Menu.Name
	}
	if lo.Contains(req.UpdateMask, "route_name") {
		update["route_name"] = req.Menu.RouteName
	}
	if lo.Contains(req.UpdateMask, "route_path") {
		update["route_path"] = req.Menu.RoutePath
	}
	if lo.Contains(req.UpdateMask, "component") {
		update["component"] = req.Menu.Component
	}
	if lo.Contains(req.UpdateMask, "i18n_key") {
		update["i18n_key"] = req.Menu.I18NKey
	}
	if lo.Contains(req.UpdateMask, "icon") {
		update["icon"] = req.Menu.Icon
	}
	if lo.Contains(req.UpdateMask, "icon_type") {
		update["icon_type"] = req.Menu.IconType
	}
	if lo.Contains(req.UpdateMask, "order") {
		update["order"] = req.Menu.Order
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.Menu.Status
	}
	if lo.Contains(req.UpdateMask, "query") {
		query := lo.Map(req.Menu.Query, func(query *systempb.MenuQuery, _ int) model.MenuQuery {
			return model.MenuQuery{
				Key:   query.Key,
				Value: query.Value,
			}
		})
		queryData, _ := json.Marshal(query)
		update["query"] = queryData
	}
	if lo.Contains(req.UpdateMask, "buttons") {
		buttons := lo.Map(req.Menu.Buttons, func(button *systempb.MenuButton, _ int) model.MenuButton {
			return model.MenuButton{
				Code: button.Code,
				Desc: button.Desc,
			}
		})
		buttonsData, _ := json.Marshal(buttons)
		update["buttons"] = buttonsData
	}
	if len(update) == 0 {
		return
	}
	if err = s.menuDao.UpdateOne(ctx, menu, update); err != nil {
		err = errors.Errorf("failed to update menu: %v", err)
		return
	}
	return
}

func (s *Service) DeleteMenus(ctx context.Context, req *model.DeleteMenusRequest) (err error) {
	if len(req.IDs) == 0 {
		err = errors.New("ids is empty")
		return
	}
	_, err = s.menuDao.Delete(ctx, &model.MenuFilter{IDs: db.In(req.IDs)})
	return err
}

func (s *Service) GetMenuList(ctx context.Context, req *model.GetMenuListRequest) (total int64, list []*model.Menu, err error) {
	f := &model.MenuFilter{
		Pid:    db.Eq(req.Pid, db.OmitIfZero[int64]()),
		Pids:   db.In(req.Pids, db.OmitIfZero[[]int64]()),
		Type:   db.Eq(req.Type, db.OmitIfZero[model.MenuType]()),
		Status: db.Eq(req.Status, db.OmitIfZero[model.EnableStatus]()),
	}
	if req.IsRoot {
		f.Pid = db.Eq(int64(0))
	} else if req.IsChild {
		f.Pid = db.Gt(int64(0))
	}
	if req.IsConstant {
		f.Constant = db.Eq(true)
	} else if req.IsUserRoute {
		f.Constant = db.Eq(false)
	}
	if req.HasButton {
		f.ButtonsLength = db.Gt(0)
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.menuDao.Count(ctx, f)
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
	list, err = s.menuDao.Find(ctx, f, options...)
	return
}

func (s *Service) GetMenu(ctx context.Context, req *model.GetMenuRequest) (menu *model.Menu, err error) {
	f := &model.MenuFilter{
		ID: db.Eq(req.ID),
	}
	menu, err = s.menuDao.FindOne(ctx, f)
	return
}
