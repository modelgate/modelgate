package hooks

import (
	"context"
	"math"

	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/internal/runtime/core"
)

type BillingHook struct {
	service relay.Service
}

var _ core.Hook = (*BillingHook)(nil)

func NewBillingHook(i do.Injector) (*BillingHook, error) {
	return &BillingHook{
		service: do.MustInvoke[relay.Service](i),
	}, nil
}

func (h *BillingHook) Name() string {
	return "billing"
}

// Before 预扣款
func (h *BillingHook) Before(ctx context.Context, c *core.Context) (err error) {
	modelInfo := c.CurrentModel
	cost := int64(math.Ceil(modelInfo.InputPrice * float64(c.PromptTokens) * float64(modelInfo.PointsPerCurrency) / float64(modelInfo.TokenNum)))
	log.Infof("prepay cost: %d", cost)

	_, err = h.service.DeductBalance(ctx, c.AccountId, cost, c.RequestId, model.LedgerTypeConsume, "reserve")
	if err != nil {
		return
	}
	c.PreCost = cost
	return
}

// After 扣款
func (h *BillingHook) After(ctx context.Context, c *core.Context) (err error) {
	modelInfo := c.CurrentModel

	var promptTokens int64
	var promptCacheTokens int64
	var completionTokens int64
	if c.Usage != nil {
		promptTokens = int64(c.Usage.PromptTokens)
		promptCacheTokens = int64(c.Usage.PromptCachedTokens)
		completionTokens = int64(c.Usage.CompletionTokens)
	} else {
		promptTokens = int64(c.PromptTokens)
		completionTokens = int64(c.CompletionTokens)
	}

	totalCost := int64(math.Ceil(modelInfo.InputPrice * float64(promptTokens) * float64(modelInfo.PointsPerCurrency) / float64(modelInfo.TokenNum)))
	if promptCacheTokens > 0 {
		totalCost += int64(math.Ceil(modelInfo.InputCachePrice * float64(promptCacheTokens) * float64(modelInfo.PointsPerCurrency) / float64(modelInfo.TokenNum)))
	}
	totalCost += int64(math.Ceil(modelInfo.OutputPrice * float64(completionTokens) * float64(modelInfo.PointsPerCurrency) / float64(modelInfo.TokenNum)))
	log.Infof("total cost: %d", totalCost)

	if v := totalCost - c.PreCost; v > 0 {
		_, err = h.service.DeductBalance(ctx, c.AccountId, v, c.RequestId, model.LedgerTypeConsume, "settle")
		if err != nil {
			return
		}
	} else if v < 0 {
		_, err = h.service.AddBalance(ctx, c.AccountId, -v, c.RequestId, model.LedgerTypeRefund, "settle")
		if err != nil {
			return
		}
	}
	c.TotalCost = totalCost
	// 记录点数使用情况
	if eErr := h.service.AddPointUsage(ctx, modelInfo.ProviderCode, modelInfo.ApiKeyId, c.AccountApiKeyId, totalCost); eErr != nil {
		log.Errorf("AddPointUsage provider_code: %s, provider_api_key: %d, account_api_key: %d, total_cost: %d, error: %v", modelInfo.ProviderCode, modelInfo.ApiKeyId, c.AccountApiKeyId, totalCost, eErr)
	}
	return
}

func (h *BillingHook) OnChunk(ctx context.Context, c *core.Context, chunk *core.StreamChunk) (err error) {
	return
}

func (h *BillingHook) OnError(ctx context.Context, c *core.Context, err error) {
}
