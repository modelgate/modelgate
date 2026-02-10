package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/dao"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/utils"
)

type Service struct {
	requestDao          relay.RequestDAO
	requestAttemptDao   relay.RequestAttemptDAO
	providerDao         relay.ProviderDAO
	providerApiKeyDao   relay.ProviderApiKeyDAO
	accountApiKeyDao    relay.AccountApiKeyDAO
	modelPricingDao     relay.ModelPricingDAO
	modelDao            relay.ModelDAO
	accountDao          relay.AccountDAO
	ledgerDao           relay.LedgerDAO
	relayUsageDao       relay.RelayUsageDAO
	relayHourlyUsageDao relay.RelayHourlyUsageDAO
	redisClient         *redis.Client
}

func New(i do.Injector) (relay.Service, error) {
	return &Service{
		requestDao:          do.MustInvoke[relay.RequestDAO](i),
		requestAttemptDao:   do.MustInvoke[relay.RequestAttemptDAO](i),
		providerDao:         do.MustInvoke[relay.ProviderDAO](i),
		providerApiKeyDao:   do.MustInvoke[relay.ProviderApiKeyDAO](i),
		accountApiKeyDao:    do.MustInvoke[relay.AccountApiKeyDAO](i),
		modelPricingDao:     do.MustInvoke[relay.ModelPricingDAO](i),
		modelDao:            do.MustInvoke[relay.ModelDAO](i),
		accountDao:          do.MustInvoke[relay.AccountDAO](i),
		ledgerDao:           do.MustInvoke[relay.LedgerDAO](i),
		relayUsageDao:       do.MustInvoke[relay.RelayUsageDAO](i),
		relayHourlyUsageDao: do.MustInvoke[relay.RelayHourlyUsageDAO](i),
		redisClient:         do.MustInvoke[*redis.Client](i),
	}, nil
}

// Init
func Init(i do.Injector) {
	dao.Init(i)
	do.Provide(i, New)

	// 启动多实例时，只让一个实例执行统计任务
	instanceID := uuid.New().String()

	s := do.MustInvoke[relay.Service](i)
	redisClient := do.MustInvoke[*redis.Client](i)
	leader := utils.NewLeaderElector(redisClient, instanceID, model.WorkerKeyLeader, time.Minute)
	leader.Run(context.Background(), func(leaderCtx context.Context) {
		s.StartWorker(leaderCtx)
	})
}
