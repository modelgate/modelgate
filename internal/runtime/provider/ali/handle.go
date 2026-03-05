package ali

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
	"github.com/modelgate/modelgate/internal/runtime/provider/anthropic"
	"github.com/modelgate/modelgate/internal/runtime/provider/openai"
	"github.com/modelgate/modelgate/pkg/utils"
)

// OpenAIHandler 阿里云百炼 OpenAI 协议处理器，继承 openai.Handler
// 仅覆写 BeforeRequest 以实现阿里云百炼特有的端点拼接：baseUrl/<path>
type OpenAIHandler struct {
	*openai.Handler
}

func NewOpenAIHandler() *OpenAIHandler {
	return &OpenAIHandler{
		Handler: openai.NewHandler(core.ProviderCodeAli),
	}
}

func (h *OpenAIHandler) BeforeRequest(ctx context.Context, c *core.Context) (err error) {
	baseUrl := strings.TrimRight(c.CurrentModel.BaseUrl, "/")
	endpoint := baseUrl + c.UrlPath

	log.Infof("ali openai handler, model: %s, endpoint: %s", c.CurrentModel.ModelCode, endpoint)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(c.InputBody))
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

// AnthropicHandler 阿里云百炼 Anthropic 协议处理器，继承 anthropic.Handler
// 仅覆写 BeforeRequest 以实现阿里云百炼特有的端点拼接：baseUrl/apps/anthropic/<path>
// DoRequest 也需覆写，因为阿里云百炼的错误格式与标准 Anthropic 不同
type AnthropicHandler struct {
	*anthropic.Handler
}

func NewAnthropicHandler() *AnthropicHandler {
	return &AnthropicHandler{
		Handler: anthropic.NewHandler(core.ProviderCodeAli),
	}
}

func (h *AnthropicHandler) BeforeRequest(ctx context.Context, c *core.Context) (err error) {
	baseUrl := strings.TrimRight(c.CurrentModel.BaseUrl, "/")
	endpoint := baseUrl + "/apps/anthropic/" + strings.TrimPrefix(c.UrlPath, "/")

	log.Infof("ali anthropic handler, model: %s, endpoint: %s", c.CurrentModel.ModelCode, endpoint)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(c.InputBody))
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

func (h *AnthropicHandler) DoRequest(ctx context.Context, c *core.Context) (err error) {
	resp, err := core.HttpClient.Do(c.HTTPRequest)
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
		log.Errorf("ali anthropic do request error, status code: %d, response: %s", resp.StatusCode, string(c.RawResponse))
		var respData struct {
			Error struct {
				Message string
			}
		}
		if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
			err = fmt.Errorf("unmarshal %s response error: %v", h.Provider(), err)
			return
		}
		err = fmt.Errorf("%s response error: %s", h.Provider(), respData.Error.Message)
	}
	return
}
