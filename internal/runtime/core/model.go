package core

const (
	ProviderCodeAnthropic   string = "anthropic"    // Anthropic
	ProviderCodeDeepSeek    string = "deepseek"     // DeepSeek
	ProviderCodeOpenAI      string = "openai"       // OpenAI
	ProviderCodeZhipu       string = "zhipu"        // 智谱 - bigmodel
	ProviderCodeZhipuClaude string = "zhipu-claude" // 智谱 - claude
)

var AllProviderCodeList = []string{
	ProviderCodeAnthropic,
	ProviderCodeDeepSeek,
	ProviderCodeOpenAI,
	ProviderCodeZhipu,
	ProviderCodeZhipuClaude,
}

// Model 模型
type Model struct {
	ModelId         int64  // ID
	ModelCode       string // 模型Code
	ProviderId      int64  // 提供商 ID
	ProviderCode    string // 提供商Code
	BaseUrl         string // 基础URL
	ApiKeyId        int64  // 提供商 API Key ID
	ApiKeyEncrypted string // 加密后的 API Key

	InputPrice        float64 // 输入价格
	InputCachePrice   float64 // 输入缓存价格
	OutputPrice       float64 // 输出价格
	TokenNum          int64   // Token 数量
	PointsPerCurrency int64   // 每个货币点数
}

// Usage 使用情况
type Usage struct {
	PromptTokens       int64
	PromptCachedTokens int64
	CompletionTokens   int64
	TotalTokens        int64
}
