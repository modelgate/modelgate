package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
)

// Provider 供应商
type Provider struct {
	db.Model

	Code    string       `gorm:"type:varchar(50);not null;default:'';uniqueIndex:uk_code"`                  // 供应商代码
	Name    string       `gorm:"type:varchar(100);not null;default:'';index:idx_name,priority:1,length:32"` // 供应商名称
	BaseUrl string       `gorm:"type:varchar(255);not null;default:''"`                                     // 接口URL
	Status  EnableStatus `gorm:"type:enum('enabled','disabled');not null;default:'enabled'"`                // 状态
}

func (Provider) TableName() string {
	return TableProvider
}

func (m *Provider) ToProto() *relaypb.Provider {
	return &relaypb.Provider{
		Id:        m.ID,
		Code:      m.Code,
		Name:      m.Name,
		BaseUrl:   m.BaseUrl,
		Status:    string(m.Status),
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

// ProviderFilter 过滤器
type ProviderFilter struct {
	ID     db.F[int64]
	IDs    db.F[[]int64] `gorm:"column:id"`
	Code   db.F[string]
	Name   db.F[string]
	Status db.F[EnableStatus]
}

type CreateProviderRequest struct {
	Provider *relaypb.Provider
}

type UpdateProviderRequest struct {
	Provider   *relaypb.Provider
	UpdateMask []string
}

type DeleteProvidersRequest struct {
	Ids []int64
}

type GetProviderListRequest struct {
	*types.PageParam

	Ids    []int64
	Name   string
	Code   string
	Status EnableStatus
}
