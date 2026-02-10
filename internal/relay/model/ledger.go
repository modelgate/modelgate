package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
)

type LedgerType string

const (
	LedgerTypeConsume LedgerType = "consume" // 消费
	LedgerTypeRefund  LedgerType = "refund"  // 退款
	LedgerTypeCharge  LedgerType = "charge"  // 充值
	LedgerTypeAdjust  LedgerType = "adjust"  // 调整
)

// Ledger 账户流水
type Ledger struct {
	ID        int64     `gorm:"type:bigint unsigned;primaryKey"`                                                                // 主键ID
	CreatedAt time.Time `gorm:"type: datetime;not null;default: CURRENT_TIMESTAMP;autoCreateTime; index:idx_account_type_time"` // 创建时间

	AccountId    int64      `gorm:"type:bigint unsigned;not null;default:0; index:idx_account_type_time"`                                    // 用户ID
	Type         LedgerType `gorm:"type:enum('consume','refund','charge','adjust');not null;default:'consume'; index:idx_account_type_time"` // 类型
	Amount       int64      `gorm:"type:bigint;not null;default:0"`                                                                          // 调整值
	BalanceAfter int64      `gorm:"type:bigint unsigned;not null;default:0"`                                                                 // 余额
	RequestId    int64      `gorm:"type:bigint unsigned;not null;default:0"`                                                                 // 请求ID
	Reason       string     `gorm:"type:varchar(255);not null;default:''"`                                                                   // 原因
}

func (Ledger) TableName() string {
	return TableLedger
}

// LedgerFilter 过滤器
type LedgerFilter struct {
	ID        db.F[int64]
	IDs       db.F[[]int64] `gorm:"column:id"`
	AccountId db.F[int64]
	Type      db.F[LedgerType]
}

func (m *Ledger) ToProto() *relaypb.Ledger {
	return &relaypb.Ledger{
		Id:           m.ID,
		AccountId:    m.AccountId,
		Type:         string(m.Type),
		Amount:       m.Amount,
		BalanceAfter: m.BalanceAfter,
		RequestId:    m.RequestId,
		Reason:       m.Reason,
		CreatedAt:    timestamppb.New(m.CreatedAt),
	}
}

type CreateLedgerRequest struct {
	Ledger *relaypb.Ledger
}

type DeleteLedgersRequest struct {
	Ids []int64
}

type GetLedgerListRequest struct {
	*types.PageParam

	AccountId int64
	Type      LedgerType
}
