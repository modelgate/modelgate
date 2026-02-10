//go:generate go run go.uber.org/mock/mockgen -source $GOFILE -destination ./service_mock.go -package $GOPACKAGE
package relay

import (
	"context"

	"github.com/modelgate/modelgate/internal/relay/model"
)

type Service interface {
	StartWorker(ctx context.Context)

	CreateAccountApiKey(ctx context.Context, req *model.CreateAccountApiKeyRequest) (*model.AccountApiKey, error)
	UpdateAccountApiKey(ctx context.Context, req *model.UpdateAccountApiKeyRequest) (*model.AccountApiKey, error)
	DeleteAccountApiKeys(ctx context.Context, req *model.DeleteAccountApiKeysRequest) error
	GetAccountApiKeyList(ctx context.Context, req *model.GetAccountApiKeyListRequest) (int64, []*model.AccountApiKey, error)
	GetAccountApiKey(ctx context.Context, apiKey string) (*model.AccountApiKey, error)

	CreateProvider(ctx context.Context, req *model.CreateProviderRequest) (*model.Provider, error)
	UpdateProvider(ctx context.Context, req *model.UpdateProviderRequest) (*model.Provider, error)
	DeleteProviders(ctx context.Context, req *model.DeleteProvidersRequest) error
	GetProviderList(ctx context.Context, req *model.GetProviderListRequest) (int64, []*model.Provider, error)

	CreateProviderApiKey(ctx context.Context, req *model.CreateProviderApiKeyRequest) (*model.ProviderApiKey, error)
	UpdateProviderApiKey(ctx context.Context, req *model.UpdateProviderApiKeyRequest) (*model.ProviderApiKey, error)
	DeleteProviderApiKeys(ctx context.Context, req *model.DeleteProviderApiKeysRequest) error
	GetProviderApiKeyList(ctx context.Context, req *model.GetProviderApiKeyListRequest) (int64, []*model.ProviderApiKey, error)

	CreateModel(ctx context.Context, req *model.CreateModelRequest) (*model.Model, error)
	UpdateModel(ctx context.Context, req *model.UpdateModelRequest) (*model.Model, error)
	DeleteModels(ctx context.Context, req *model.DeleteModelsRequest) error
	GetModelList(ctx context.Context, req *model.GetModelListRequest) (int64, []*model.Model, error)
	ResolveModel(ctx context.Context, provider string, modelCode string) (info *model.ResolvedModel, err error)

	CreateModelPricing(ctx context.Context, req *model.CreateModelPricingRequest) (*model.ModelPricing, error)
	UpdateModelPricing(ctx context.Context, req *model.UpdateModelPricingRequest) (*model.ModelPricing, error)
	DeleteModelPricings(ctx context.Context, req *model.DeleteModelPricingsRequest) error
	GetModelPricingList(ctx context.Context, req *model.GetModelPricingListRequest) (int64, []*model.ModelPricing, error)

	CreateLedger(ctx context.Context, req *model.CreateLedgerRequest) (*model.Ledger, error)
	DeleteLedgers(ctx context.Context, req *model.DeleteLedgersRequest) error
	GetLedgerList(ctx context.Context, req *model.GetLedgerListRequest) (int64, []*model.Ledger, error)

	DeductBalance(ctx context.Context, accountId int64, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error)
	AddBalance(ctx context.Context, accountId int64, amount int64, requestId int64, typ model.LedgerType, reason string) (ledger *model.Ledger, err error)

	CreateAccount(ctx context.Context, req *model.CreateAccountRequest) (*model.Account, error)
	UpdateAccount(ctx context.Context, req *model.UpdateAccountRequest) (*model.Account, error)
	DeleteAccounts(ctx context.Context, req *model.DeleteAccountsRequest) error
	GetAccountList(ctx context.Context, req *model.GetAccountListRequest) (int64, []*model.Account, error)

	CreateRequest(ctx context.Context, req *model.CreateRequestRequest) (*model.Request, error)
	UpdateRequestCompleted(ctx context.Context, req *model.UpdateRequestCompletedRequest) error
	DeleteRequests(ctx context.Context, req *model.DeleteRequestsRequest) error
	GetRequestList(ctx context.Context, req *model.GetRequestListRequest) (int64, []*model.Request, error)

	GetRelayHourlyUsageList(ctx context.Context, req *model.GetRelayHourlyUsageListRequest) (int64, []*model.RelayHourlyUsage, error)

	AddRequestUsage(ctx context.Context, providerCode string, metric string, value int64) error
	AddPointUsage(ctx context.Context, providerCode string, providerApiKeyId, accountApiKeyId int64, value int64) error
	GetRelayInfo(ctx context.Context) (*model.RelayInfo, error)
	GetRelayUsage(ctx context.Context) (*model.RelayUsage, error)
}
