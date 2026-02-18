package core

import (
	"net/http"
	"sync"

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
	IsAnthropic  bool
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

// Pool Context 对象池
var Pool = sync.Pool{
	New: func() any {
		return new(Context)
	},
}

// Get 从池中获取 Context
func Get() *Context {
	return Pool.Get().(*Context)
}

// Put 将 Context 归还到池中
func Put(ctx *Context) {
	ctx.Reset()
	Pool.Put(ctx)
}

// Reset 重置 Context 状态
func (ctx *Context) Reset() {
	ctx.AccountId = 0
	ctx.AccountApiKeyId = 0
	ctx.RequestUUID = nil
	ctx.AttemptNo = 0
	ctx.RequestId = 0
	ctx.UrlPath = ""
	ctx.ProviderCode = ""
	ctx.ModelCode = ""
	ctx.CurrentModel = nil
	ctx.PromptTokens = 0
	ctx.CompletionTokens = 0
	ctx.Usage = nil
	ctx.ActualModel = ""
	ctx.PreCost = 0
	ctx.TotalCost = 0
	ctx.Header = nil
	ctx.InputBody = nil
	ctx.HTTPRequest = nil
	ctx.HTTPResponse = nil
	ctx.RawResponse = nil
	ctx.IsStream = false
	ctx.StreamWriter = nil
	ctx.LastErr = nil
}
