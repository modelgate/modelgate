package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/modelgate/modelgate/pkg/db"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
)

// 系统内部统一使用点数作为计费单位，点数与人民币的换算比例为 1_000_000:1（1_000_000 点 = 1 元）。

// 平台当前支持多种大模型服务，不同模型提供方采用的计费币种并不一致：
// - 部分模型（如 DeepSeek）使用 **人民币（CNY）**计费
// - 部分模型（如 OpenAI）使用 **美元（USD）**计费

// 平台采用 点数制计费模型，将不同币种的模型成本统一折算为点数进行结算。
// 点数是平台内部唯一的计费与扣减单位，用户账户余额仅以点数形式存在。
// 点数与人民币的换算规则
// 1_000_000 点 = 1 元人民币
// 即：1 元人民币 = 1_000_000 点
// 该比例固定不变，用于保证计费精度和系统稳定性。

// Account 账户
type Account struct {
	db.Model

	Nickname string       `gorm:"type:varchar(64);not null;default:'';"`                     // 昵称
	Name     string       `gorm:"type:varchar(64);not null;default:'';uniqueIndex:uk_name"`  // 账户名称
	Balance  int64        `gorm:"type:bigint unsigned;not null;default:0;"`                  // 余额: 点数
	Status   EnableStatus `gorm:"type:enum('enabled','disabled');not null;default:enabled;"` // 状态
}

func (Account) TableName() string {
	return TableAccount
}

type AccountFilter struct {
	ID       db.F[int64]
	IDs      db.F[[]int64] `gorm:"column:id"`
	Name     db.F[string]
	Nickname db.F[string]
	Status   db.F[EnableStatus]
}

func (m *Account) ToProto() *relaypb.Account {
	return &relaypb.Account{
		Id:        m.ID,
		Nickname:  m.Nickname,
		Name:      m.Name,
		Balance:   m.Balance,
		Status:    string(m.Status),
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

type CreateAccountRequest struct {
	Account *relaypb.Account
}

type UpdateAccountRequest struct {
	Account    *relaypb.Account
	UpdateMask []string
}

type DeleteAccountsRequest struct {
	Ids []int64
}

type GetAccountListRequest struct {
	*types.PageParam

	Id       int64
	Ids      []int64
	Name     string
	Nickname string
	Status   EnableStatus
}
