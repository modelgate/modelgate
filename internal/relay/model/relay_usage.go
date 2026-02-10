package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
)

// RelayUsage 转发使用量
type RelayUsage struct {
	db.Model

	TotalRequest int64 `gorm:"type:bigint;not null;default:0"` // 总请求数
	TotalPoint   int64 `gorm:"type:bigint;not null;default:0"` // 总点数
}

func (RelayUsage) TableName() string {
	return TableRelayUsage
}

func (r *RelayUsage) ToProto() *relaypb.RelayUsage {
	return &relaypb.RelayUsage{
		Id:           r.ID,
		TotalRequest: r.TotalRequest,
		TotalPoint:   r.TotalPoint,
		CreatedAt:    timestamppb.New(r.CreatedAt),
		UpdatedAt:    timestamppb.New(r.UpdatedAt),
	}
}

type RelayUsageFilter struct {
	ID db.F[int64]
}

type RelayUsageInfo struct {
	TotalRequest int64 // 总请求数
	TotalPoint   int64 // 总点数
}
