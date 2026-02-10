package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) DeductBalance(ctx context.Context, accountId int64, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error) {
	if amount <= 0 {
		err = fmt.Errorf("amount must be less than 0")
		return
	}
	account, err := s.accountDao.FindOneByID(ctx, accountId)
	if err != nil {
		return
	}
	if account.Status != model.EnableStatusEnabled {
		err = fmt.Errorf("account is disabled")
		return
	}
	if account.Balance < amount {
		err = fmt.Errorf("insufficient balance")
		return
	}
	return s.accountDao.DeductBalance(ctx, accountId, amount, requestId, typ, reason)
}

func (s *Service) AddBalance(ctx context.Context, accountId int64, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error) {
	if amount <= 0 {
		err = fmt.Errorf("amount must be less than 0")
		return
	}
	account, err := s.accountDao.FindOneByID(ctx, accountId)
	if err != nil {
		return
	}
	if account.Status != model.EnableStatusEnabled {
		err = fmt.Errorf("account is disabled")
		return
	}
	return s.accountDao.IncreaseBalance(ctx, accountId, amount, requestId, typ, reason)
}

func (s *Service) GetAccountList(ctx context.Context, req *model.GetAccountListRequest) (total int64, list []*model.Account, err error) {
	f := &model.AccountFilter{
		IDs:    db.In(req.Ids, db.OmitIfZero[[]int64]()),
		Name:   db.Eq(req.Name, db.OmitIfZero[string]()),
		Status: db.Eq(req.Status, db.OmitIfZero[model.EnableStatus]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.accountDao.Count(ctx, f)
		if err != nil {
			return
		}
		if !db.HasRecrods(total, req.PageParam.Page, req.PageParam.PageSize) {
			return
		}
		options = append(options,
			db.WithPaging(req.PageParam.Page, req.PageParam.PageSize),
			db.WithOrder(req.PageParam.OrderBy, nil))
	}
	list, err = s.accountDao.Find(ctx, f, options...)
	return
}

func (s *Service) CreateAccount(ctx context.Context, req *model.CreateAccountRequest) (info *model.Account, err error) {
	info, err = s.accountDao.FindOne(ctx, &model.AccountFilter{Name: db.Eq(req.Account.Name)})
	if db.IsDbError(err) {
		return
	}
	if info != nil {
		err = fmt.Errorf("account already exists, name: %s", req.Account.Name)
		return
	}
	info = &model.Account{
		Name:     req.Account.Name,
		Nickname: req.Account.Nickname,
		Balance:  req.Account.Balance,
		Status:   model.EnableStatus(req.Account.Status),
	}
	err = s.accountDao.Create(ctx, info)
	return
}

func (s *Service) UpdateAccount(ctx context.Context, req *model.UpdateAccountRequest) (info *model.Account, err error) {
	info, err = s.accountDao.FindOneByID(ctx, req.Account.Id)
	if db.IsDbError(err) {
		return
	}
	if info == nil {
		err = fmt.Errorf("account not found, id: %d", req.Account.Id)
		return
	}
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "name") {
		update["name"] = req.Account.Name
	}
	if lo.Contains(req.UpdateMask, "nickname") {
		update["nickname"] = req.Account.Nickname
	}
	if lo.Contains(req.UpdateMask, "balance") {
		update["balance"] = req.Account.Balance
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.Account.Status
	}
	if len(update) == 0 {
		err = fmt.Errorf("no fields to update")
		return
	}
	err = s.accountDao.UpdateOne(ctx, info, update)
	return
}

func (s *Service) DeleteAccounts(ctx context.Context, req *model.DeleteAccountsRequest) (err error) {
	_, err = s.accountDao.Delete(ctx, &model.AccountFilter{IDs: db.In(req.Ids)})
	return
}
