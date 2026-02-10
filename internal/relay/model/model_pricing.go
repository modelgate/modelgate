package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
)

// ModelPricing 模型价格
// 同一个模型如果有 2 个价格，取生效时间最近的且未失效的价格
type ModelPricing struct {
	ID        int64     `gorm:"type:bigint unsigned;primaryKey" json:"id,string"`                                    // 主键ID
	CreatedAt time.Time `gorm:"type: datetime;not null;default: CURRENT_TIMESTAMP;autoCreateTime" json:"created_at"` // 创建时间

	ProviderCode      string       `gorm:"type:varchar(50);not null;default:'';uniqueIndex:uk_provider_model_effective"`  // 供应商 code
	ModelCode         string       `gorm:"type:varchar(100);not null;default:'';uniqueIndex:uk_provider_model_effective"` // 模型 code
	Currency          Currency     `gorm:"type:enum('USD','CNY','POINT');not null;default:'USD'"`                         // 货币单位
	PointsPerCurrency int64        `gorm:"type:bigint unsigned;not null;default:1"`                                       // 1 单位货币 = N Points
	TokenNum          int64        `gorm:"type:bigint unsigned;not null;default:1000000"`                                 // 价格对应的 token 数，如: 1_000_000
	InputPrice        float64      `gorm:"type:decimal(10,4) unsigned;not null;default:0"`                                // 每 1 token unit input token 价格
	InputCachePrice   float64      `gorm:"type:decimal(10,4) unsigned;not null;default:0"`                                // 每 1 token unit input token 缓存价格
	OutputPrice       float64      `gorm:"type:decimal(10,4) unsigned;not null;default:0"`                                // 每 1 token unit output token 价格
	EffectiveFrom     time.Time    `gorm:"type:datetime;not null;uniqueIndex:uk_provider_model_effective"`                // 生效时间
	EffectiveTo       time.Time    `gorm:"type:datetime;not null;" json:"effective_to"`                                   // 失效时间
	Status            EnableStatus `gorm:"type:enum('enabled','disabled');not null;default:'enabled'"`                    // 状态
}

func (ModelPricing) TableName() string {
	return TableModelPricing
}

func (m *ModelPricing) ToProto() *relaypb.ModelPricing {
	return &relaypb.ModelPricing{
		Id:                m.ID,
		ProviderCode:      m.ProviderCode,
		ModelCode:         m.ModelCode,
		Currency:          string(m.Currency),
		PointsPerCurrency: m.PointsPerCurrency,
		TokenNum:          m.TokenNum,
		InputPrice:        float32(m.InputPrice),
		InputCachePrice:   float32(m.InputCachePrice),
		OutputPrice:       float32(m.OutputPrice),
		EffectiveFrom:     timestamppb.New(m.EffectiveFrom),
		EffectiveTo:       timestamppb.New(m.EffectiveTo),
		Status:            string(m.Status),
		CreatedAt:         timestamppb.New(m.CreatedAt),
	}
}

// ModelPricingFilter 过滤器
type ModelPricingFilter struct {
	ID            db.F[int64]
	IDs           db.F[[]int64] `gorm:"column:id"`
	ProviderCode  db.F[string]
	ModelCode     db.F[string]
	Currency      db.F[Currency]
	Status        db.F[EnableStatus]
	EffectiveFrom db.F[time.Time]
	EffectiveTo   db.F[time.Time]
}

type CreateModelPricingRequest struct {
	ModelPricing *relaypb.ModelPricing
}

type UpdateModelPricingRequest struct {
	ModelPricing *relaypb.ModelPricing
	UpdateMask   []string
}

type DeleteModelPricingsRequest struct {
	Ids []int64
}

type GetModelPricingListRequest struct {
	*types.PageParam

	ProviderCode  string
	ModelCode     string
	Currency      Currency
	Status        EnableStatus
	EffectiveFrom time.Time
	EffectiveTo   time.Time
}
