package dao

import (
	"context"
	"errors"

	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type AccountDao struct {
	*db.BaseDAO[model.Account, model.AccountFilter]
}

func NewAccountDao(i do.Injector) (relay.AccountDAO, error) {
	dbConn := do.MustInvoke[*gorm.DB](i)
	return &AccountDao{
		BaseDAO: db.NewBaseDAO[model.Account, model.AccountFilter](dbConn),
	}, nil
}

// DeductBalance 事务扣减账户余额
func (d *AccountDao) DeductBalance(ctx context.Context, accountId int64, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error) {
	if amount <= 0 {
		err = errors.New("amount must be greater than 0")
		return
	}
	switch typ {
	case model.LedgerTypeConsume,
		model.LedgerTypeAdjust:
	default:
		err = errors.New("invalid ledger type")
		return
	}
	err = d.GetDB().Transaction(func(tx *gorm.DB) (tErr error) {
		res := tx.Model(&model.Account{}).Where("id = ? and balance >= ?", accountId, amount).Update("balance", gorm.Expr("balance - ?", amount))
		if res.Error != nil {
			return
		} else if res.RowsAffected == 0 {
			err = errors.New("insufficient balance")
			return
		}
		var account *model.Account
		tErr = tx.Model(&model.Account{}).Where("id = ?", accountId).First(&account).Error
		if tErr != nil {
			return
		}
		ledger = &model.Ledger{
			AccountId:    accountId,
			Type:         typ,
			Amount:       -amount,
			BalanceAfter: account.Balance,
			RequestId:    requestId,
			Reason:       reason,
		}
		if tErr = tx.Create(ledger).Error; err != nil {
			return
		}
		return nil
	})
	return
}

// IncreaseBalance 事务增加账户余额
func (d *AccountDao) IncreaseBalance(ctx context.Context, accountId int64, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error) {
	if amount <= 0 {
		err = errors.New("amount must be greater than 0")
		return
	}
	switch typ {
	case model.LedgerTypeRefund,
		model.LedgerTypeCharge,
		model.LedgerTypeAdjust:
	default:
		err = errors.New("invalid ledger type")
		return
	}
	err = d.GetDB().Transaction(func(tx *gorm.DB) (tErr error) {
		tErr = tx.Model(&model.Account{}).Where("id = ?", accountId).Update("balance", gorm.Expr("balance + ?", amount)).Error
		if tErr != nil {
			return
		}
		var account *model.Account
		tErr = tx.Model(&model.Account{}).Where("id = ?", accountId).First(&account).Error
		if tErr != nil {
			return
		}
		ledger = &model.Ledger{
			AccountId:    accountId,
			Amount:       amount,
			BalanceAfter: account.Balance,
			Type:         typ,
			RequestId:    requestId,
			Reason:       reason,
		}
		tErr = tx.Create(ledger).Error
		return tErr
	})
	return
}
