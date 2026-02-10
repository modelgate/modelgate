package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
	"github.com/samber/lo"
)

type ModelStatus string

const (
	ModelStatusEnabled    ModelStatus = "enabled"    // 正常
	ModelStatusDisabled   ModelStatus = "disabled"   // 禁用
	ModelStatusDeprecated ModelStatus = "deprecated" // 废弃
)

// Model 模型
type Model struct {
	db.Model

	ProviderId   int64       `gorm:"type:bigint unsigned;not null;default:0"`                                          // 供应商ID
	ProviderCode string      `gorm:"type:varchar(50);not null;default:'';uniqueIndex:uk_provider_model"`               // 供应商代码
	Code         string      `gorm:"type:varchar(50);not null;default:'';uniqueIndex:uk_provider_model"`               // 模型代码
	ActualCode   string      `gorm:"type:varchar(50);not null;default:''"`                                             // 实际模型代码
	Name         string      `gorm:"type:varchar(100);not null;default:'';index:idx_name_status,priority:2,length:32"` // 模型名称
	Priority     int         `gorm:"type:int;not null;default:1"`                                                      // 优先级，越小越优先
	Weight       int         `gorm:"type:int;not null;default:100"`                                                    // 权重
	Status       ModelStatus `gorm:"type:enum('enabled','disabled','deprecated');not null;default:'enabled'"`          // 状态
}

func (Model) TableName() string {
	return TableModel
}

func (m *Model) GetWeight() int {
	return m.Weight
}

// GetActualCode 模型真实代码
func (m *Model) GetActualCode() string {
	return lo.Ternary(m.ActualCode != "", m.ActualCode, m.Code)
}

func (m *Model) ToProto() *relaypb.Model {
	return &relaypb.Model{
		Id:           m.ID,
		ProviderId:   m.ProviderId,
		ProviderCode: m.ProviderCode,
		Code:         m.Code,
		ActualCode:   m.ActualCode,
		Name:         m.Name,
		Priority:     int64(m.Priority),
		Weight:       int64(m.Weight),
		Status:       string(m.Status),
		CreatedAt:    timestamppb.New(m.CreatedAt),
		UpdatedAt:    timestamppb.New(m.UpdatedAt),
	}
}

// ModelFilter 过滤器
type ModelFilter struct {
	ID           db.F[int64]
	IDs          db.F[[]int64] `gorm:"column:id"`
	ProviderId   db.F[int64]
	ProviderCode db.F[string]
	VirtualCode  db.F[string]
	Code         db.F[string]
	Name         db.F[string]
	Status       db.F[ModelStatus]
}

type CreateModelRequest struct {
	Model *relaypb.Model
}

type UpdateModelRequest struct {
	Model      *relaypb.Model
	UpdateMask []string
}

type DeleteModelsRequest struct {
	Ids []int64
}

type GetModelListRequest struct {
	*types.PageParam

	Ids          []int64
	Name         string
	Code         string
	ProviderCode string
	VirtualCode  string
	Status       ModelStatus
}

// ResolvedModel 解析后的模型
type ResolvedModel struct {
	ModelId         int64  // ID
	ModelCode       string // 模型Code
	ProviderId      int64  // 提供商 ID
	ProviderCode    string // 提供商Code
	BaseUrl         string // 基础URL
	ApiKeyId        int64  // 提供商 API Key ID
	ApiKeyEncrypted string // 加密后的 API Key

	InputPrice        float64 // 输入价格
	InputCachePrice   float64 // 输入缓存价格
	OutputPrice       float64 // 输出价格
	TokenNum          int64   // Token 数量
	PointsPerCurrency int64   // 每个货币点数
}
