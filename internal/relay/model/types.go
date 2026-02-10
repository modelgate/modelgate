package model

const (
	TableAccount          = "accounts"
	TableAccountApiKey    = "account_api_keys"
	TableRequest          = "requests"
	TableRequestStat      = "request_stats"
	TableRequestAttempt   = "request_attempts"
	TableLedger           = "ledgers"
	TableProvider         = "providers"
	TableProviderApiKey   = "provider_api_keys"
	TableModelPricing     = "model_pricings"
	TableModel            = "models"
	TableRelayStat        = "relay_stats"
	TableRelayUsage       = "relay_usages"
	TableRelayHourlyUsage = "relay_hourly_usages"
)

const (
	ProviderCodeOpenAI      string = "openai"       // OpenAI
	ProviderCodeAnthropic   string = "anthropic"    // Anthropic
	ProviderCodeDeepSeek    string = "deepseek"     // DeepSeek
	ProviderCodeZhipu       string = "zhipu"        // 智谱 - bigmodel
	ProviderCodeZhipuClaude string = "zhipu-claude" // 智谱 - claude
)

var AllProviderCodeList = []string{
	ProviderCodeOpenAI,
	ProviderCodeAnthropic,
	ProviderCodeDeepSeek,
	ProviderCodeZhipu,
	ProviderCodeZhipuClaude,
}

// Currency 货币单位
type Currency string

const (
	CurrencyUSD   Currency = "USD"   // 美元
	CurrencyCNY   Currency = "CNY"   // 人民币
	CurrencyPOINT Currency = "POINT" // 点数
)

// ApiKeyStatus API密钥状态
type ApiKeyStatus string

const (
	ApiKeyStatusEnabled  ApiKeyStatus = "enabled"  // 正常
	ApiKeyStatusDisabled ApiKeyStatus = "disabled" // 禁用
	ApiKeyStatusCooldown ApiKeyStatus = "cooldown" // 冷却
	ApiKeyStatusRevoked  ApiKeyStatus = "revoked"  // 撤销
)

// EnableStatus 启用状态
type EnableStatus string

const (
	EnableStatusEnabled  EnableStatus = "enabled"  // 正常
	EnableStatusDisabled EnableStatus = "disabled" // 禁用
)

const (
	WorkerKeyLeader = "relay:worker:leader"
)

const (
	UsageProviderPrefix       = "usage:provider:"
	UsageAccountApiKeyPrefix  = "usage:account_api_key:"
	UsageProviderApiKeyPrefix = "usage:provider_api_key:"
)

var AllUsagePrefixs = []string{
	UsageProviderPrefix,
	UsageAccountApiKeyPrefix,
	UsageProviderApiKeyPrefix,
}

const (
	MetricTotal   = "total"
	MetricSuccess = "success"
	MetricFailed  = "failed"
	MetricUsage   = "usage"
)
