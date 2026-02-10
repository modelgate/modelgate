package service

import (
	"context"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) CreateLedger(ctx context.Context, req *model.CreateLedgerRequest) (info *model.Ledger, err error) {
	info = &model.Ledger{
		AccountId:    req.Ledger.AccountId,
		Type:         model.LedgerType(req.Ledger.Type),
		Amount:       req.Ledger.Amount,
		BalanceAfter: req.Ledger.BalanceAfter,
		RequestId:    req.Ledger.RequestId,
		Reason:       req.Ledger.Reason,
	}
	err = s.ledgerDao.Create(ctx, info)
	return
}

func (s *Service) DeleteLedgers(ctx context.Context, req *model.DeleteLedgersRequest) (err error) {
	_, err = s.ledgerDao.Delete(ctx, &model.LedgerFilter{IDs: db.In(req.Ids)})
	return
}

func (s *Service) GetLedgerList(ctx context.Context, req *model.GetLedgerListRequest) (total int64, list []*model.Ledger, err error) {
	f := &model.LedgerFilter{
		AccountId: db.Eq(req.AccountId, db.OmitIfZero[int64]()),
		Type:      db.Eq(req.Type, db.OmitIfZero[model.LedgerType]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.ledgerDao.Count(ctx, f)
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
	list, err = s.ledgerDao.Find(ctx, f, options...)
	return
}
