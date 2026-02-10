package model

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
)

// AccountApiKey 账号API密钥
type AccountApiKey struct {
	db.Model

	AccountId  int64        `gorm:"type:bigint unsigned;not null;default:0;index:idx_account_status"`                              // 账号ID
	KeyName    string       `gorm:"type:varchar(50);not null;default:'';"`                                                         // API密钥名称
	KeyPrefix  string       `gorm:"type:varchar(20);not null;default:'';"`                                                         // API密钥前缀
	KeySuffix  string       `gorm:"type:varchar(10);not null;default:'';"`                                                         // API密钥后缀
	KeyHash    string       `gorm:"type:varchar(64);not null;default:'';uniqueIndex:uk_key_hash"`                                  // API密钥哈希
	Status     ApiKeyStatus `gorm:"type:enum('enabled','disabled','revoked');not null;default:'enabled';index:idx_account_status"` // 状态
	Scope      string       `gorm:"type:json;default:null"`                                                                        // 权限
	QuoteUsed  int64        `gorm:"type:bigint unsigned;not null;default:0"`                                                       // 已使用
	QuoteLimit *int64       `gorm:"type:bigint unsigned;default:null"`                                                             // 限额，null 不限
	RateLimit  *int         `gorm:"type:int unsigned;default:null"`                                                                // QPS 限流, null不限
	LastUsedAt *time.Time   `gorm:"type:datetime;default:null"`                                                                    // 最近一次使用时间
	ExpiredAt  *time.Time   `gorm:"type:datetime;default:null"`                                                                    // 过期时间

	Key string `gorm:"-"`
}

func (AccountApiKey) TableName() string {
	return TableAccountApiKey
}

type AccountApiKeyFilter struct {
	ID        db.F[int64]
	IDs       db.F[[]int64] `gorm:"column:id"`
	AccountId db.F[int64]
	KeyHash   db.F[string]
	Status    db.F[ApiKeyStatus]
	ExpiredAt db.F[*time.Time]
	Keyword   string
}

func (m *AccountApiKey) ToProto() *relaypb.AccountApiKey {
	return &relaypb.AccountApiKey{
		Id:         m.ID,
		AccountId:  m.AccountId,
		KeyName:    m.KeyName,
		Key:        fmt.Sprintf("%s...%s", m.KeyPrefix, m.KeySuffix),
		Status:     string(m.Status),
		Scope:      m.Scope,
		QuoteUsed:  m.QuoteUsed,
		QuoteLimit: lo.FromPtr(m.QuoteLimit),
		RateLimit:  int64(lo.FromPtr(m.RateLimit)),
		LastUsedAt: timestamppb.New(lo.FromPtr(m.LastUsedAt)),
		ExpiredAt:  timestamppb.New(lo.FromPtr(m.ExpiredAt)),
		CreatedAt:  timestamppb.New(m.CreatedAt),
		UpdatedAt:  timestamppb.New(m.UpdatedAt),
	}
}

type CreateAccountApiKeyRequest struct {
	AccountApiKey *relaypb.AccountApiKey
}

type UpdateAccountApiKeyRequest struct {
	AccountApiKey *relaypb.AccountApiKey
	UpdateMask    []string
}

type DeleteAccountApiKeysRequest struct {
	Ids []int64
}

type GetAccountApiKeyListRequest struct {
	*types.PageParam

	AccountId int64
	Keyword   string
	Status    ApiKeyStatus
}
