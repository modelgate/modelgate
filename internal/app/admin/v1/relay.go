package v1

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"connectrpc.com/connect"
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

func (s *RelayService) CreateProvider(ctx context.Context, req *connect.Request[v1pb.CreateProviderRequest]) (resp *connect.Response[relaypb.Provider], err error) {
	provider, err := s.relayService.CreateProvider(ctx, &model.CreateProviderRequest{Provider: req.Msg.Provider})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(provider.ToProto())
	return resp, nil
}

func (s *RelayService) UpdateProvider(ctx context.Context, req *connect.Request[v1pb.UpdateProviderRequest]) (resp *connect.Response[relaypb.Provider], err error) {
	provider, err := s.relayService.UpdateProvider(ctx, &model.UpdateProviderRequest{
		Provider:   req.Msg.Provider,
		UpdateMask: req.Msg.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(provider.ToProto())
	return resp, nil
}

func (s *RelayService) DeleteProviders(ctx context.Context, req *connect.Request[v1pb.DeleteProvidersRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteProviders(ctx, &model.DeleteProvidersRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetProviderList(ctx context.Context, req *connect.Request[v1pb.GetProviderListRequest]) (resp *connect.Response[v1pb.GetProviderListResponse], err error) {
	total, list, err := s.relayService.GetProviderList(ctx, &model.GetProviderListRequest{
		PageParam: types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		Name:      strings.TrimSpace(req.Msg.Name),
		Code:      strings.TrimSpace(req.Msg.Code),
		Status:    model.EnableStatus(strings.TrimSpace(req.Msg.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(
		&v1pb.GetProviderListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.Provider, _ int) *relaypb.Provider {
				return item.ToProto()
			}),
		})
	return resp, nil
}

func (s *RelayService) GetProviderCodeList(ctx context.Context, req *connect.Request[v1pb.GetProviderCodeListRequest]) (resp *connect.Response[v1pb.GetProviderCodeListResponse], err error) {
	resp = connect.NewResponse(&v1pb.GetProviderCodeListResponse{Records: core.AllProviderCodeList})
	return resp, nil
}

func (s *RelayService) CreateModel(ctx context.Context, req *connect.Request[v1pb.CreateModelRequest]) (resp *connect.Response[relaypb.Model], err error) {
	model, err := s.relayService.CreateModel(ctx, &model.CreateModelRequest{Model: req.Msg.Model})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(model.ToProto())
	return resp, nil
}

func (s *RelayService) UpdateModel(ctx context.Context, req *connect.Request[v1pb.UpdateModelRequest]) (resp *connect.Response[relaypb.Model], err error) {
	model, err := s.relayService.UpdateModel(ctx, &model.UpdateModelRequest{
		Model:      req.Msg.Model,
		UpdateMask: req.Msg.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(model.ToProto())
	return resp, nil
}

func (s *RelayService) DeleteModels(ctx context.Context, req *connect.Request[v1pb.DeleteModelsRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteModels(ctx, &model.DeleteModelsRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetModelList(ctx context.Context, req *connect.Request[v1pb.GetModelListRequest]) (resp *connect.Response[v1pb.GetModelListResponse], err error) {
	total, list, err := s.relayService.GetModelList(ctx, &model.GetModelListRequest{
		PageParam:    types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		Name:         strings.TrimSpace(req.Msg.Name),
		Code:         strings.TrimSpace(req.Msg.Code),
		ProviderCode: strings.TrimSpace(req.Msg.ProviderCode),
		Status:       model.ModelStatus(strings.TrimSpace(req.Msg.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(
		&v1pb.GetModelListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.Model, _ int) *relaypb.Model {
				return item.ToProto()
			}),
		})
	return resp, nil
}

func (s *RelayService) CreateProviderApiKey(ctx context.Context, req *connect.Request[v1pb.CreateProviderApiKeyRequest]) (resp *connect.Response[relaypb.ProviderApiKey], err error) {
	providerApiKey, err := s.relayService.CreateProviderApiKey(ctx, &model.CreateProviderApiKeyRequest{ProviderApiKey: req.Msg.ProviderApiKey})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(providerApiKey.ToProto())
	return resp, nil
}

func (s *RelayService) UpdateProviderApiKey(ctx context.Context, req *connect.Request[v1pb.UpdateProviderApiKeyRequest]) (resp *connect.Response[relaypb.ProviderApiKey], err error) {
	providerApiKey, err := s.relayService.UpdateProviderApiKey(ctx, &model.UpdateProviderApiKeyRequest{
		ProviderApiKey: req.Msg.ProviderApiKey,
		UpdateMask:     req.Msg.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(providerApiKey.ToProto())
	return resp, nil
}

func (s *RelayService) DeleteProviderApiKeys(ctx context.Context, req *connect.Request[v1pb.DeleteProviderApiKeysRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteProviderApiKeys(ctx, &model.DeleteProviderApiKeysRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetProviderApiKeyList(ctx context.Context, req *connect.Request[v1pb.GetProviderApiKeyListRequest]) (resp *connect.Response[v1pb.GetProviderApiKeyListResponse], err error) {
	total, list, err := s.relayService.GetProviderApiKeyList(ctx, &model.GetProviderApiKeyListRequest{
		PageParam:    types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		ProviderId:   req.Msg.ProviderId,
		ProviderCode: strings.TrimSpace(req.Msg.ProviderCode),
		Name:         strings.TrimSpace(req.Msg.Name),
		Status:       model.ApiKeyStatus(strings.TrimSpace(req.Msg.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(
		&v1pb.GetProviderApiKeyListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.ProviderApiKey, _ int) *relaypb.ProviderApiKey {
				return item.ToProto()
			}),
		})
	return resp, nil
}

func (s *RelayService) CreateModelPricing(ctx context.Context, req *connect.Request[v1pb.CreateModelPricingRequest]) (resp *connect.Response[relaypb.ModelPricing], err error) {
	modelPricing, err := s.relayService.CreateModelPricing(ctx, &model.CreateModelPricingRequest{ModelPricing: req.Msg.ModelPricing})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(modelPricing.ToProto())
	return resp, nil
}

func (s *RelayService) UpdateModelPricing(ctx context.Context, req *connect.Request[v1pb.UpdateModelPricingRequest]) (resp *connect.Response[relaypb.ModelPricing], err error) {
	modelPricing, err := s.relayService.UpdateModelPricing(ctx, &model.UpdateModelPricingRequest{
		ModelPricing: req.Msg.ModelPricing,
		UpdateMask:   req.Msg.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(modelPricing.ToProto())
	return resp, nil
}

func (s *RelayService) DeleteModelPricings(ctx context.Context, req *connect.Request[v1pb.DeleteModelPricingsRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteModelPricings(ctx, &model.DeleteModelPricingsRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetModelPricingList(ctx context.Context, req *connect.Request[v1pb.GetModelPricingListRequest]) (resp *connect.Response[v1pb.GetModelPricingListResponse], err error) {
	total, list, err := s.relayService.GetModelPricingList(ctx, &model.GetModelPricingListRequest{
		PageParam:     types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		ProviderCode:  strings.TrimSpace(req.Msg.ProviderCode),
		ModelCode:     strings.TrimSpace(req.Msg.ModelCode),
		Currency:      model.Currency(strings.TrimSpace(req.Msg.Currency)),
		EffectiveFrom: req.Msg.EffectiveFrom.AsTime(),
		EffectiveTo:   req.Msg.EffectiveTo.AsTime(),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(
		&v1pb.GetModelPricingListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.ModelPricing, _ int) *relaypb.ModelPricing {
				return item.ToProto()
			}),
		})
	return resp, nil
}

func (s *RelayService) CreateLedger(ctx context.Context, req *connect.Request[v1pb.CreateLedgerRequest]) (resp *connect.Response[relaypb.Ledger], err error) {
	ledger, err := s.relayService.CreateLedger(ctx, &model.CreateLedgerRequest{Ledger: req.Msg.Ledger})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(ledger.ToProto())
	return resp, nil
}

func (s *RelayService) DeleteLedgers(ctx context.Context, req *connect.Request[v1pb.DeleteLedgersRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteLedgers(ctx, &model.DeleteLedgersRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetLedgerList(ctx context.Context, req *connect.Request[v1pb.GetLedgerListRequest]) (resp *connect.Response[v1pb.GetLedgerListResponse], err error) {
	total, list, err := s.relayService.GetLedgerList(ctx, &model.GetLedgerListRequest{
		PageParam: types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		AccountId: req.Msg.AccountId,
		Type:      model.LedgerType(strings.TrimSpace(req.Msg.Type)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	accountIds := lo.UniqMap(list, func(item *model.Ledger, _ int) int64 { return item.AccountId })
	_, accountList, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{Ids: accountIds})
	if err != nil {
		return
	}
	accountMap := lo.Associate(accountList, func(item *model.Account) (int64, string) {
		return item.ID, item.Name
	})
	resp = connect.NewResponse(
		&v1pb.GetLedgerListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.Ledger, _ int) *relaypb.Ledger {
				info := item.ToProto()
				info.AccountName = lo.Ternary(accountMap[item.AccountId] != "", accountMap[item.AccountId], "-")
				return info
			}),
		})
	return resp, nil
}

func (s *RelayService) CreateAccountApiKey(ctx context.Context, req *connect.Request[v1pb.CreateAccountApiKeyRequest]) (resp *connect.Response[relaypb.AccountApiKey], err error) {
	accountApiKey, err := s.relayService.CreateAccountApiKey(ctx, &model.CreateAccountApiKeyRequest{AccountApiKey: req.Msg.AccountApiKey})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	info := accountApiKey.ToProto()
	info.Key = accountApiKey.Key // 添加时，返回真实 key
	resp = connect.NewResponse(info)
	return resp, nil
}

func (s *RelayService) UpdateAccountApiKey(ctx context.Context, req *connect.Request[v1pb.UpdateAccountApiKeyRequest]) (resp *connect.Response[relaypb.AccountApiKey], err error) {
	accountApiKey, err := s.relayService.UpdateAccountApiKey(ctx, &model.UpdateAccountApiKeyRequest{
		AccountApiKey: req.Msg.AccountApiKey,
		UpdateMask:    req.Msg.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(accountApiKey.ToProto())
	return resp, nil
}

func (s *RelayService) DeleteAccountApiKeys(ctx context.Context, req *connect.Request[v1pb.DeleteAccountApiKeysRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteAccountApiKeys(ctx, &model.DeleteAccountApiKeysRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetAccountApiKeyList(ctx context.Context, req *connect.Request[v1pb.GetAccountApiKeyListRequest]) (resp *connect.Response[v1pb.GetAccountApiKeyListResponse], err error) {
	total, list, err := s.relayService.GetAccountApiKeyList(ctx, &model.GetAccountApiKeyListRequest{
		PageParam: types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		AccountId: req.Msg.AccountId,
		Keyword:   strings.TrimSpace(req.Msg.Keyword),
		Status:    model.ApiKeyStatus(strings.TrimSpace(req.Msg.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	accountIds := lo.UniqMap(list, func(item *model.AccountApiKey, _ int) int64 { return item.AccountId })
	_, accountList, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{Ids: accountIds})
	if err != nil {
		return
	}
	accountMap := lo.Associate(accountList, func(item *model.Account) (int64, string) {
		return item.ID, item.Name
	})
	resp = connect.NewResponse(
		&v1pb.GetAccountApiKeyListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.AccountApiKey, _ int) *relaypb.AccountApiKey {
				info := item.ToProto()
				info.AccountName = lo.Ternary(accountMap[item.AccountId] != "", accountMap[item.AccountId], "-")
				return info
			}),
		})
	return resp, nil
}

func (s *RelayService) CreateAccount(ctx context.Context, req *connect.Request[v1pb.CreateAccountRequest]) (resp *connect.Response[relaypb.Account], err error) {
	account, err := s.relayService.CreateAccount(ctx, &model.CreateAccountRequest{Account: req.Msg.Account})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(account.ToProto())
	return resp, nil
}

func (s *RelayService) UpdateAccount(ctx context.Context, req *connect.Request[v1pb.UpdateAccountRequest]) (resp *connect.Response[relaypb.Account], err error) {
	account, err := s.relayService.UpdateAccount(ctx, &model.UpdateAccountRequest{
		Account:    req.Msg.Account,
		UpdateMask: req.Msg.UpdateMask.Paths,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(account.ToProto())
	return resp, nil
}

func (s *RelayService) DeleteAccounts(ctx context.Context, req *connect.Request[v1pb.DeleteAccountsRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteAccounts(ctx, &model.DeleteAccountsRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetAccountList(ctx context.Context, req *connect.Request[v1pb.GetAccountListRequest]) (resp *connect.Response[v1pb.GetAccountListResponse], err error) {
	total, list, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{
		PageParam: types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		Name:      strings.TrimSpace(req.Msg.Name),
		Nickname:  strings.TrimSpace(req.Msg.Nickname),
		Status:    model.EnableStatus(strings.TrimSpace(req.Msg.Status)),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(
		&v1pb.GetAccountListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.Account, _ int) *relaypb.Account {
				return item.ToProto()
			}),
		})
	return resp, nil
}

func (s *RelayService) DeleteRequests(ctx context.Context, req *connect.Request[v1pb.DeleteRequestsRequest]) (resp *connect.Response[emptypb.Empty], err error) {
	if err = s.relayService.DeleteRequests(ctx, &model.DeleteRequestsRequest{Ids: req.Msg.Ids}); err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&emptypb.Empty{})
	return resp, nil
}

func (s *RelayService) GetRelayUsage(ctx context.Context, req *connect.Request[v1pb.GetRelayUsageRequest]) (resp *connect.Response[v1pb.GetRelayUsageResponse], err error) {
	startTime := req.Msg.StartTime.AsTime().UTC().Local()
	endTime := req.Msg.EndTime.AsTime().UTC().Local()

	if endTime.Sub(startTime) > time.Hour*24*31 {
		err = connect.NewError(connect.CodeInvalidArgument, errors.New("time range is too long, max 31 days"))
		return
	}

	_, list, err := s.relayService.GetRelayHourlyUsageList(ctx, &model.GetRelayHourlyUsageListRequest{
		PageParam: types.NewPageParam(1, 10000, ""),
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	groupMap := lo.GroupBy(list, func(item *model.RelayHourlyUsage) string { return item.ProviderCode })
	providerCodes := lo.Keys(groupMap)
	sort.Strings(providerCodes)

	// 返回默认值，防止前端没有数据报错
	if len(providerCodes) == 0 {
		providerCodes = append(providerCodes, core.AllProviderCodeList...)
	}

	series := make([]*relaypb.UsageSerie, 0, len(providerCodes))
	for _, providerCode := range providerCodes {
		valueList := groupMap[providerCode]
		valueMap := lo.Associate(valueList, func(item *model.RelayHourlyUsage) (string, int64) {
			switch req.Msg.ChartType {
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
	resp = connect.NewResponse(&v1pb.GetRelayUsageResponse{Series: series})
	return resp, nil
}

func (s *RelayService) GetTotalRelayUsage(ctx context.Context, req *connect.Request[v1pb.GetTotalRelayUsageRequest]) (resp *connect.Response[v1pb.GetTotalRelayUsageResponse], err error) {
	usage, err := s.relayService.GetRelayUsage(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&v1pb.GetTotalRelayUsageResponse{
		TotalRequest: usage.TotalRequest,
		TotalPoint:   usage.TotalPoint,
	})
	return resp, nil
}

func (s *RelayService) GetRelayInfo(ctx context.Context, req *connect.Request[v1pb.GetRelayInfoRequest]) (resp *connect.Response[v1pb.GetRelayInfoResponse], err error) {
	relayInfo, err := s.relayService.GetRelayInfo(ctx)
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	resp = connect.NewResponse(&v1pb.GetRelayInfoResponse{
		ProviderCount: relayInfo.ProviderCount,
		ModelCount:    relayInfo.ModelCount,
		ApiKeyCount:   relayInfo.ApiKeyCount,
	})
	return resp, nil
}

func (s *RelayService) GetRequestList(ctx context.Context, req *connect.Request[v1pb.GetRequestListRequest]) (resp *connect.Response[v1pb.GetRequestListResponse], err error) {
	total, list, err := s.relayService.GetRequestList(ctx, &model.GetRequestListRequest{
		PageParam:        types.NewPageParam(int64(req.Msg.Current), int64(req.Msg.Size), req.Msg.OrderBy),
		AccountId:        req.Msg.AccountId,
		ProviderCode:     strings.TrimSpace(req.Msg.ProviderCode),
		ModelCode:        strings.TrimSpace(req.Msg.ModelCode),
		Status:           model.RequestStatus(strings.TrimSpace(req.Msg.Status)),
		CompletedAtStart: req.Msg.CompletedAtStart.AsTime(),
		CompletedAtEnd:   req.Msg.CompletedAtEnd.AsTime(),
	})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
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
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	providerMap := lo.Associate(providerList, func(item *model.Provider) (int64, string) {
		return item.ID, item.Code
	})
	_, modelList, err := s.relayService.GetModelList(ctx, &model.GetModelListRequest{Ids: modelIds})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	modelMap := lo.Associate(modelList, func(item *model.Model) (int64, string) {
		return item.ID, item.Code
	})
	_, users, err := s.relayService.GetAccountList(ctx, &model.GetAccountListRequest{Ids: accountIds})
	if err != nil {
		err = connect.NewError(connect.CodeInternal, err)
		return
	}
	userMap := lo.Associate(users, func(item *model.Account) (int64, string) {
		return item.ID, item.Name
	})

	resp = connect.NewResponse(
		&v1pb.GetRequestListResponse{
			Current: req.Msg.Current,
			Size:    req.Msg.Size,
			Total:   uint32(total),
			Records: lo.Map(list, func(item *model.Request, _ int) *relaypb.Request {
				info := item.ToProto()
				info.ProviderCode = lo.Ternary(providerMap[item.ProviderId] != "", providerMap[item.ProviderId], "-")
				info.ModelCode = lo.Ternary(modelMap[item.ModelId] != "", modelMap[item.ModelId], "-")
				info.AccountName = lo.Ternary(userMap[item.AccountId] != "", userMap[item.AccountId], "-")
				return info
			}),
		})
	return resp, nil
}
