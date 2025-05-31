package handlers

import (
	"fmt"
	"strings"
	"time"

	"news_reporter/client"
	"news_reporter/models"
)

type SearchHandler struct {
	openaiClient *client.OpenAIClient
}

// NewSearchHandler 新しい検索ハンドラーを作成
func NewSearchHandler(openaiClient *client.OpenAIClient) *SearchHandler {
	return &SearchHandler{
		openaiClient: openaiClient,
	}
}

// HandleSearch 検索を処理
func (h *SearchHandler) HandleSearch(query string) error {
	currentDate := time.Now().Format("2006年1月2日 15:04")
	fmt.Printf("🔍 最新情報を検索中: %s (%s時点)\n", query, currentDate)
	fmt.Println(strings.Repeat("-", 50))

	// 検索を実行
	result, err := h.openaiClient.Search(query)
	if err != nil {
		return fmt.Errorf("検索に失敗しました: %w", err)
	}

	// 結果を表示
	h.displayResult(result)

	return nil
}

// displayResult 検索結果を表示
func (h *SearchHandler) displayResult(result *models.SearchResult) {
	fmt.Printf("📊 最新検索結果 (%s取得)\n", result.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("=", 50))

	// Web検索結果を表示
	if len(result.Results) > 0 {
		fmt.Printf("\n🌐 最新Web検索結果 (%d件):\n", len(result.Results))
		fmt.Println(strings.Repeat("-", 30))

		for i, searchResult := range result.Results {
			fmt.Printf("\n%d. %s\n", i+1, searchResult.Title)
			fmt.Printf("   🔗 %s\n", searchResult.URL)
			if searchResult.Snippet != "" {
				// スニペットを適切な長さで改行
				snippet := h.formatSnippet(searchResult.Snippet, 80)
				fmt.Printf("   📄 %s\n", snippet)
			}
		}
	} else {
		fmt.Println("\n⚠️  最新のWeb検索結果が見つかりませんでした")
	}

	// AI要約を表示
	if result.Summary != "" {
		fmt.Printf("\n🤖 最新情報AI要約:\n")
		fmt.Println(strings.Repeat("-", 30))
		summary := h.formatText(result.Summary, 80)
		fmt.Printf("%s\n", summary)
	}

	fmt.Println(strings.Repeat("=", 50))
}

// formatSnippet スニペットをフォーマット
func (h *SearchHandler) formatSnippet(snippet string, maxWidth int) string {
	if len(snippet) <= maxWidth {
		return snippet
	}

	// 長いスニペットを適切に切り詰める
	words := strings.Fields(snippet)
	var result strings.Builder
	currentLength := 0

	for _, word := range words {
		if currentLength+len(word)+1 > maxWidth-3 { // "..."を考慮
			result.WriteString("...")
			break
		}

		if currentLength > 0 {
			result.WriteString(" ")
			currentLength++
		}

		result.WriteString(word)
		currentLength += len(word)
	}

	return result.String()
}

// formatText テキストを指定幅でフォーマット
func (h *SearchHandler) formatText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(line) <= maxWidth {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		// 長い行を分割
		words := strings.Fields(line)
		currentLength := 0
		var currentLine strings.Builder

		for _, word := range words {
			if currentLength+len(word)+1 > maxWidth && currentLength > 0 {
				result.WriteString(currentLine.String())
				result.WriteString("\n")
				currentLine.Reset()
				currentLength = 0
			}

			if currentLength > 0 {
				currentLine.WriteString(" ")
				currentLength++
			}

			currentLine.WriteString(word)
			currentLength += len(word)
		}

		if currentLine.Len() > 0 {
			result.WriteString(currentLine.String())
			result.WriteString("\n")
		}
	}

	return strings.TrimSuffix(result.String(), "\n")
}
