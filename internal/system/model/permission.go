package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"

	"github.com/modelgate/modelgate/pkg/db"
	systempb "github.com/modelgate/modelgate/pkg/proto/model/system"
	"github.com/modelgate/modelgate/pkg/types"
)

// Permission 权限
type Permission struct {
	db.Model

	Name string         `gorm:"type:varchar(100);not null;default:'';"`
	Code string         `gorm:"type:varchar(100);not null;default:'';uniqueIndex:uk_code"`
	Data datatypes.JSON `gorm:"type:json"`
	Desc string         `gorm:"type:varchar(100);not null;default:'';"`
}

func (Permission) TableName() string {
	return TablePermissions
}

func (m *Permission) ToProto() *systempb.Permission {
	return &systempb.Permission{
		Id:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		Desc:      m.Desc,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

type ApiPerm struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

func (a *ApiPerm) ToProto() *systempb.ApiPerm {
	return &systempb.ApiPerm{
		Path:   a.Path,
		Method: a.Method,
	}
}

type PermissionFilter struct {
	ID    db.F[int64]
	IDs   db.F[[]int64] `gorm:"column:id"`
	Name  db.F[string]
	Code  db.F[string]
	Codes db.F[[]string] `gorm:"column:code"`
}

type CreatePermissionRequest struct {
	Permission *systempb.Permission
}

type UpdatePermissionRequest struct {
	Permission *systempb.Permission
	UpdateMask []string
}

type DeletePermissionsRequest struct {
	IDs []int64
}

type GetPermissionListRequest struct {
	*types.PageParam

	Name string
	Code string
}
