package gemini

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/config"
	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/internal/runtime/provider/openai"
	"github.com/modelgate/modelgate/pkg/utils"
)

// OpenAIHandler Gemini OpenAI 协议处理器，继承 openai.Handler
// 仅覆写 BeforeRequest 以实现 Gemini 特有的端点拼接和认证方式
type OpenAIHandler struct {
	*openai.Handler
}

func NewOpenAIHandler() *OpenAIHandler {
	return &OpenAIHandler{
		Handler: openai.NewHandler(core.ProviderCodeGemini),
	}
}

func (h *OpenAIHandler) BeforeRequest(ctx context.Context, c *core.Context) (err error) {
	baseUrl := strings.TrimRight(c.CurrentModel.BaseUrl, "/")
	endpoint := baseUrl + c.UrlPath

	log.Infof("gemini openai handler, model: %s, endpoint: %s", c.CurrentModel.ModelCode, endpoint)

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