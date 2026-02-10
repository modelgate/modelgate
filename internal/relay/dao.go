//go:generate go run go.uber.org/mock/mockgen -source $GOFILE -destination ./dao_mock.go -package $GOPACKAGE
package relay

import (
	"context"

	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/pkg/db"
)

type RequestDAO interface {
	Create(ctx context.Context, m *model.Request) error
	Save(ctx context.Context, m *model.Request) error
	Update(ctx context.Context, filter *model.RequestFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Request, update map[string]any) error
	Count(ctx context.Context, f *model.RequestFilter) (total int64, err error)
	Find(ctx context.Context, f *model.RequestFilter, opts ...db.Option) (ms []*model.Request, err error)
	FindOne(ctx context.Context, f *model.RequestFilter, opts ...db.Option) (*model.Request, error)
	FindOneByID(ctx context.Context, id int64) (m *model.Request, err error)
	Delete(ctx context.Context, filter *model.RequestFilter) (int64, error)
}

type RequestAttemptDAO interface {
	Create(ctx context.Context, m *model.RequestAttempt) error
	Save(ctx context.Context, m *model.RequestAttempt) error
	Update(ctx context.Context, filter *model.RequestAttemptFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.RequestAttempt, update map[string]any) error
	Count(ctx context.Context, f *model.RequestAttemptFilter) (total int64, err error)
	Find(ctx context.Context, f *model.RequestAttemptFilter, opts ...db.Option) (ms []*model.RequestAttempt, err error)
	FindOne(ctx context.Context, f *model.RequestAttemptFilter, opts ...db.Option) (*model.RequestAttempt, error)
	FindOneByID(ctx context.Context, id int64) (m *model.RequestAttempt, err error)
	Delete(ctx context.Context, filter *model.RequestAttemptFilter) (int64, error)
}

type ProviderDAO interface {
	Create(ctx context.Context, m *model.Provider) error
	Save(ctx context.Context, m *model.Provider) error
	Update(ctx context.Context, filter *model.ProviderFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Provider, update map[string]any) error
	Count(ctx context.Context, f *model.ProviderFilter) (total int64, err error)
	Find(ctx context.Context, f *model.ProviderFilter, opts ...db.Option) (ms []*model.Provider, err error)
	FindOne(ctx context.Context, f *model.ProviderFilter, opts ...db.Option) (*model.Provider, error)
	FindOneByID(ctx context.Context, id int64) (m *model.Provider, err error)
	Delete(ctx context.Context, filter *model.ProviderFilter) (int64, error)
}

type ProviderApiKeyDAO interface {
	Create(ctx context.Context, m *model.ProviderApiKey) error
	Save(ctx context.Context, m *model.ProviderApiKey) error
	Update(ctx context.Context, filter *model.ProviderApiKeyFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.ProviderApiKey, update map[string]any) error
	Count(ctx context.Context, f *model.ProviderApiKeyFilter) (total int64, err error)
	Find(ctx context.Context, f *model.ProviderApiKeyFilter, opts ...db.Option) (ms []*model.ProviderApiKey, err error)
	FindOne(ctx context.Context, f *model.ProviderApiKeyFilter, opts ...db.Option) (*model.ProviderApiKey, error)
	FindOneByID(ctx context.Context, id int64) (m *model.ProviderApiKey, err error)
	Delete(ctx context.Context, filter *model.ProviderApiKeyFilter) (int64, error)
}

type AccountApiKeyDAO interface {
	Create(ctx context.Context, m *model.AccountApiKey) error
	Save(ctx context.Context, m *model.AccountApiKey) error
	Update(ctx context.Context, filter *model.AccountApiKeyFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.AccountApiKey, update map[string]any) error
	Count(ctx context.Context, f *model.AccountApiKeyFilter) (total int64, err error)
	Find(ctx context.Context, f *model.AccountApiKeyFilter, opts ...db.Option) (ms []*model.AccountApiKey, err error)
	FindOne(ctx context.Context, f *model.AccountApiKeyFilter, opts ...db.Option) (*model.AccountApiKey, error)
	FindOneByID(ctx context.Context, id int64) (m *model.AccountApiKey, err error)
	Delete(ctx context.Context, filter *model.AccountApiKeyFilter) (int64, error)
}

type ModelPricingDAO interface {
	Create(ctx context.Context, m *model.ModelPricing) error
	Save(ctx context.Context, m *model.ModelPricing) error
	Update(ctx context.Context, filter *model.ModelPricingFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.ModelPricing, update map[string]any) error
	Count(ctx context.Context, f *model.ModelPricingFilter) (total int64, err error)
	Find(ctx context.Context, f *model.ModelPricingFilter, opts ...db.Option) (ms []*model.ModelPricing, err error)
	FindOne(ctx context.Context, f *model.ModelPricingFilter, opts ...db.Option) (*model.ModelPricing, error)
	FindOneByID(ctx context.Context, id int64) (m *model.ModelPricing, err error)
	Delete(ctx context.Context, filter *model.ModelPricingFilter) (int64, error)
}

type ModelDAO interface {
	Create(ctx context.Context, m *model.Model) error
	Save(ctx context.Context, m *model.Model) error
	Update(ctx context.Context, filter *model.ModelFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Model, update map[string]any) error
	Count(ctx context.Context, f *model.ModelFilter) (total int64, err error)
	Find(ctx context.Context, f *model.ModelFilter, opts ...db.Option) (ms []*model.Model, err error)
	FindOne(ctx context.Context, f *model.ModelFilter, opts ...db.Option) (*model.Model, error)
	FindOneByID(ctx context.Context, id int64) (m *model.Model, err error)
	Delete(ctx context.Context, filter *model.ModelFilter) (int64, error)
}

type AccountDAO interface {
	Create(ctx context.Context, m *model.Account) error
	Save(ctx context.Context, m *model.Account) error
	Update(ctx context.Context, filter *model.AccountFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Account, update map[string]any) error
	Count(ctx context.Context, f *model.AccountFilter) (total int64, err error)
	Find(ctx context.Context, f *model.AccountFilter, opts ...db.Option) (ms []*model.Account, err error)
	FindOne(ctx context.Context, f *model.AccountFilter, opts ...db.Option) (*model.Account, error)
	FindOneByID(ctx context.Context, id int64) (m *model.Account, err error)
	Delete(ctx context.Context, filter *model.AccountFilter) (int64, error)

	DeductBalance(ctx context.Context, accountId, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error)
	IncreaseBalance(ctx context.Context, accountId, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error)
}

type LedgerDAO interface {
	Create(ctx context.Context, m *model.Ledger) error
	Save(ctx context.Context, m *model.Ledger) error
	Update(ctx context.Context, filter *model.LedgerFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.Ledger, update map[string]any) error
	Count(ctx context.Context, f *model.LedgerFilter) (total int64, err error)
	Find(ctx context.Context, f *model.LedgerFilter, opts ...db.Option) (ms []*model.Ledger, err error)
	FindOne(ctx context.Context, f *model.LedgerFilter, opts ...db.Option) (*model.Ledger, error)
	FindOneByID(ctx context.Context, id int64) (m *model.Ledger, err error)
	Delete(ctx context.Context, filter *model.LedgerFilter) (int64, error)
}

type RelayUsageDAO interface {
	Create(ctx context.Context, m *model.RelayUsage) error
	Save(ctx context.Context, m *model.RelayUsage) error
	Update(ctx context.Context, filter *model.RelayUsageFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.RelayUsage, update map[string]any) error
	Count(ctx context.Context, f *model.RelayUsageFilter) (total int64, err error)
	Find(ctx context.Context, f *model.RelayUsageFilter, opts ...db.Option) (ms []*model.RelayUsage, err error)
	FindOne(ctx context.Context, f *model.RelayUsageFilter, opts ...db.Option) (*model.RelayUsage, error)
	FindOneByID(ctx context.Context, id int64) (m *model.RelayUsage, err error)
	Delete(ctx context.Context, filter *model.RelayUsageFilter) (int64, error)
}

type RelayHourlyUsageDAO interface {
	Create(ctx context.Context, m *model.RelayHourlyUsage) error
	Save(ctx context.Context, m *model.RelayHourlyUsage) error
	Update(ctx context.Context, filter *model.RelayHourlyUsageFilter, update map[string]any) (int64, error)
	UpdateOne(ctx context.Context, m *model.RelayHourlyUsage, update map[string]any) error
	Count(ctx context.Context, f *model.RelayHourlyUsageFilter) (total int64, err error)
	Find(ctx context.Context, f *model.RelayHourlyUsageFilter, opts ...db.Option) (ms []*model.RelayHourlyUsage, err error)
	FindOne(ctx context.Context, f *model.RelayHourlyUsageFilter, opts ...db.Option) (*model.RelayHourlyUsage, error)
	FindOneByID(ctx context.Context, id int64) (m *model.RelayHourlyUsage, err error)
	Delete(ctx context.Context, filter *model.RelayHourlyUsageFilter) (int64, error)
}
