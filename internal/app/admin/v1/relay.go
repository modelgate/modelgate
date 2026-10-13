package v1

import (
	"context"
	"sort"
	"strings"
	"time"

	"connectrpc.com/connect/v2"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/internal/runtime/core"
	v1pb "github.com/modelgate/modelgate/pkg/proto/admin/v1"
	relaypb "github.com/modelgate/modelgate/pkg/proto/model/relay"
	"github.com/modelgate/modelgate/pkg/types"
	"github.com/modelgate/modelgate/pkg/utils"
)

type RelayService struct {
	v1pb.UnimplementedRelayServiceHandler
	relayService relay.Service
}

func NewRelayService(i do.Injector) (*RelayService, error) {
	return &RelayService{
		relayService: do.MustInvoke[relay.Service](i),
	}, nil
}

func (s *RelayService) CreateProvider(ctx context.Context, req *v1pb.CreateProviderRequest) (resp *relaypb.Provider, err error) {
	provider, err := s.relayService.CreateProvider(ctx, &model.CreateProviderRequest{Provider: req.Provider})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = provider.ToProto()
	return resp, nil
}

func (s *RelayService) UpdateProvider(ctx context.Context, req *v1pb.UpdateProviderRequest) (resp *relaypb.Provider, err error) {
	provider, err := s.relayService.UpdateProvider(ctx, &model.UpdateProviderRequest{
		Provider:   req.Provider,
		UpdateMask: req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = provider.ToProto()
	return resp, nil
}

func (s *RelayService) DeleteProviders(ctx context.Context, req *v1pb.DeleteProvidersRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteProviders(ctx, &model.DeleteProvidersRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetProviderList(ctx context.Context, req *v1pb.GetProviderListRequest) (resp *v1pb.GetProviderListResponse, err error) {
	total, list, err := s.relayService.GetProviderList(ctx, &model.GetProviderListRequest{
		PageParam: types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		Name:      strings.TrimSpace(req.Name),
		Code:      strings.TrimSpace(req.Code),
		Status:    model.EnableStatus(strings.TrimSpace(req.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetProviderListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.Provider, _ int) *relaypb.Provider {
			return item.ToProto()
		}),
	}
	return resp, nil
}

func (s *RelayService) GetProviderCodeList(ctx context.Context, req *v1pb.GetProviderCodeListRequest) (resp *v1pb.GetProviderCodeListResponse, err error) {
	resp = &v1pb.GetProviderCodeListResponse{Records: core.AllProviderCodeList}
	return resp, nil
}

func (s *RelayService) CreateModel(ctx context.Context, req *v1pb.CreateModelRequest) (resp *relaypb.Model, err error) {
	model, err := s.relayService.CreateModel(ctx, &model.CreateModelRequest{Model: req.Model})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = model.ToProto()
	return resp, nil
}

func (s *RelayService) UpdateModel(ctx context.Context, req *v1pb.UpdateModelRequest) (resp *relaypb.Model, err error) {
	model, err := s.relayService.UpdateModel(ctx, &model.UpdateModelRequest{
		Model:      req.Model,
		UpdateMask: req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = model.ToProto()
	return resp, nil
}

func (s *RelayService) DeleteModels(ctx context.Context, req *v1pb.DeleteModelsRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteModels(ctx, &model.DeleteModelsRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetModelList(ctx context.Context, req *v1pb.GetModelListRequest) (resp *v1pb.GetModelListResponse, err error) {
	total, list, err := s.relayService.GetModelList(ctx, &model.GetModelListRequest{
		PageParam:    types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		Name:         strings.TrimSpace(req.Name),
		Code:         strings.TrimSpace(req.Code),
		ProviderCode: strings.TrimSpace(req.ProviderCode),
		Status:       model.ModelStatus(strings.TrimSpace(req.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetModelListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.Model, _ int) *relaypb.Model {
			return item.ToProto()
		}),
	}
	return resp, nil
}

func (s *RelayService) CreateProviderApiKey(ctx context.Context, req *v1pb.CreateProviderApiKeyRequest) (resp *relaypb.ProviderApiKey, err error) {
	providerApiKey, err := s.relayService.CreateProviderApiKey(ctx, &model.CreateProviderApiKeyRequest{ProviderApiKey: req.ProviderApiKey})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = providerApiKey.ToProto()
	return resp, nil
}

func (s *RelayService) UpdateProviderApiKey(ctx context.Context, req *v1pb.UpdateProviderApiKeyRequest) (resp *relaypb.ProviderApiKey, err error) {
	providerApiKey, err := s.relayService.UpdateProviderApiKey(ctx, &model.UpdateProviderApiKeyRequest{
		ProviderApiKey: req.ProviderApiKey,
		UpdateMask:     req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = providerApiKey.ToProto()
	return resp, nil
}

func (s *RelayService) DeleteProviderApiKeys(ctx context.Context, req *v1pb.DeleteProviderApiKeysRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteProviderApiKeys(ctx, &model.DeleteProviderApiKeysRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetProviderApiKeyList(ctx context.Context, req *v1pb.GetProviderApiKeyListRequest) (resp *v1pb.GetProviderApiKeyListResponse, err error) {
	total, list, err := s.relayService.GetProviderApiKeyList(ctx, &model.GetProviderApiKeyListRequest{
		PageParam:    types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		ProviderId:   req.ProviderId,
		ProviderCode: strings.TrimSpace(req.ProviderCode),
		Name:         strings.TrimSpace(req.Name),
		Status:       model.ApiKeyStatus(strings.TrimSpace(req.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetProviderApiKeyListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.ProviderApiKey, _ int) *relaypb.ProviderApiKey {
			return item.ToProto()
		}),
	}
	return resp, nil
}

func (s *RelayService) CreateModelPricing(ctx context.Context, req *v1pb.CreateModelPricingRequest) (resp *relaypb.ModelPricing, err error) {
	modelPricing, err := s.relayService.CreateModelPricing(ctx, &model.CreateModelPricingRequest{ModelPricing: req.ModelPricing})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = modelPricing.ToProto()
	return resp, nil
}

func (s *RelayService) UpdateModelPricing(ctx context.Context, req *v1pb.UpdateModelPricingRequest) (resp *relaypb.ModelPricing, err error) {
	modelPricing, err := s.relayService.UpdateModelPricing(ctx, &model.UpdateModelPricingRequest{
		ModelPricing: req.ModelPricing,
		UpdateMask:   req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = modelPricing.ToProto()
	return resp, nil
}

func (s *RelayService) DeleteModelPricings(ctx context.Context, req *v1pb.DeleteModelPricingsRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteModelPricings(ctx, &model.DeleteModelPricingsRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetModelPricingList(ctx context.Context, req *v1pb.GetModelPricingListRequest) (resp *v1pb.GetModelPricingListResponse, err error) {
	total, list, err := s.relayService.GetModelPricingList(ctx, &model.GetModelPricingListRequest{
		PageParam:     types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		ProviderCode:  strings.TrimSpace(req.ProviderCode),
		ModelCode:     strings.TrimSpace(req.ModelCode),
		Currency:      model.Currency(strings.TrimSpace(req.Currency)),
		EffectiveFrom: req.EffectiveFrom.AsTime(),
		EffectiveTo:   req.EffectiveTo.AsTime(),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetModelPricingListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.ModelPricing, _ int) *relaypb.ModelPricing {
			return item.ToProto()
		}),
	}
	return resp, nil
}

func (s *RelayService) CreateLedger(ctx context.Context, req *v1pb.CreateLedgerRequest) (resp *relaypb.Ledger, err error) {
	ledger, err := s.relayService.CreateLedger(ctx, &model.CreateLedgerRequest{Ledger: req.Ledger})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = ledger.ToProto()
	return resp, nil
}

func (s *RelayService) DeleteLedgers(ctx context.Context, req *v1pb.DeleteLedgersRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteLedgers(ctx, &model.DeleteLedgersRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetLedgerList(ctx context.Context, req *v1pb.GetLedgerListRequest) (resp *v1pb.GetLedgerListResponse, err error) {
	total, list, err := s.relayService.GetLedgerList(ctx, &model.GetLedgerListRequest{
		PageParam: types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		AccountId: req.AccountId,
		Type:      model.LedgerType(strings.TrimSpace(req.Type)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	accountIds := lo.UniqMap(list, func(item *model.Ledger, _ int) int64 { return item.AccountId })
	_, accountList, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{Ids: accountIds})
	if err != nil {
		return
	}
	accountMap := lo.Associate(accountList, func(item *model.Account) (int64, string) {
		return item.ID, item.DisplayName()
	})
	resp = &v1pb.GetLedgerListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.Ledger, _ int) *relaypb.Ledger {
			info := item.ToProto()
			info.AccountName = lo.Ternary(accountMap[item.AccountId] != "", accountMap[item.AccountId], "-")
			return info
		}),
	}
	return resp, nil
}

func (s *RelayService) CreateAccountApiKey(ctx context.Context, req *v1pb.CreateAccountApiKeyRequest) (resp *relaypb.AccountApiKey, err error) {
	accountApiKey, err := s.relayService.CreateAccountApiKey(ctx, &model.CreateAccountApiKeyRequest{AccountApiKey: req.AccountApiKey})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	info := accountApiKey.ToProto()
	info.Key = accountApiKey.Key // 添加时，返回真实 key
	resp = info
	return resp, nil
}

func (s *RelayService) UpdateAccountApiKey(ctx context.Context, req *v1pb.UpdateAccountApiKeyRequest) (resp *relaypb.AccountApiKey, err error) {
	accountApiKey, err := s.relayService.UpdateAccountApiKey(ctx, &model.UpdateAccountApiKeyRequest{
		AccountApiKey: req.AccountApiKey,
		UpdateMask:    req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = accountApiKey.ToProto()
	return resp, nil
}

func (s *RelayService) DeleteAccountApiKeys(ctx context.Context, req *v1pb.DeleteAccountApiKeysRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteAccountApiKeys(ctx, &model.DeleteAccountApiKeysRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetAccountApiKeyList(ctx context.Context, req *v1pb.GetAccountApiKeyListRequest) (resp *v1pb.GetAccountApiKeyListResponse, err error) {
	total, list, err := s.relayService.GetAccountApiKeyList(ctx, &model.GetAccountApiKeyListRequest{
		PageParam: types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		AccountId: req.AccountId,
		Keyword:   strings.TrimSpace(req.Keyword),
		Status:    model.ApiKeyStatus(strings.TrimSpace(req.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	accountIds := lo.UniqMap(list, func(item *model.AccountApiKey, _ int) int64 { return item.AccountId })
	_, accountList, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{Ids: accountIds})
	if err != nil {
		return
	}
	accountMap := lo.Associate(accountList, func(item *model.Account) (int64, string) {
		return item.ID, item.DisplayName()
	})
	resp = &v1pb.GetAccountApiKeyListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.AccountApiKey, _ int) *relaypb.AccountApiKey {
			info := item.ToProto()
			info.AccountName = lo.Ternary(accountMap[item.AccountId] != "", accountMap[item.AccountId], "-")
			return info
		}),
	}
	return resp, nil
}

func (s *RelayService) CreateAccount(ctx context.Context, req *v1pb.CreateAccountRequest) (resp *relaypb.Account, err error) {
	account, err := s.relayService.CreateAccount(ctx, &model.CreateAccountRequest{Account: req.Account})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = account.ToProto()
	return resp, nil
}

func (s *RelayService) UpdateAccount(ctx context.Context, req *v1pb.UpdateAccountRequest) (resp *relaypb.Account, err error) {
	account, err := s.relayService.UpdateAccount(ctx, &model.UpdateAccountRequest{
		Account:    req.Account,
		UpdateMask: req.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = account.ToProto()
	return resp, nil
}

func (s *RelayService) DeleteAccounts(ctx context.Context, req *v1pb.DeleteAccountsRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteAccounts(ctx, &model.DeleteAccountsRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetAccountList(ctx context.Context, req *v1pb.GetAccountListRequest) (resp *v1pb.GetAccountListResponse, err error) {
	total, list, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{
		PageParam: types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		Status:    model.EnableStatus(strings.TrimSpace(req.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetAccountListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.Account, _ int) *relaypb.Account {
			return item.ToProto()
		}),
	}
	return resp, nil
}

func (s *RelayService) DeleteRequests(ctx context.Context, req *v1pb.DeleteRequestsRequest) (resp *emptypb.Empty, err error) {
	if err = s.relayService.DeleteRequests(ctx, &model.DeleteRequestsRequest{Ids: req.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &emptypb.Empty{}
	return resp, nil
}

func (s *RelayService) GetRelayUsage(ctx context.Context, req *v1pb.GetRelayUsageRequest) (resp *v1pb.GetRelayUsageResponse, err error) {
	startTime := req.StartTime.AsTime().UTC().Local()
	endTime := req.EndTime.AsTime().UTC().Local()

	if endTime.Sub(startTime) > time.Hour*24*31 {
		err = connect.NewError(connect.CodeInvalidArgument, "time range is too long, max 31 days")
		return
	}

	_, list, err := s.relayService.GetRelayHourlyUsageList(ctx, &model.GetRelayHourlyUsageListRequest{
		PageParam: types.NewPageParam(1, 10000, ""),
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	groupMap := lo.GroupBy(list, func(item *model.RelayHourlyUsage) string { return item.ProviderCode })
	providerCodes := lo.Keys(groupMap)
	sort.Strings(providerCodes)

	// 返回默认值，防止前端没有数据报错
	if len(providerCodes) == 0 {
		providerCodes = append(providerCodes, core.AllProviderCodeList...)
	}

	// 计算总量
	var total int64
	switch req.ChartType {
	case "point":
		total = lo.SumBy(list, func(item *model.RelayHourlyUsage) int64 { return item.TotalPoint })
	default:
		total = lo.SumBy(list, func(item *model.RelayHourlyUsage) int64 { return item.TotalRequest })
	}

	series := make([]*relaypb.UsageSerie, 0, len(providerCodes))
	for _, providerCode := range providerCodes {
		valueList := groupMap[providerCode]
		valueMap := lo.Associate(valueList, func(item *model.RelayHourlyUsage) (string, int64) {
			switch req.ChartType {
			case "point":
				return utils.FormatTime(item.Time, "YmdH"), item.TotalPoint
			default:
				return utils.FormatTime(item.Time, "YmdH"), item.TotalRequest
			}
		})

		var dataList []*relaypb.UsageItem
		for t := startTime; t.Before(endTime); t = t.Add(time.Hour) {
			dataList = append(dataList, &relaypb.UsageItem{
				Label: utils.FormatTime(t, "Y-m-d H"),
				Value: valueMap[utils.FormatTime(t, "YmdH")],
			})
		}
		series = append(series, &relaypb.UsageSerie{
			Name: providerCode,
			Data: dataList,
		})
	}
	resp = &v1pb.GetRelayUsageResponse{
		Series: series,
		Total:  total,
	}
	return resp, nil
}

func (s *RelayService) GetTotalRelayUsage(ctx context.Context, req *v1pb.GetTotalRelayUsageRequest) (resp *v1pb.GetTotalRelayUsageResponse, err error) {
	usage, err := s.relayService.GetRelayUsage(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetTotalRelayUsageResponse{
		TotalRequest: usage.TotalRequest,
		TotalPoint:   usage.TotalPoint,
	}
	return resp, nil
}

func (s *RelayService) GetRelayInfo(ctx context.Context, req *v1pb.GetRelayInfoRequest) (resp *v1pb.GetRelayInfoResponse, err error) {
	relayInfo, err := s.relayService.GetRelayInfo(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	resp = &v1pb.GetRelayInfoResponse{
		ProviderCount: relayInfo.ProviderCount,
		ModelCount:    relayInfo.ModelCount,
		ApiKeyCount:   relayInfo.ApiKeyCount,
	}
	return resp, nil
}

func (s *RelayService) GetRequestList(ctx context.Context, req *v1pb.GetRequestListRequest) (resp *v1pb.GetRequestListResponse, err error) {
	total, list, err := s.relayService.GetRequestList(ctx, &model.GetRequestListRequest{
		PageParam:        types.NewPageParam(int64(req.Current), int64(req.Size), req.OrderBy),
		AccountId:        req.AccountId,
		ProviderCode:     strings.TrimSpace(req.ProviderCode),
		ModelCode:        strings.TrimSpace(req.ModelCode),
		Status:           model.RequestStatus(strings.TrimSpace(req.Status)),
		CompletedAtStart: req.CompletedAtStart.AsTime(),
		CompletedAtEnd:   req.CompletedAtEnd.AsTime(),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	var providerIds []int64
	var modelIds []int64
	var accountIds []int64
	for _, item := range list {
		if !lo.Contains(providerIds, item.ProviderId) {
			providerIds = append(providerIds, item.ProviderId)
		}
		if !lo.Contains(modelIds, item.ModelId) {
			modelIds = append(modelIds, item.ModelId)
		}
		if !lo.Contains(accountIds, item.AccountId) {
			accountIds = append(accountIds, item.AccountId)
		}
	}

	_, providerList, err := s.relayService.GetProviderList(ctx, &model.GetProviderListRequest{Ids: providerIds})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	providerMap := lo.Associate(providerList, func(item *model.Provider) (int64, string) {
		return item.ID, item.Code
	})
	_, modelList, err := s.relayService.GetModelList(ctx, &model.GetModelListRequest{Ids: modelIds})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	modelMap := lo.Associate(modelList, func(item *model.Model) (int64, string) {
		return item.ID, item.Code
	})
	_, users, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{Ids: accountIds})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err.Error()).WithCause(err)
		return
	}
	userMap := lo.Associate(users, func(item *model.Account) (int64, string) {
		return item.ID, item.DisplayName()
	})

	resp = &v1pb.GetRequestListResponse{
		Current: req.Current,
		Size:    req.Size,
		Total:   uint32(total),
		Records: lo.Map(list, func(item *model.Request, _ int) *relaypb.Request {
			info := item.ToProto()
			info.ProviderCode = lo.Ternary(providerMap[item.ProviderId] != "", providerMap[item.ProviderId], "-")
			info.ModelCode = lo.Ternary(modelMap[item.ModelId] != "", modelMap[item.ModelId], "-")
			info.AccountName = lo.Ternary(userMap[item.AccountId] != "", userMap[item.AccountId], "-")
			return info
		}),
	}
	return resp, nil
}
