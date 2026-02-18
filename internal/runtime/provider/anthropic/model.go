package anthropic

const (
	MessageStart      = "message_start"
	MessageDelta      = "message_delta"
	MessageStop       = "message_stop"
	ContentBlockStart = "content_block_start"
	ContentBlockDelta = "content_block_delta"
	ContentBlockStop  = "content_block_stop"
)

type Response struct {
	Id           string    `json:"id"`
	Type         string    `json:"type"`
	Role         string    `json:"role"`
	Content      []Content `json:"content"`
	Model        string    `json:"model"`
	StopReason   *string   `json:"stop_reason"`
	StopSequence *string   `json:"stop_sequence"`
	Usage        Usage     `json:"usage"`
	Error        Error     `json:"error"`
}

type StreamResponse struct {
	Type         string    `json:"type"`
	Message      *Response `json:"message"`
	Index        int       `json:"index"`
	ContentBlock *Content  `json:"content_block"`
	Delta        *Delta    `json:"delta"`
	Usage        *Usage    `json:"usage"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Content struct {
	Type   string       `json:"type"`
	Text   string       `json:"text,omitempty"`
	Source *ImageSource `json:"source,omitempty"`
	// tool_calls
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Input     any    `json:"input,omitempty"`
	Content   string `json:"content,omitempty"`
	ToolUseId string `json:"tool_use_id,omitempty"`
}

type ImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// Billing and rate-limit usage.

// Anthropic's API bills and rate-limits by token counts, as tokens represent the underlying cost to our systems.

// Under the hood, the API transforms requests into a format suitable for the model. The model's output then goes through a parsing stage before becoming an API response.
// As a result, the token counts in usage will not match one-to-one with the exact visible content of an API request or response.

// For example, output_tokens will be non-zero, even for an empty string response from Claude.

// Total input tokens in a request is the summation of input_tokens, cache_creation_input_tokens, and cache_read_input_tokens.

type Usage struct {
	CacheCreation            CacheCreation `json:"cache_creation"`
	CacheCreationInputTokens int64         `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64         `json:"cache_read_input_tokens"`
	InputTokens              int64         `json:"input_tokens"`
	OutputTokens             int64         `json:"output_tokens"`
	ServerToolUse            ServerToolUse `json:"server_tool_use"`
	ServiceTier              string        `json:"service_tier"`
}

type CacheCreation struct {
	Ephemeral1hInputTokens int64 `json:"ephemeral_1h_input_tokens"`
	Ephemeral5mInputTokens int64 `json:"ephemeral_5m_input_tokens"`
}

type ServerToolUse struct {
	WebSearchRequests int64 `json:"web_search_requests"`
}

type Delta struct {
	Type         string  `json:"type"`
	Text         string  `json:"text"`
	PartialJson  string  `json:"partial_json,omitempty"`
	StopReason   *string `json:"stop_reason"`
	StopSequence *string `json:"stop_sequence"`
}
