package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
	"github.com/modelgate/modelgate/pkg/utils"
)

func (s *Service) CreateAccountApiKey(ctx context.Context, req *model.CreateAccountApiKeyRequest) (info *model.AccountApiKey, err error) {
	scope, err := json.Marshal(req.AccountApiKey.Scope)
	if err != nil {
		return
	}
	rawKey := utils.GenApiKey()
	prefix, suffix := utils.MaskApiKey(rawKey)
	info = &model.AccountApiKey{
		AccountId:  req.AccountApiKey.AccountId,
		KeyName:    req.AccountApiKey.KeyName,
		KeyPrefix:  prefix,
		KeySuffix:  suffix,
		KeyHash:    utils.Sha256Hex(rawKey),
		Key:        rawKey,
		Scope:      string(scope),
		QuoteLimit: lo.ToPtr(req.AccountApiKey.QuoteLimit),
		RateLimit:  lo.ToPtr(int(req.AccountApiKey.RateLimit)),
		ExpiredAt:  lo.ToPtr(req.AccountApiKey.ExpiredAt.AsTime()),
		Status:     lo.Ternary(req.AccountApiKey.Status != "", model.ApiKeyStatus(req.AccountApiKey.Status), model.ApiKeyStatusEnabled),
	}
	err = s.accountApiKeyDao.Create(ctx, info)
	return
}

func (s *Service) UpdateAccountApiKey(ctx context.Context, req *model.UpdateAccountApiKeyRequest) (info *model.AccountApiKey, err error) {
	info, err = s.accountApiKeyDao.FindOneByID(ctx, req.AccountApiKey.Id)
	if db.IsDbError(err) {
		return
	}
	if info == nil {
		err = fmt.Errorf("account api key not found, id: %d", req.AccountApiKey.Id)
		return
	}
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "key_name") {
		update["key_name"] = req.AccountApiKey.KeyName
	}
	if lo.Contains(req.UpdateMask, "scope") {
		update["scope"] = req.AccountApiKey.Scope
	}
	if lo.Contains(req.UpdateMask, "quote_limit") {
		update["quote_limit"] = req.AccountApiKey.QuoteLimit
	}
	if lo.Contains(req.UpdateMask, "rate_limit") {
		update["rate_limit"] = req.AccountApiKey.RateLimit
	}
	if lo.Contains(req.UpdateMask, "expired_at") {
		update["expired_at"] = req.AccountApiKey.ExpiredAt.AsTime()
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.AccountApiKey.Status
	}
	if len(update) == 0 {
		return
	}
	err = s.accountApiKeyDao.UpdateOne(ctx, info, update)
	return
}

func (s *Service) DeleteAccountApiKeys(ctx context.Context, req *model.DeleteAccountApiKeysRequest) (err error) {
	_, err = s.accountApiKeyDao.Delete(ctx, &model.AccountApiKeyFilter{IDs: db.In(req.Ids)})
	return
}

func (s *Service) GetAccountApiKeyList(ctx context.Context, req *model.GetAccountApiKeyListRequest) (total int64, list []*model.AccountApiKey, err error) {
	f := &model.AccountApiKeyFilter{
		AccountId: db.Eq(req.AccountId, db.OmitIfZero[int64]()),
		Status:    db.Eq(req.Status, db.OmitIfZero[model.ApiKeyStatus]()),
		Keyword:   req.Keyword,
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.accountApiKeyDao.Count(ctx, f)
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
	list, err = s.accountApiKeyDao.Find(ctx, f, options...)
	return
}

func (s *Service) GetAccountApiKey(ctx context.Context, apiKey string) (accountApiKey *model.AccountApiKey, err error) {
	if len(apiKey) < 16 {
		err = fmt.Errorf("invalid api key")
		return
	}

	keyHash := utils.Sha256Hex(apiKey)
	f := &model.AccountApiKeyFilter{
		KeyHash:   db.Eq(keyHash),
		ExpiredAt: db.GteOrNull(lo.ToPtr(time.Now())),
	}

	accountApiKey, err = s.accountApiKeyDao.FindOne(ctx, f)
	if err != nil {
		log.Errorf("find account api key error: %v", err)
		return
	}
	if accountApiKey.Status != model.ApiKeyStatusEnabled {
		err = errors.New("account api key is not enabled")
		return
	}
	return
}
