package models

import "time"

// ResponseRequest Responses APIへのリクエスト構造体
type ResponseRequest struct {
	Model       string      `json:"model"`
	Input       []InputItem `json:"input"`
	Tools       []Tool      `json:"tools,omitempty"`
	ToolChoice  string      `json:"tool_choice,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
	Temperature float64     `json:"temperature,omitempty"`
}

// InputItem 入力アイテム
type InputItem struct {
	Type    string `json:"type"`
	Role    string `json:"role,omitempty"`
	Content string `json:"content"`
}

// Tool ツール定義
type Tool struct {
	Type string `json:"type"`
}

// FunctionTool 関数ツール
type FunctionTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// WebSearchOptions Web検索オプション
type WebSearchOptions struct{}

// ResponseData レスポンスデータ
type ResponseData struct {
	ID       string            `json:"id"`
	Object   string            `json:"object"`
	Created  int64             `json:"created"`
	Model    string            `json:"model"`
	Choices  []Choice          `json:"choices"`
	Usage    *Usage            `json:"usage,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Choice レスポンスの選択肢
type Choice struct {
	Index        int         `json:"index"`
	Message      *Message    `json:"message,omitempty"`
	Delta        *Message    `json:"delta,omitempty"`
	FinishReason *string     `json:"finish_reason,omitempty"`
	Logprobs     interface{} `json:"logprobs,omitempty"`
}

// Message メッセージ
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall ツール呼び出し
type ToolCall struct {
	ID        string               `json:"id"`
	Type      string               `json:"type"`
	Function  *FunctionCallDetail  `json:"function,omitempty"`
	WebSearch *WebSearchCallDetail `json:"web_search,omitempty"`
}

// FunctionCallDetail 関数呼び出し詳細
type FunctionCallDetail struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// WebSearchCallDetail Web検索呼び出し詳細
type WebSearchCallDetail struct {
	Query   string            `json:"query"`
	Results []WebSearchResult `json:"results,omitempty"`
}

// WebSearchResult Web検索結果
type WebSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// Usage 使用量情報
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamEvent ストリーミングイベント
type StreamEvent struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}

// SearchResult 検索結果の統合表現
type SearchResult struct {
	Query     string            `json:"query"`
	Results   []WebSearchResult `json:"results"`
	Summary   string            `json:"summary,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}
