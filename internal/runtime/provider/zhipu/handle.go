package zhipu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/provider/openai"
	"github.com/modelgate/modelgate/pkg/utils"
)

type Handler struct {
	base     *openai.Handler
	provider string
}

func NewHandler(provider string) *Handler {
	return &Handler{
		base:     openai.NewHandler(provider),
		provider: provider,
	}
}

func (h *Handler) Provider() string {
	return h.provider
}

// BeforeRequest 构建请求参数
func (h *Handler) BeforeRequest(ctx context.Context, c *core.Context) (err error) {
	var endpoint string
	baseUrl := strings.TrimRight(c.CurrentModel.BaseUrl, "/")
	if h.provider == core.ProviderCodeZhipuClaude {
		endpoint = baseUrl + "/api/anthropic/" + strings.TrimPrefix(c.UrlPath, "/")
	} else {
		endpoint = baseUrl + "/api/paas/v4/" + strings.TrimPrefix(c.UrlPath, "/v1/")
	}

	log.Infof("%s handler before request, provider: %s, model: %s, endpoint: %s", h.provider, c.CurrentModel.ProviderCode, c.CurrentModel.ModelCode, endpoint)

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
		log.Errorf("%s provider do request error: %v", h.provider, string(c.RawResponse))
		if h.provider == core.ProviderCodeZhipuClaude {
			var respData struct {
				Error struct {
					Message string
				}
			}
			if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
				err = fmt.Errorf("unmarshal %s response error: %v", h.provider, err)
				return
			}
			err = fmt.Errorf("%s response error: %s", h.provider, respData.Error.Message)
		} else {
			var respData struct {
				Detail string
			}
			if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
				err = fmt.Errorf("unmarshal %s response error: %v", h.provider, err)
				return
			}
			err = fmt.Errorf("%s response error: %s", h.provider, respData.Detail)
		}
		return
	}
	return
}

// AfterResponse 处理响应结果
func (h *Handler) AfterResponse(ctx context.Context, c *core.Context) error {
	return h.base.AfterResponse(ctx, c)
}

// DoStream 发送流式请求
func (h *Handler) DoStream(ctx context.Context, c *core.Context) (stream core.Stream, err error) {
	return h.base.DoStream(ctx, c)
}
