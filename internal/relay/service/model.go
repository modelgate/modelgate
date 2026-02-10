package service

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
	"github.com/modelgate/modelgate/pkg/utils"
)

func (s *Service) CreateModel(ctx context.Context, req *model.CreateModelRequest) (info *model.Model, err error) {
	provider, err := s.providerDao.FindOneByID(ctx, req.Model.ProviderId)
	if err != nil {
		return
	}
	if provider.Status != model.EnableStatusEnabled {
		err = fmt.Errorf("provider not enabled")
		return
	}
	info = &model.Model{
		ProviderId:   provider.ID,
		ProviderCode: provider.Code,
		Name:         req.Model.Name,
		Code:         req.Model.Code,
		ActualCode:   req.Model.ActualCode,
		Priority:     int(req.Model.Priority),
		Weight:       int(req.Model.Weight),
		Status:       model.ModelStatus(req.Model.Status),
	}
	err = s.modelDao.Create(ctx, info)
	return
}

func (s *Service) UpdateModel(ctx context.Context, req *model.UpdateModelRequest) (info *model.Model, err error) {
	info, err = s.modelDao.FindOneByID(ctx, req.Model.Id)
	if err != nil {
		return
	}
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "provider_id") {
		var provider *model.Provider
		provider, err = s.providerDao.FindOneByID(ctx, req.Model.ProviderId)
		if err != nil {
			return
		}
		if provider.Status != model.EnableStatusEnabled {
			err = fmt.Errorf("provider not enabled")
			return
		}
		update["provider_id"] = provider.ID
		update["provider_code"] = provider.Code
	}
	if lo.Contains(req.UpdateMask, "name") {
		update["name"] = req.Model.Name
	}
	if lo.Contains(req.UpdateMask, "code") {
		update["code"] = req.Model.Code
	}
	if lo.Contains(req.UpdateMask, "actual_code") {
		update["actual_code"] = req.Model.ActualCode
	}
	if lo.Contains(req.UpdateMask, "priority") {
		update["priority"] = req.Model.Priority
	}
	if lo.Contains(req.UpdateMask, "weight") {
		update["weight"] = req.Model.Weight
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.Model.Status
	}
	if len(update) == 0 {
		err = fmt.Errorf("no fields to update")
		return
	}
	err = s.modelDao.UpdateOne(ctx, info, update)
	return
}

func (s *Service) DeleteModels(ctx context.Context, req *model.DeleteModelsRequest) (err error) {
	_, err = s.modelDao.Delete(ctx, &model.ModelFilter{IDs: db.In(req.Ids)})
	return
}

func (s *Service) GetModelList(ctx context.Context, req *model.GetModelListRequest) (total int64, list []*model.Model, err error) {
	f := &model.ModelFilter{
		Name:         db.Like(req.Name+"%", db.OmitIf(func(s string) bool { return s == "%" })),
		Code:         db.Eq(req.Code, db.OmitIfZero[string]()),
		ProviderCode: db.Eq(req.ProviderCode, db.OmitIfZero[string]()),
		Status:       db.Eq(req.Status, db.OmitIfZero[model.ModelStatus]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.modelDao.Count(ctx, f)
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
	list, err = s.modelDao.Find(ctx, f, options...)
	return
}

// ResolveModel 解析模型
func (s *Service) ResolveModel(ctx context.Context, provider string, modelCode string) (info *model.ResolvedModel, err error) {
	modelInfo, err := s.pickModel(ctx, provider, modelCode)
	if err != nil {
		return
	}
	log.Infof("picked model, provide: %s, model, %s", modelInfo.ProviderCode, modelInfo.Code)
	providerInfo, err := s.providerDao.FindOneByID(ctx, modelInfo.ProviderId)
	if err != nil {
		return
	}
	if providerInfo.Status != model.EnableStatusEnabled {
		err = fmt.Errorf("provider not enabled")
		return
	}
	keyInfo, err := s.pickProviderApiKey(ctx, modelInfo.ProviderId)
	if err != nil {
		return
	}
	actualCode := modelInfo.GetActualCode()
	f := &model.ModelPricingFilter{
		ProviderCode:  db.Eq(modelInfo.ProviderCode),
		ModelCode:     db.Eq(actualCode),
		EffectiveFrom: db.Lte(time.Now()),
		EffectiveTo:   db.Gt(time.Now()),
	}
	modelPrice, err := s.modelPricingDao.FindOne(ctx, f, db.WithOrder("-effective_from", nil))
	if err != nil {
		return
	}
	info = &model.ResolvedModel{
		// 模型
		ModelId:      modelInfo.ID,
		ModelCode:    actualCode,
		ProviderId:   modelInfo.ProviderId,
		ProviderCode: modelInfo.ProviderCode,
		// 供应商
		BaseUrl: providerInfo.BaseUrl,
		// 供应商ApiKey
		ApiKeyId:        keyInfo.ID,
		ApiKeyEncrypted: keyInfo.KeyEncrypted,
		// 价格
		InputPrice:        modelPrice.InputPrice,
		InputCachePrice:   modelPrice.InputCachePrice,
		OutputPrice:       modelPrice.OutputPrice,
		TokenNum:          modelPrice.TokenNum,
		PointsPerCurrency: modelPrice.PointsPerCurrency,
	}
	return
}

// pickModel 按照权重随机选择一个模型
func (s *Service) pickModel(ctx context.Context, providerCode string, modelCode string) (info *model.Model, err error) {
	f := &model.ModelFilter{
		ProviderCode: db.Eq(providerCode, db.OmitIfZero[string]()),
		Code:         db.Eq(modelCode, db.OmitIfZero[string]()),
		Status:       db.Eq(model.ModelStatusEnabled),
	}
	var list []*model.Model
	list, err = s.modelDao.Find(ctx, f, db.WithOrder("priority", nil))
	if err != nil {
		return
	}
	if len(list) == 0 {
		err = fmt.Errorf("model not found，provider: %s, model: %s", providerCode, modelCode)
		return
	}
	if len(list) == 1 {
		return list[0], nil
	}
	minPriority := list[0].Priority
	modelList := lo.FlatMap(list, func(info *model.Model, index int) []*model.Model {
		if info.Priority == minPriority {
			return []*model.Model{info}
		} else {
			return []*model.Model{}
		}
	})
	return utils.PickByWeight(modelList), nil
}

// pickProviderApiKey 按照权重随机选择一个API Key
func (s *Service) pickProviderApiKey(ctx context.Context, providerId int64) (keyInfo *model.ProviderApiKey, err error) {
	f := &model.ProviderApiKeyFilter{
		ProviderId: db.Eq(providerId),
		Status:     db.Eq(model.ApiKeyStatusEnabled),
	}
	var list []*model.ProviderApiKey
	list, err = s.providerApiKeyDao.Find(ctx, f)
	if err != nil {
		return
	}
	if len(list) == 0 {
		err = fmt.Errorf("provider api key not found，provider: %d", providerId)
		return
	}
	return utils.PickByWeight(list), nil
}
