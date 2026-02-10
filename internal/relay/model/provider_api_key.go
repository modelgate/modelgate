package model

import (
	"fmt"
	"time"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProviderApiKey
type ProviderApiKey struct {
	db.Model

	ProviderId    int64        `gorm:"type:bigint;not null;default:0;index:idx_provider_status"`                                        // 供应商ID
	ProviderCode  string       `gorm:"type:varchar(50);not null;default:''"`                                                            // 供应商代码
	Name          string       `gorm:"type:varchar(100);not null;default:''"`                                                           // API Key名称
	KeyPrefix     string       `gorm:"type:varchar(20);not null;default:''"`                                                            // API Key前缀
	KeySuffix     string       `gorm:"type:varchar(10);not null;default:''"`                                                            // API Key后缀
	KeyEncrypted  string       `gorm:"type:varchar(512);not null;default:''"`                                                           // API Key值
	Weight        int          `gorm:"type:int;not null;default:100"`                                                                   // 权重
	Status        ApiKeyStatus `gorm:"type:enum('enabled','disabled','cooldown');not null;default:'enabled';index:idx_provider_status"` // 状态
	QuoteUsed     int64        `gorm:"type:bigint unsigned;not null;default:0"`                                                         // 已使用
	QuoteLimit    *int64       `gorm:"type:bigint unsigned;default:null"`                                                               // 限额，null 不限
	RateLimit     *int         `gorm:"type:int unsigned;default:null"`                                                                  // QPS 限流, null不限
	LastUsedAt    *time.Time   `gorm:"type:datetime"`                                                                                   // 最后使用时间
	FaildCount    int          `gorm:"type:int;not null;default:0"`                                                                     // 失败次数
	CoolDownUntil *time.Time   `gorm:"type:datetime"`                                                                                   // 冷却结束时间
}

// TableName 表名
func (ProviderApiKey) TableName() string {
	return TableProviderApiKey
}

func (m *ProviderApiKey) GetWeight() int {
	return m.Weight
}

func (m *ProviderApiKey) ToProto() *relaypb.ProviderApiKey {
	item := &relaypb.ProviderApiKey{
		Id:            m.ID,
		ProviderId:    m.ProviderId,
		ProviderCode:  m.ProviderCode,
		Name:          m.Name,
		Key:           fmt.Sprintf("%s...%s", m.KeyPrefix, m.KeySuffix),
		Weight:        int64(m.Weight),
		Status:        string(m.Status),
		FaildCount:    int64(m.FaildCount),
		LastUsedAt:    timestamppb.New(lo.FromPtr(m.LastUsedAt)),
		CoolDownUntil: timestamppb.New(lo.FromPtr(m.CoolDownUntil)),
		CreatedAt:     timestamppb.New(m.CreatedAt),
		UpdatedAt:     timestamppb.New(m.UpdatedAt),
	}
	return item
}

// ProviderApiKeyFilter 过滤器
type ProviderApiKeyFilter struct {
	ID           db.F[int64]
	IDs          db.F[[]int64] `gorm:"column:id"`
	ProviderId   db.F[int64]
	ProviderCode db.F[string]
	Name         db.F[string]
	Status       db.F[ApiKeyStatus]
}

type CreateProviderApiKeyRequest struct {
	ProviderApiKey *relaypb.ProviderApiKey
}

type UpdateProviderApiKeyRequest struct {
	ProviderApiKey *relaypb.ProviderApiKey
	UpdateMask     []string
}

type DeleteProviderApiKeysRequest struct {
	Ids []int64
}

type GetProviderApiKeyListRequest struct {
	*types.PageParam

	ProviderId   int64
	ProviderCode string
	Name         string
	Status       ApiKeyStatus
}
