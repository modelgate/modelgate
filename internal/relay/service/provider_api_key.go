package service

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
	"github.com/modelgate/modelgate/pkg/utils"
)

func (s *Service) CreateProviderApiKey(ctx context.Context, req *model.CreateProviderApiKeyRequest) (info *model.ProviderApiKey, err error) {
	info, err = s.providerApiKeyDao.FindOneByID(ctx, req.ProviderApiKey.Id)
	if db.IsDbError(err) {
		return
	}
	if info != nil {
		err = fmt.Errorf("provider api key already exists, id: %d", req.ProviderApiKey.Id)
		return
	}
	provider, err := s.providerDao.FindOneByID(ctx, req.ProviderApiKey.ProviderId)
	if db.IsDbError(err) {
		return
	}
	if provider == nil {
		err = fmt.Errorf("provider not found, id: %d", req.ProviderApiKey.ProviderId)
		return
	}
	keyEncrypted, keyPrefix, keySuffix, err := s.encryptKey(req.ProviderApiKey.Key)
	if err != nil {
		return
	}
	info = &model.ProviderApiKey{
		ProviderId:   provider.ID,
		ProviderCode: provider.Code,
		Name:         req.ProviderApiKey.Name,
		KeyPrefix:    keyPrefix,
		KeySuffix:    keySuffix,
		KeyEncrypted: keyEncrypted,
		Weight:       int(req.ProviderApiKey.Weight),
		Status:       model.ApiKeyStatus(req.ProviderApiKey.Status),
	}
	err = s.providerApiKeyDao.Create(ctx, info)
	return
}

func (s *Service) encryptKey(key string) (keyEncrypted, prefix, suffix string, err error) {
	if len(key) < 20 {
		err = fmt.Errorf("key length must be greater than 16")
		return
	}
	prefix = key[:12]
	suffix = key[len(key)-4:]
	keyEncrypted, err = utils.EncryptAESGCM([]byte(key), []byte(config.GetConfig().Secret.Key))
	return
}

func (s *Service) UpdateProviderApiKey(ctx context.Context, req *model.UpdateProviderApiKeyRequest) (info *model.ProviderApiKey, err error) {
	info, err = s.providerApiKeyDao.FindOneByID(ctx, req.ProviderApiKey.Id)
	if db.IsDbError(err) {
		return
	}
	if info == nil {
		err = fmt.Errorf("provider api key not found, id: %d", req.ProviderApiKey.Id)
		return
	}
	update := make(map[string]any)
	if lo.Contains(req.UpdateMask, "provider_id") {
		var provider *model.Provider
		provider, err = s.providerDao.FindOneByID(ctx, req.ProviderApiKey.ProviderId)
		if err != nil {
			return
		}
		update["provider_id"] = provider.ID
		update["provider_code"] = provider.Code
	}
	if lo.Contains(req.UpdateMask, "name") {
		update["name"] = req.ProviderApiKey.Name
	}
	if lo.Contains(req.UpdateMask, "key") {
		var keyEncrypted, keyPrefix, keySuffix string
		keyEncrypted, keyPrefix, keySuffix, err = s.encryptKey(req.ProviderApiKey.Key)
		if err != nil {
			return
		}
		update["key_prefix"] = keyPrefix
		update["key_suffix"] = keySuffix
		update["key_encrypted"] = keyEncrypted
	}
	if lo.Contains(req.UpdateMask, "weight") {
		update["weight"] = req.ProviderApiKey.Weight
	}
	if lo.Contains(req.UpdateMask, "status") {
		update["status"] = req.ProviderApiKey.Status
	}
	if len(update) == 0 {
		err = fmt.Errorf("no fields to update")
		return
	}
	err = s.providerApiKeyDao.UpdateOne(ctx, info, update)
	return
}

func (s *Service) DeleteProviderApiKeys(ctx context.Context, req *model.DeleteProviderApiKeysRequest) (err error) {
	_, err = s.providerApiKeyDao.Delete(ctx, &model.ProviderApiKeyFilter{IDs: db.In(req.Ids)})
	return
}

func (s *Service) GetProviderApiKeyList(ctx context.Context, req *model.GetProviderApiKeyListRequest) (total int64, list []*model.ProviderApiKey, err error) {
	f := &model.ProviderApiKeyFilter{
		ProviderId:   db.Eq(req.ProviderId, db.OmitIfZero[int64]()),
		ProviderCode: db.Eq(req.ProviderCode, db.OmitIfZero[string]()),
		Name:         db.Like(req.Name+"%", db.OmitIf(func(s string) bool { return s == "%" })),
		Status:       db.Eq(req.Status, db.OmitIfZero[model.ApiKeyStatus]()),
	}
	var options []db.Option
	if req.PageParam != nil {
		total, err = s.providerApiKeyDao.Count(ctx, f)
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
	list, err = s.providerApiKeyDao.Find(ctx, f, options...)
	return
}
