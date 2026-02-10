package service

import (
	"context"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) GetRelayHourlyUsageList(ctx context.Context, req *model.GetRelayHourlyUsageListRequest) (total int64, list []*model.RelayHourlyUsage, err error) {
	f := &model.RelayHourlyUsageFilter{}
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
		f.Time = db.Between(req.StartTime, req.EndTime)
	}

	var options []db.Option
	if req.PageParam != nil {
		total, err = s.relayHourlyUsageDao.Count(ctx, f)
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
	list, err = s.relayHourlyUsageDao.Find(ctx, f, options...)
	return
}
