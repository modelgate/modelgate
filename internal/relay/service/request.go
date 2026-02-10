package service

import (
	"context"
	"time"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) CreateRequest(ctx context.Context, req *model.CreateRequestRequest) (m *model.Request, err error) {
	request := model.Request{
		RequestUUID:      req.RequestUUID,
		AccountId:        req.AccountId,
		AccountApiKeyId:  req.AccountApiKeyId,
		ProviderId:       req.ProviderId,
		ProviderApiKeyId: req.ProviderApiKeyId,
		ModelId:          req.ModelId,
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
		Status:           model.RequestStatusPending,
	}
	err = s.requestDao.Create(ctx, &request)
	if err != nil {
		return
	}
	requestAttempt := model.RequestAttempt{
		RequestUUID:      req.RequestUUID,
		AttemptNo:        req.AttemptNo,
		AccountId:        req.AccountId,
		AccountApiKeyId:  req.AccountApiKeyId,
		ProviderId:       req.ProviderId,
		ProviderApiKeyId: req.ProviderApiKeyId,
		ModelId:          req.ModelId,
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
		Status:           model.RequestStatusPending,
	}
	err = s.requestAttemptDao.Create(ctx, &requestAttempt)
	if err != nil {
		return
	}
	// 统计
	s.AddRequestUsage(ctx, req.ProviderCode, model.MetricTotal, 1)
	return
}

func (s *Service) UpdateRequestCompleted(ctx context.Context, req *model.UpdateRequestCompletedRequest) (err error) {
	update := map[string]any{
		"prompt_tokens":     req.PromptTokens,
		"completion_tokens": req.CompletionTokens,
		"total_tokens":      req.TotalTokens,
		"status":            req.Status,
		"error_code":        req.ErrorCode,
		"error_message":     req.ErrorMessage,
		"completed_at":      time.Now(),
	}
	if req.ActualModel != "" {
		update["actual_model"] = req.ActualModel
	}
	_, err = s.requestDao.Update(ctx, &model.RequestFilter{RequestUUID: db.Eq(req.RequestUUID)}, update)
	if err != nil {
		return
	}
	_, err = s.requestAttemptDao.Update(ctx, &model.RequestAttemptFilter{RequestUUID: db.Eq(req.RequestUUID), AttemptNo: db.Eq(req.AttemptNo)}, update)
	if err != nil {
		return
	}
	// 统计厂商请求成功、失败数量
	if req.Status == model.RequestStatusSuccess {
		s.AddRequestUsage(ctx, req.ProviderCode, model.MetricSuccess, 1)
	} else {
		s.AddRequestUsage(ctx, req.ProviderCode, model.MetricFailed, 1)
	}
	return
}

func (s *Service) DeleteRequests(ctx context.Context, req *model.DeleteRequestsRequest) (err error) {
	_, err = s.requestDao.Delete(ctx, &model.RequestFilter{IDs: db.In(req.Ids)})
	return
}

func (s *Service) GetRequestList(ctx context.Context, req *model.GetRequestListRequest) (total int64, list []*model.Request, err error) {
	f := &model.RequestFilter{
		ProviderCode: db.Eq(req.ProviderCode, db.OmitIfZero[string]()),
		AccountId:    db.Eq(req.AccountId, db.OmitIfZero[int64]()),
		Object:       db.Eq(req.Object, db.OmitIfZero[string]()),
		Status:       db.Eq(req.Status, db.OmitIfZero[model.RequestStatus]()),
	}

	var options []db.Option
	if req.PageParam != nil {
		total, err = s.requestDao.Count(ctx, f)
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
	list, err = s.requestDao.Find(ctx, f, options...)
	return
}
