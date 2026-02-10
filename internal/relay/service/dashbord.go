package service

import (
	"context"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

func (s *Service) GetRelayInfo(ctx context.Context) (*model.RelayInfo, error) {
	providerCount, err := s.providerDao.Count(ctx, &model.ProviderFilter{Status: db.Eq(model.EnableStatusEnabled)})
	if err != nil {
		return nil, err
	}
	modelCount, err := s.modelDao.Count(ctx, &model.ModelFilter{Status: db.Eq(model.ModelStatusEnabled)})
	if err != nil {
		return nil, err
	}
	apiKeyCount, err := s.providerApiKeyDao.Count(ctx, &model.ProviderApiKeyFilter{Status: db.Eq(model.ApiKeyStatusEnabled)})
	if err != nil {
		return nil, err
	}
	return &model.RelayInfo{
		ProviderCount: providerCount,
		ModelCount:    modelCount,
		ApiKeyCount:   apiKeyCount,
	}, nil
}
