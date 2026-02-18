package minimax

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/provider/anthropic"
	"github.com/modelgate/modelgate/internal/runtime/provider/openai"
	"github.com/modelgate/modelgate/pkg/utils"
)

// OpenAIHandler MiniMax OpenAI 协议处理器，继承 openai.Handler
// 仅覆写 BeforeRequest 以实现 MiniMax 特有的端点拼接：baseUrl/<path>
type OpenAIHandler struct {
	*openai.Handler
}

func NewOpenAIHandler() *OpenAIHandler {
	return &OpenAIHandler{
		Handler: openai.NewHandler(core.ProviderCodeMinimax),
	}
}

func (h *OpenAIHandler) BeforeRequest(ctx context.Context, c *core.Context) (err error) {
	baseUrl := strings.TrimRight(c.CurrentModel.BaseUrl, "/")
	endpoint := baseUrl + c.UrlPath

	log.Infof("minimax openai handler, model: %s, endpoint: %s", c.CurrentModel.ModelCode, endpoint)

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

// AnthropicHandler MiniMax Anthropic 协议处理器，继承 anthropic.Handler
// 仅覆写 BeforeRequest 以实现 MiniMax 特有的端点拼接：baseUrl/anthropic/<path>
type AnthropicHandler struct {
	*anthropic.Handler
}

func NewAnthropicHandler() *AnthropicHandler {
	return &AnthropicHandler{
		Handler: anthropic.NewHandler(core.ProviderCodeMinimax),
	}
}

func (h *AnthropicHandler) BeforeRequest(ctx context.Context, c *core.Context) (err error) {
	baseUrl := strings.TrimRight(c.CurrentModel.BaseUrl, "/")
	endpoint := baseUrl + "/anthropic/" + strings.TrimPrefix(c.UrlPath, "/")

	log.Infof("minimax anthropic handler, model: %s, endpoint: %s", c.CurrentModel.ModelCode, endpoint)

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
