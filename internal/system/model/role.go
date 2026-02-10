package model

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"

	"github.com/modelgate/modelgate/pkg/db"
	systempb "github.com/modelgate/modelgate/pkg/proto/model/system"
	"github.com/modelgate/modelgate/pkg/types"
)

// Role 角色
type Role struct {
	db.Model
	Name         string         `gorm:"type:varchar(100);not null;default:'';"`
	Code         string         `gorm:"type:varchar(100);not null;default:'';uniqueIndex:uk_role_code;"`
	IsSuperAdmin bool           `gorm:"type:tinyint(1);not null;default:0"`
	Description  string         `gorm:"type:varchar(1000);default:''"`
	Status       EnableStatus   `gorm:"type:enum('enabled', 'disabled');not null;default:'enabled'"`
	Permission   datatypes.JSON `gorm:"type:json"`
}

func (Role) TableName() string {
	return TableRoles
}

func (r *Role) ToProto() *systempb.Role {
	info := &systempb.Role{
		Id:           r.ID,
		Name:         r.Name,
		Code:         r.Code,
		Description:  r.Description,
		IsSuperAdmin: r.IsSuperAdmin,
		Status:       string(r.Status),
		CreatedAt:    timestamppb.New(r.CreatedAt),
		UpdatedAt:    timestamppb.New(r.UpdatedAt),
	}
	return info
}

func (r *Role) GetPermission() (*RolePermission, error) {
	var permission RolePermission
	if err := json.Unmarshal(r.Permission, &permission); err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *Role) SetPermission(permission *RolePermission) error {
	data, err := json.Marshal(permission)
	if err != nil {
		return err
	}
	r.Permission = datatypes.JSON(data)
	return nil
}

type RolePermission struct {
	Home    string   `json:"home"`
	MenuIds []int64  `json:"menu_ids"`
	Buttons []string `json:"buttons"`
}

func (p *RolePermission) ToProto() *systempb.RolePermission {
	return &systempb.RolePermission{
		Home:    p.Home,
		MenuIds: p.MenuIds,
		Buttons: p.Buttons,
	}
}

type RoleFilter struct {
	ID          db.F[int64]
	IDs         db.F[[]int64] `gorm:"column:id"`
	Name        db.F[string]
	Code        db.F[string]
	Description db.F[string]
	Status      db.F[EnableStatus]
}

type CreateRoleRequest struct {
	Role *systempb.Role
}

type UpdateRoleRequest struct {
	Role       *systempb.Role
	UpdateMask []string
}

type UpdateRolePermissionRequest struct {
	ID         int64
	Home       string
	MenuIds    []int64
	Buttons    []string
	UpdateMask []string
}

type DeleteRolesRequest struct {
	IDs []int64
}

type GetRoleListRequest struct {
	*types.PageParam

	IDs         []int64
	Name        string
	Code        string
	Description string
	Status      EnableStatus
}

type GetRoleRequest struct {
	ID int64
}
