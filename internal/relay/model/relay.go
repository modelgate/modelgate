package model

import (
	"net/http"

	"github.com/modelgate/modelgate/internal/runtime/core"
)

// RelayRequest 转发请求
type RelayRequest struct {
	Path         string
	Provider     string
	Model        string
	ApiKeyId     int64
	AccountId    int64
	IsStream     bool
	StreamWriter core.StreamWriter
	Request      *http.Request
	InputBody    []byte
}

// RelayResponse 转发响应
type RelayResponse struct {
	HTTPResponse *http.Response
	RawResponse  []byte
}

type RelayInfo struct {
	ProviderCount int64
	ModelCount    int64
	ApiKeyCount   int64
}
