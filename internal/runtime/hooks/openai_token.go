package hooks

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/openai/openai-go"
	"github.com/pkoukk/tiktoken-go"
	"github.com/samber/do/v2"
	log "github.com/sirupsen/logrus"

	"github.com/modelgate/modelgate/internal/runtime/core"
)

var tokenEncoderMap = make(map[string]*tiktoken.Tiktoken)
var locker sync.RWMutex

// OpenAITokenHook 计算 Token
type OpenAITokenHook struct {
}

var _ core.Hook = (*OpenAITokenHook)(nil)

func NewOpenAITokenHook(i do.Injector) (*OpenAITokenHook, error) {
	return &OpenAITokenHook{}, nil
}

func (h *OpenAITokenHook) Name() string {
	return "openai_token"
}

// Before 执行前
func (h *OpenAITokenHook) Before(ctx context.Context, c *core.Context) (err error) {
	if c.CurrentModel == nil {
		err = errors.New("model info is nil")
		return
	}
	var reqBody struct {
		Messages []openai.ChatCompletionMessage `json:"messages"`
	}
	if err = json.Unmarshal(c.InputBody, &reqBody); err != nil {
		return
	}
	var text strings.Builder
	for _, message := range reqBody.Messages {
		text.WriteString(message.Content)
	}
	// providerId := c.ModelInfo.ProviderId
	tokenNum, err := h.countTokenText(c.CurrentModel.ModelCode, text.String())
	if err != nil {
		return
	}

	c.PromptTokens = tokenNum
	log.Info("prompt token num: ", tokenNum)
	return
}

// After 执行后
func (h *OpenAITokenHook) After(ctx context.Context, c *core.Context) (err error) {
	if c.IsStream || c.Usage != nil {
		return
	}
	var respData struct {
		Choices []openai.ChatCompletionChoice `json:"choices"`
	}
	if err = json.Unmarshal(c.RawResponse, &respData); err != nil {
		return
	}
	var text strings.Builder
	for _, item := range respData.Choices {
		text.WriteString(item.Message.Content)
	}
	tokenNum, err := h.countTokenText(c.CurrentModel.ModelCode, text.String())
	if err != nil {
		return
	}
	log.Info("token num: ", tokenNum)
	return
}

func (h *OpenAITokenHook) OnChunk(ctx context.Context, c *core.Context, chunk *core.StreamChunk) (err error) {
	if chunk.Finish {
		return
	}
	// 如果有的话，解析返回的usage
	var respData openai.ChatCompletionChunk
	if err = json.Unmarshal([]byte(chunk.Data), &respData); err != nil {
		log.Errorf("json unmarshal data %s, error: %v", chunk.Data, err)
		return
	}
	c.ActualModel = respData.Model
	// 计算token
	if len(respData.Choices) > 0 {
		var tokenNum int
		tokenNum, err = h.countTokenText(c.CurrentModel.ModelCode, respData.Choices[0].Delta.Content)
		if err != nil {
			return
		}
		c.CompletionTokens += tokenNum
	}
	if respData.Usage.TotalTokens > 0 {
		c.Usage = &core.Usage{
			PromptTokens:     respData.Usage.PromptTokens,
			CompletionTokens: respData.Usage.CompletionTokens,
			TotalTokens:      respData.Usage.TotalTokens,
		}
	}
	return
}

func (h *OpenAITokenHook) OnError(ctx context.Context, c *core.Context, err error) {
}

func (h *OpenAITokenHook) countTokenText(model, text string) (num int, err error) {
	if len(text) == 0 {
		return
	}
	tt, err := h.getTokenEncoder(model)
	if err != nil {
		return
	}
	num = len(tt.Encode(text, nil, nil))
	return
}

func (h *OpenAITokenHook) getTokenEncoder(model string) (*tiktoken.Tiktoken, error) {
	locker.RLock()
	tt, ok := tokenEncoderMap[model]
	locker.RUnlock()
	if ok {
		return tt, nil
	}
	var modelName string
	if strings.HasPrefix(model, "gpt-3.5") {
		modelName = "gpt-3.5-turbo"
	} else if strings.HasPrefix(model, "gpt-4o") {
		modelName = "gpt-4o"
	} else if strings.HasPrefix(model, "gpt-4") {
		modelName = "gpt-4"
	} else {
		modelName = "gpt-3.5-turbo"
	}
	locker.Lock()
	defer locker.Unlock()
	tt, err := tiktoken.EncodingForModel(modelName)
	if err != nil {
		return nil, err
	}
	tokenEncoderMap[model] = tt
	return tt, nil
}
