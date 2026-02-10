package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) CreateModelPricing(ctx context.Context, req *model.CreateModelPricingRequest) (info *model.ModelPricing, err error) {
	info = &model.ModelPricing{
		ProviderCode:      req.ModelPricing.ProviderCode,
		ModelCode:         req.ModelPricing.ModelCode,
		Currency:          model.Currency(req.ModelPricing.Currency),
		PointsPerCurrency: req.ModelPricing.PointsPerCurrency,
		TokenNum:          req.ModelPricing.TokenNum,
		InputPrice:        float64(req.ModelPricing.InputPrice),
		InputCachePrice:   float64(req.ModelPricing.InputCachePrice),
		OutputPrice:       float64(req.ModelPricing.OutputPrice),
		Status:            model.EnableStatus(req.ModelPricing.Status),
		EffectiveFrom:     req.ModelPricing.EffectiveFrom.AsTime(),
		EffectiveTo:       req.ModelPricing.EffectiveTo.AsTime(),
	}
	err = s.modelPricingDao.Create(ctx, info)
	return
}

func (s *Service) UpdateModelPricing(ctx context.Context, req *model.UpdateModelPricingRequest) (info *model.ModelPricing, err error) {
	info, err = s.modelPricingDao.FindOneByID(ctx, req.ModelPricing.Id)
	if err != nil {
		return nil, err
	}
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "provider_code") {
		update["provider_code"] = req.ModelPricing.ProviderCode
	}
	if lo.Contains(req.UpdateMask, "model_code") {
		update["model_code"] = req.ModelPricing.ModelCode
	}
	if lo.Contains(req.UpdateMask, "currency") {
		update["currency"] = req.ModelPricing.Currency
	}
	if lo.Contains(req.UpdateMask, "points_per_currency") {
		update["points_per_currency"] = req.ModelPricing.PointsPerCurrency
	}
	if lo.Contains(req.UpdateMask, "token_num") {
		update["token_num"] = req.ModelPricing.TokenNum
	}
	if lo.Contains(req.UpdateMask, "input_price") {
		update["input_price"] = req.ModelPricing.InputPrice
	}
	if lo.Contains(req.UpdateMask, "input_cache_price") {
		update["input_cache_price"] = req.ModelPricing.InputCachePrice
	}
	if lo.Contains(req.UpdateMask, "output_price") {
		update["output_price"] = req.ModelPricing.OutputPrice
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.ModelPricing.Status
	}
	if lo.Contains(req.UpdateMask, "effective_from") {
		update["effective_from"] = req.ModelPricing.EffectiveFrom.AsTime()
	}
	if lo.Contains(req.UpdateMask, "effective_to") {
		update["effective_to"] = req.ModelPricing.EffectiveTo.AsTime()
	}
	if len(update) == 0 {
		err = fmt.Errorf("no fields to update")
		return
	}
	err = s.modelPricingDao.UpdateOne(ctx, info, update)
	return
}

func (s *Service) DeleteModelPricings(ctx context.Context, req *model.DeleteModelPricingsRequest) error {
	_, err := s.modelPricingDao.Delete(ctx, &model.ModelPricingFilter{IDs: db.In(req.Ids)})
	return err
}

func (s *Service) GetModelPricingList(ctx context.Context, req *model.GetModelPricingListRequest) (total int64, list []*model.ModelPricing, err error) {
	f := &model.ModelPricingFilter{
		ProviderCode: db.Like(req.ProviderCode+"%", db.OmitIf(func(s string) bool { return s == "%" })),
		ModelCode:    db.Like(req.ModelCode+"%", db.OmitIf(func(s string) bool { return s == "%" })),
		Currency:     db.Eq(req.Currency, db.OmitIfZero[model.Currency]()),
		Status:       db.Eq(req.Status, db.OmitIfZero[model.EnableStatus]()),
		// EffectiveFrom: db.Gte(req.EffectiveFrom, db.OmitIfZero[time.Time]()),
		// EffectiveTo:   db.Lte(req.EffectiveTo, db.OmitIfZero[time.Time]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.modelPricingDao.Count(ctx, f)
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
	list, err = s.modelPricingDao.Find(ctx, f, options...)
	return
}
