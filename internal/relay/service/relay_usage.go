package service

import (
	"context"

	"github.com/modelgate/modelgate/internal/relay/model"
)

func (s *Service) GetRelayUsage(ctx context.Context) (*model.RelayUsage, error) {
	return s.relayUsageDao.FindOneByID(ctx, 1)
}
