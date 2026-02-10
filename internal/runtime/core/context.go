package core

import (
	"net/http"

	"github.com/modelgate/modelgate/pkg/utils"
)

// Context 上下文
type Context struct {
	AccountId       int64
	AccountApiKeyId int64

	RequestUUID  utils.UUIDv7
	AttemptNo    int
	RequestId    int64
	UrlPath      string
	ProviderCode string
	ModelCode    string
	CurrentModel *Model // 模型

	PromptTokens     int // 输入Token数
	CompletionTokens int // 输出Token数

	Usage       *Usage //  提供商返回模型使用情况
	ActualModel string // 厂商返回的模型

	PreCost   int64 // 预先扣费
	TotalCost int64 // 实际扣费

	Header    http.Header
	InputBody []byte // 统一输入

	// HTTP
	HTTPRequest  *http.Request
	HTTPResponse *http.Response
	RawResponse  []byte

	// 流
	IsStream     bool
	StreamWriter StreamWriter

	LastErr error
}
