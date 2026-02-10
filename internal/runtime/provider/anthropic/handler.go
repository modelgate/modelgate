package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	endpoint := c.CurrentModel.BaseUrl + "/v1/messages"
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
	anthropicVersion := "2023-06-01"
	if c.Header.Get("anthropic-version") != "" {
		anthropicVersion = c.Header.Get("anthropic-version")
	}
	anthropicBeta := "messages-2023-12-15"
	if strings.HasPrefix(c.CurrentModel.ModelCode, "claude-3-5-sonnet") {
		anthropicBeta = "max-tokens-3-5-sonnet-2024-07-15"
	}
	req.Header.Set("x-api-key", string(apiKey))
	req.Header.Set("anthropic-version", anthropicVersion)
	req.Header.Set("anthropic-beta", anthropicBeta)
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
		var respData struct {
			RequestId string
			Type      string
			Error     struct {
				Type    string
				Message string
			}
		}
		if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
			err = fmt.Errorf("unmarshal %s response error: %v", h.provider, err)
			return
		}
		err = fmt.Errorf("%s response error: %s", h.provider, respData.Error.Message)
		return
	}
	return
}

// AfterResponse 处理响应结果
func (h *Handler) AfterResponse(ctx context.Context, c *core.Context) (err error) {
	var respData struct {
		Model string `json:"model"`
		Usage Usage  `json:"usage"`
	}
	if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
		err = fmt.Errorf("unmarshal response error: %v", err)
		return
	}
	c.ActualModel = respData.Model
	if usage := respData.Usage; usage.InputTokens > 0 || usage.OutputTokens > 0 {
		c.Usage = &core.Usage{
			PromptTokens:     usage.InputTokens,
			CompletionTokens: usage.OutputTokens,
			TotalTokens:      usage.InputTokens + usage.OutputTokens,
		}
	}
	return
}
