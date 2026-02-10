package model

import (
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/utils"
)

// RequestAttempt 请求记录
type RequestAttempt struct {
	db.Model

	RequestUUID      utils.UUIDv7  `gorm:"type:binary(16);not null;uniqueIndex:uk_request_uuid_attempt_no"`
	AttemptNo        int           `gorm:"type:tinyint(4) unsigned;not null;default:0;uniqueIndex:uk_request_uuid_attempt_no"`
	ActualModel      string        `gorm:"type:varchar(100);not null;default:''"`
	AccountId        int64         `gorm:"type:bigint unsigned;not null;default:0"`
	AccountApiKeyId  int64         `gorm:"type:bigint unsigned;not null;default:0"`
	ProviderId       int64         `gorm:"type:bigint unsigned;not null;default:0"`
	ProviderApiKeyId int64         `gorm:"type:bigint unsigned;not null;default:0"`
	ModelId          int64         `gorm:"type:bigint unsigned;not null;default:0"`
	PromptTokens     int64         `gorm:"type:bigint unsigned;not null;default:0"`
	CompletionTokens int64         `gorm:"type:bigint unsigned;not null;default:0"`
	TotalTokens      int64         `gorm:"type:bigint unsigned;not null;default:0"`
	Status           RequestStatus `gorm:"type:enum('pending','success','failed','cancelled');not null;default:'pending'"`
	CompletedAt      *time.Time    `gorm:"type:datetime(3);"`
	ErrorCode        int           `gorm:"type:int unsigned;not null;default:0"`
	ErrorMessage     string        `gorm:"type:varchar(1000);not null;default:''"`
}

// TableName 表名
func (RequestAttempt) TableName() string {
	return TableRequestAttempt
}

// ToProto 转换为 proto
func (m *RequestAttempt) ToProto() *relaypb.RequestAttempt {
	info := &relaypb.RequestAttempt{
		Id:               m.ID,
		RequestUuid:      m.RequestUUID.String(),
		AttemptNo:        int64(m.AttemptNo),
		ActualModel:      m.ActualModel,
		AccountId:        m.AccountId,
		AccountApiKeyId:  m.AccountApiKeyId,
		ProviderId:       m.ProviderId,
		ProviderApiKeyId: m.ProviderApiKeyId,
		ModelId:          m.ModelId,
		PromptTokens:     m.PromptTokens,
		CompletionTokens: m.CompletionTokens,
		TotalTokens:      m.TotalTokens,
		Status:           string(m.Status),
		ErrorCode:        int64(m.ErrorCode),
		ErrorMessage:     m.ErrorMessage,
		CreatedAt:        timestamppb.New(m.CreatedAt),
		UpdatedAt:        timestamppb.New(m.UpdatedAt),
	}
	if m.CompletedAt != nil {
		info.CompletedAt = timestamppb.New(*m.CompletedAt)
		info.ElapsedTime = lo.FromPtr(m.CompletedAt).Sub(m.CreatedAt).Milliseconds()
	}
	return info
}

type RequestAttemptFilter struct {
	ID           db.F[int64]
	IDs          db.F[[]int64] `gorm:"column:id"`
	RequestUUID  db.F[utils.UUIDv7]
	AttemptNo    db.F[int]
	Object       db.F[string]
	AccountId    db.F[int64]
	ProviderId   db.F[int64]
	ModelId      db.F[int64]
	Status       db.F[RequestStatus]
	CompletedAt  db.F[time.Time]
	ProviderCode db.F[string]
	ModelCode    db.F[string]
}
