package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
)

// RelayHourlyUsage 每小时使用量
type RelayHourlyUsage struct {
	db.Model

	Time         time.Time `gorm:"type:datetime;not null;uniqueIndex:uk_time_provider"`
	ProviderCode string    `gorm:"type:varchar(64);not null;default:'';uniqueIndex:uk_time_provider"`
	TotalRequest int64     `gorm:"type:bigint;not null;default:0"` // 总请求数
	TotalSuccess int64     `gorm:"type:bigint;not null;default:0"` // 总成功数
	TotalFailed  int64     `gorm:"type:bigint;not null;default:0"` // 总失败数
	TotalPoint   int64     `gorm:"type:bigint;not null;default:0"` // 总点数
}

func (RelayHourlyUsage) TableName() string {
	return TableRelayHourlyUsage
}

func (r *RelayHourlyUsage) ToProto() *relaypb.RelayUsage {
	return &relaypb.RelayUsage{
		Id:           r.ID,
		TotalRequest: r.TotalRequest,
		TotalSuccess: r.TotalSuccess,
		TotalFailed:  r.TotalFailed,
		TotalPoint:   r.TotalPoint,
		CreatedAt:    timestamppb.New(r.CreatedAt),
		UpdatedAt:    timestamppb.New(r.UpdatedAt),
	}
}

type RelayHourlyUsageFilter struct {
	Time db.F[time.Time]
}

type GetRelayHourlyUsageListRequest struct {
	*types.PageParam

	StartTime time.Time
	EndTime   time.Time
}
