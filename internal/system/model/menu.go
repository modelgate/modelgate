package model

import (
	"encoding/json"

	"github.com/samber/lo"
	"gorm.io/datatypes"

	"github.com/modelgate/modelgate/pkg/db"
	systempb "github.com/modelgate/modelgate/pkg/proto/model/system"
	"github.com/modelgate/modelgate/pkg/types"
)

type MenuButton struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type MenuQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MenuType uint8

const (
	MenuTypeDirectory MenuType = iota + 1
	MenuTypeMenu
)

type IconType uint8

const (
	IconTypeIconify IconType = iota + 1
	IconTypeLocal
)

// Menu 菜单
type Menu struct {
	db.Model

	Pid             int64          `gorm:"type:bigint unsigned;not null;default:0"`                          // 父菜单id
	Type            MenuType       `gorm:"type:tinyint unsigned;not null;default:1"`                         // 菜单类型
	Name            string         `gorm:"type:varchar(200);not null;default:''"`                            // 菜单名称
	RouteName       string         `gorm:"type:varchar(200);not null;default:'';uniqueIndex:uk_router_name"` // 路由名称
	RoutePath       string         `gorm:"type:varchar(1000);not null;default:''"`                           // 路由路径
	Component       string         `gorm:"type:varchar(1000);not null;default:''"`                           // 页面组件
	I18nKey         string         `gorm:"type:varchar(200);not null;default:''"`                            // 国际化key
	Order           int64          `gorm:"type:bigint unsigned;not null;default:0"`                          // 排序
	IconType        IconType       `gorm:"type:tinyint unsigned;not null;default:1"`                         // 图标类型
	Icon            string         `gorm:"type:varchar(200);not null;default:''"`                            // 图标
	Status          EnableStatus   `gorm:"type:enum('enabled', 'disabled');not null;default:'enabled'"`      // 状态
	KeepAlive       bool           `gorm:"type:tinyint unsigned;not null;default:0"`                         // 缓存路由
	Constant        bool           `gorm:"type:tinyint unsigned;not null;default:0"`                         // 常量路由
	Href            string         `gorm:"type:varchar(1000);not null;default:''"`                           // 外链
	HideInMenu      bool           `gorm:"type:tinyint unsigned;not null;default:0"`                         // 隐藏菜单
	ActiveMenu      string         `gorm:"type:varchar(1000);not null;default:''"`                           // 激活菜单
	MultiTab        bool           `gorm:"type:tinyint unsigned;not null;default:0"`                         // 支持多页签
	FixedIndexInTab int64          `gorm:"type:bigint unsigned;not null;default:0"`                          // 固定索引在标签页
	Query           datatypes.JSON `gorm:"type:json"`                                                        // 路由参数
	Buttons         datatypes.JSON `gorm:"type:json"`                                                        // 按钮
}

func (Menu) TableName() string {
	return TableMenus
}

func (m *Menu) ToProto() *systempb.Menu {
	info := &systempb.Menu{
		Id:              m.ID,
		Pid:             m.Pid,
		Name:            m.Name,
		Type:            systempb.MenuType(m.Type),
		RouteName:       m.RouteName,
		RoutePath:       m.RoutePath,
		Component:       m.Component,
		I18NKey:         m.I18nKey,
		Order:           uint32(m.Order),
		Icon:            m.Icon,
		IconType:        systempb.IconType(m.IconType),
		Status:          string(m.Status),
		KeepAlive:       m.KeepAlive,
		Constant:        m.Constant,
		Href:            m.Href,
		HideInMenu:      m.HideInMenu,
		ActiveMenu:      m.ActiveMenu,
		MultiTab:        m.MultiTab,
		FixedIndexInTab: uint32(m.FixedIndexInTab),
	}
	var query []MenuQuery
	if err := json.Unmarshal(m.Query, &query); err == nil {
		info.Query = lo.Map(query, func(query MenuQuery, _ int) *systempb.MenuQuery {
			return &systempb.MenuQuery{
				Key:   query.Key,
				Value: query.Value,
			}
		})
	}
	var buttons []MenuButton
	if err := json.Unmarshal(m.Buttons, &buttons); err == nil {
		info.Buttons = lo.Map(buttons, func(button MenuButton, _ int) *systempb.MenuButton {
			return &systempb.MenuButton{
				Code: button.Code,
				Desc: button.Desc,
			}
		})
	}
	return info
}

func (m *Menu) ToMenuRouteProto() *systempb.MenuRoute {
	return &systempb.MenuRoute{
		Id:        m.ID,
		Name:      m.RouteName,
		Path:      m.RoutePath,
		Component: m.Component,
		Meta: &systempb.MenuMeta{
			Title:      m.RouteName,
			Constant:   m.Constant,
			HideInMenu: m.HideInMenu,
			I18NKey:    m.I18nKey,
			Icon:       m.Icon,
			Order:      uint32(m.Order),
		},
	}
}

type MenuFilter struct {
	ID            db.F[int64]
	IDs           db.F[[]int64] `gorm:"column:id"`
	Pid           db.F[int64]
	Pids          db.F[[]int64] `gorm:"column:pid"`
	Type          db.F[MenuType]
	ButtonsLength db.F[int] `gorm:"column:JSON_LENGTH(buttons)"`
	Constant      db.F[bool]
	Status        db.F[EnableStatus]
}

type CreateMenuRequest struct {
	Menu *systempb.Menu
}

type UpdateMenuRequest struct {
	Menu       *systempb.Menu
	UpdateMask []string
}

type DeleteMenusRequest struct {
	IDs []int64
}

type GetMenuListRequest struct {
	*types.PageParam

	IsRoot      bool
	IsChild     bool
	IsConstant  bool
	IsUserRoute bool
	HasButton   bool
	Pid         int64
	Pids        []int64
	Name        string
	Status      EnableStatus
	Type        MenuType
}

type GetMenuRequest struct {
	ID int64
}

type GetPageListRequest struct {
	*types.PageParam
}
