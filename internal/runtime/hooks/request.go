package hooks

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/relay/model"
	"github.com/modelgate/modelgate/internal/runtime/core"
)

// RequestHook 请求Hook
type RequestHook struct {
	service relay.Service
}

var _ core.Hook = (*RequestHook)(nil)

// NewRequestHook 创建请求Hook
func NewRequestHook(i do.Injector) (*RequestHook, error) {
	return &RequestHook{
		service: do.MustInvoke[relay.Service](i),
	}, nil
}

func (h *RequestHook) Name() string {
	return "request"
}

// Before 执行前
func (h *RequestHook) Before(ctx context.Context, c *core.Context) (err error) {
	_, err = h.service.CreateRequest(ctx, &model.CreateRequestRequest{
		RequestUUID:      c.RequestUUID,
		AttemptNo:        c.AttemptNo,
		AccountId:        c.AccountId,
		ProviderCode:     c.CurrentModel.ProviderCode,
		AccountApiKeyId:  c.AccountApiKeyId,
		ProviderId:       c.CurrentModel.ProviderId,
		ProviderApiKeyId: c.CurrentModel.ApiKeyId,
		ModelId:          c.CurrentModel.ModelId,
		PromptTokens:     0,
		CompletionTokens: 0,
		TotalTokens:      0,
		Status:           model.RequestStatusPending,
	})
	return
}

// After 执行后
func (h *RequestHook) After(ctx context.Context, c *core.Context) (err error) {
	var req = model.UpdateRequestCompletedRequest{
		RequestUUID:  c.RequestUUID,
		AttemptNo:    c.AttemptNo,
		ProviderCode: c.CurrentModel.ProviderCode,
		ActualModel:  lo.Ternary(c.ActualModel != "", c.ActualModel, c.CurrentModel.ModelCode),
	}
	if c.Usage != nil {
		req.PromptTokens = c.Usage.PromptTokens
		req.CompletionTokens = c.Usage.CompletionTokens
		req.TotalTokens = c.Usage.TotalTokens
	} else {
		req.PromptTokens = int64(c.PromptTokens)
		req.CompletionTokens = int64(c.CompletionTokens)
		req.TotalTokens = int64(c.PromptTokens + c.CompletionTokens)
	}
	if c.LastErr != nil {
		req.Status = model.RequestStatusFailed
		req.ErrorMessage = c.LastErr.Error()
		if c.HTTPResponse != nil {
			req.ErrorCode = c.HTTPResponse.StatusCode
		}
	} else {
		req.Status = model.RequestStatusSuccess
	}
	err = h.service.UpdateRequestCompleted(ctx, &req)
	return
}

// OnChunk 流chunk
func (h *RequestHook) OnChunk(ctx context.Context, c *core.Context, chunk *core.StreamChunk) (err error) {
	return
}

// OnError 流错误
func (h *RequestHook) OnError(ctx context.Context, c *core.Context, err error) {
}
