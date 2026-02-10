package dao

import (
	"github.com/samber/do/v2"
)

// Init 注册Dao
func Init(i do.Injector) {
	do.Provide(i, NewRequestDao)
	do.Provide(i, NewRequestAttemptDao)
	do.Provide(i, NewProviderDao)
	do.Provide(i, NewProviderApiKeyDao)
	do.Provide(i, NewAccountApiKeyDao)
	do.Provide(i, NewModelPricingDao)
	do.Provide(i, NewModelDao)
	do.Provide(i, NewAccountDao)
	do.Provide(i, NewLedgerDao)
	do.Provide(i, NewRelayHourlyUsageDao)
	do.Provide(i, NewRelayUsageDao)
}
