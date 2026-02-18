package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openai/openai-go"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/pkg/utils"
)

type Handler struct {
	provider string
}

func NewHandler(provider string) *Handler {
	return &Handler{
		provider: provider,
	}
}

func (h *Handler) Provider() string {
	return h.provider
}

// BeforeRequest 构建请求参数
func (h *Handler) BeforeRequest(ctx context.Context, c *core.Context) (err error) {
	endpoint := c.CurrentModel.BaseUrl + c.UrlPath
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewReader(c.InputBody),
	)
	if err != nil {
		return
	}

	apiKey, err := utils.DecryptAESGCM(c.CurrentModel.ApiKeyEncrypted, []byte(config.GetConfig().Secret.Key))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	c.HTTPRequest = req
	return
}

// DoRequest 发送请求，并处理结果
func (h *Handler) DoRequest(ctx context.Context, c *core.Context) (err error) {
	resp, err := http.DefaultClient.Do(c.HTTPRequest)
	if err != nil {
		return
	}
	c.HTTPResponse = resp

	defer resp.Body.Close()
	c.RawResponse, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error(h.provider, string(c.RawResponse))
		err = h.parseResponseError(ctx, c)
		return
	}
	return
}

func (h *Handler) parseResponseError(_ context.Context, c *core.Context) (err error) {
	switch h.provider {
	case core.ProviderCodeOpenAI:
		var respData struct {
			Error struct {
				Message string
				Type    string
			}
		}
		if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
			err = fmt.Errorf("unmarshal %s response error: %v", h.provider, err)
			return
		}
		err = fmt.Errorf("%s response error: %s", h.provider, respData.Error.Message)
		return
	case core.ProviderCodeDeepSeek:
		err = fmt.Errorf("%s response error: %s", h.provider, string(c.RawResponse))
		return
	default:
		err = fmt.Errorf("%s response error: %s", h.provider, string(c.RawResponse))
		return
	}
}

// AfterResponse 处理响应结果
func (h *Handler) AfterResponse(ctx context.Context, c *core.Context) (err error) {
	var respData struct {
		Model string                 `json:"model"`
		Usage openai.CompletionUsage `json:"usage"`
	}
	if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
		return
	}
	c.ActualModel = respData.Model
	if respData.Usage.TotalTokens > 0 {
		c.Usage = &core.Usage{
			PromptTokens:       respData.Usage.PromptTokens,
			PromptCachedTokens: respData.Usage.PromptTokensDetails.CachedTokens,
			CompletionTokens:   respData.Usage.CompletionTokens,
			TotalTokens:        respData.Usage.TotalTokens,
		}
	}
	return
}

// DoStream 发送流式请求
func (h *Handler) DoStream(ctx context.Context, c *core.Context) (stream core.Stream, err error) {
	resp, err := http.DefaultClient.Do(c.HTTPRequest)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai stream error: %s", b)
	}

	c.HTTPResponse = resp
	return NewStreamReceiver(resp.Body), nil
}
