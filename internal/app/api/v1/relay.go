package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/samber/lo"

	"github.com/modelgate/modelgate/internal/relay"
	"github.com/modelgate/modelgate/internal/runtime"
	"github.com/modelgate/modelgate/internal/runtime/core"
	"github.com/modelgate/modelgate/pkg/common"
	"github.com/modelgate/modelgate/pkg/utils"
)

type RelayService struct {
	relayService relay.Service
	runtime      *runtime.Runtime
}

func NewRelayService(i do.Injector) (*RelayService, error) {
	return &RelayService{
		relayService: do.MustInvoke[relay.Service](i),
		runtime:      do.MustInvoke[*runtime.Runtime](i),
	}, nil
}

func (s *RelayService) Run(c *gin.Context) {
	err := s.run(c, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (s *RelayService) RunWithProvider(c *gin.Context) {
	err := s.run(c, c.Param("provider"), c.Param("path"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (s *RelayService) run(c *gin.Context, relayProvider, relayPath string) (err error) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	defer c.Request.Body.Close()
	providerCode, modelCode, stream, inputData, err := s.parseInputBody(data)
	if err != nil {
		return
	}
	providerCode = lo.Ternary(relayProvider != "", relayProvider, providerCode)
	currentModel, err := s.relayService.ResolveModel(c, providerCode, modelCode)
	if err != nil {
		return
	}
	cModel := &core.Model{
		ModelId:           currentModel.ModelId,
		ModelCode:         currentModel.ModelCode,
		ProviderId:        currentModel.ProviderId,
		ProviderCode:      currentModel.ProviderCode,
		BaseUrl:           currentModel.BaseUrl,
		ApiKeyId:          currentModel.ApiKeyId,
		ApiKeyEncrypted:   currentModel.ApiKeyEncrypted,
		InputPrice:        currentModel.InputPrice,
		InputCachePrice:   currentModel.InputCachePrice,
		OutputPrice:       currentModel.OutputPrice,
		TokenNum:          currentModel.TokenNum,
		PointsPerCurrency: currentModel.PointsPerCurrency,
	}
	rCtx := core.Get()
	defer core.Put(rCtx)
	rCtx.RequestUUID = utils.NewUUIDv7()
	rCtx.CurrentModel = cModel
	rCtx.AccountApiKeyId = common.GetApiKeyId(c)
	rCtx.AccountId = common.GetAccountId(c)
	rCtx.UrlPath = lo.Ternary(relayPath != "", relayPath, c.Request.URL.Path)
	rCtx.InputBody = inputData
	rCtx.Header = c.Request.Header
	if stream {
		rCtx.IsStream = true
		rCtx.StreamWriter = newGinSSEWriter(c)
	}
	if err = s.runtime.Run(c, rCtx); err != nil {
		return
	}
	if stream {
		return
	}
	for k, v := range rCtx.HTTPResponse.Header {
		c.Writer.Header().Set(k, v[0])
	}
	c.Writer.WriteHeader(rCtx.HTTPResponse.StatusCode)
	io.Copy(c.Writer, bytes.NewReader(rCtx.RawResponse))
	return
}

func (s *RelayService) parseInputBody(data []byte) (providerCode, modelCode string, stream bool, inputData []byte, err error) {
	reqBody := make(map[string]any)
	if err = json.Unmarshal(data, &reqBody); err != nil {
		return
	}
	providerCode, _ = reqBody["provider"].(string)
	modelCode, _ = reqBody["model"].(string)
	stream, _ = reqBody["stream"].(bool)
	// provider 不是标准请求参数，去除
	if providerCode != "" {
		delete(reqBody, "provider")
	}
	inputData, err = json.Marshal(reqBody)
	return
}

type GinSSEWriter struct {
	w       gin.ResponseWriter
	flusher http.Flusher
}

func newGinSSEWriter(c *gin.Context) *GinSSEWriter {
	return &GinSSEWriter{
		w:       c.Writer,
		flusher: c.Writer.(http.Flusher),
	}
}

func (g *GinSSEWriter) Open() error {
	g.w.Header().Set("Content-Type", "text/event-stream")
	g.w.Header().Set("Cache-Control", "no-cache")
	g.w.Header().Set("Connection", "keep-alive")
	g.w.Header().Set("X-Accel-Buffering", "no") // Nginx
	g.w.WriteHeader(http.StatusOK)
	return nil
}

// Write 写入
func (g *GinSSEWriter) Write(chunk *core.StreamChunk) error {
	if chunk.Finish {
		_, err := g.w.Write([]byte("data: [DONE]\n\n"))
		g.flusher.Flush()
		return err
	}

	buf := bytes.NewBuffer([]byte("data: "))
	buf.Write([]byte(chunk.Data))
	buf.Write([]byte("\n\n"))

	_, err := g.w.Write(buf.Bytes())
	g.flusher.Flush()
	return err
}

// Close 关闭
func (g *GinSSEWriter) Close() error {
	return nil
}
