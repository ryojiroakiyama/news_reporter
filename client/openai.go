package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"news_reporter/config"
	"news_reporter/models"
)

type OpenAIClient struct {
	config     *config.Config
	httpClient *http.Client
}

// NewOpenAIClient 新しいOpenAIクライアントを作成
func NewOpenAIClient(cfg *config.Config) *OpenAIClient {
	return &OpenAIClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Search Web検索を実行
func (c *OpenAIClient) Search(query string) (*models.SearchResult, error) {
	// 現在の日付を取得
	currentDate := time.Now().Format("2006年1月2日")

	// システムメッセージで最新情報検索を指示
	systemMessage := fmt.Sprintf(`あなたは最新のニュースと情報を検索するアシスタントです。
現在の日付: %s

以下の指示に従ってください：
1. 必ずweb_search_previewツールを使用して、最新の情報を検索してください
2. 検索結果から、今日（%s）またはできるだけ最近の情報を優先してください
3. 古い情報（1週間以上前）は避け、最新のニュースに焦点を当ててください
4. 検索結果を日本語で要約し、情報源のURLも含めてください
5. 情報の日付が明確でない場合は、その旨を明記してください`, currentDate, currentDate)

	// ユーザークエリを現在の日付と組み合わせて強化
	enhancedQuery := fmt.Sprintf("【%s時点】%s（最新情報・今日のニュース）", currentDate, query)

	// リクエストボディを構築
	request := models.ResponseRequest{
		Model: "gpt-4o-mini",
		Input: []models.InputItem{
			{
				Type:    "message",
				Role:    "system",
				Content: systemMessage,
			},
			{
				Type:    "message",
				Role:    "user",
				Content: enhancedQuery,
			},
		},
		Tools: []models.Tool{
			{
				Type: "web_search_preview",
			},
		},
		ToolChoice:  "required",
		Stream:      true,
		Temperature: 0.3, // より一貫性のある結果のために温度を下げる
	}

	// JSONエンコード
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", c.config.BaseURL+"/responses", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.OpenAIAPIKey)
	req.Header.Set("Accept", "text/event-stream")

	// リクエストを送信
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// ストリーミングレスポンスを処理
	return c.processStreamResponse(resp.Body, query)
}

// processStreamResponse ストリーミングレスポンスを処理
func (c *OpenAIClient) processStreamResponse(body io.Reader, query string) (*models.SearchResult, error) {
	scanner := bufio.NewScanner(body)
	result := &models.SearchResult{
		Query:     query,
		Results:   make([]models.WebSearchResult, 0),
		Timestamp: time.Now(),
	}

	var responseContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Server-Sent Eventsの形式をパース
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// [DONE]は終了シグナル
			if data == "[DONE]" {
				break
			}

			// JSONをパース
			var event map[string]interface{}
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				// パースエラーは無視して続行
				continue
			}

			// イベントタイプに基づいて処理
			eventType, ok := event["type"].(string)
			if !ok {
				continue
			}

			switch eventType {
			case "response.output_text.delta":
				// テキストデルタを処理
				if delta, ok := event["delta"].(string); ok {
					responseContent.WriteString(delta)
				}
			case "response.output_text.annotation.added":
				// アノテーションを処理（Web検索結果など）
				if err := c.processAnnotation(event, result); err != nil {
					fmt.Printf("Warning: failed to process annotation: %v\n", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading stream: %w", err)
	}

	// 要約を設定
	result.Summary = responseContent.String()

	return result, nil
}

// processAnnotation アノテーションを処理してWeb検索結果を抽出
func (c *OpenAIClient) processAnnotation(event map[string]interface{}, result *models.SearchResult) error {
	annotation, ok := event["annotation"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid annotation data")
	}

	// URL引用の場合のみ処理
	annotationType, ok := annotation["type"].(string)
	if !ok || annotationType != "url_citation" {
		return nil
	}

	// Web検索結果を作成
	searchResult := models.WebSearchResult{}

	if title, ok := annotation["title"].(string); ok {
		searchResult.Title = title
	}
	if url, ok := annotation["url"].(string); ok {
		searchResult.URL = url
	}

	// スニペットは含まれていない可能性があるため、タイトルを使用
	if searchResult.Title != "" {
		searchResult.Snippet = "Web検索結果から引用"
	}

	// 重複チェック
	for _, existing := range result.Results {
		if existing.URL == searchResult.URL {
			return nil // 重複は追加しない
		}
	}

	result.Results = append(result.Results, searchResult)
	return nil
}
