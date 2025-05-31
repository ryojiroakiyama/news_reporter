package main

import (
	"fmt"
	"os"
	"strings"

	"news_reporter/client"
	"news_reporter/config"
	"news_reporter/handlers"
)

func showHelp() {
	fmt.Println("📰 News Reporter - 最新ニュース検索アプリ")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("OpenAI Responses API と web_search_preview を使用して")
	fmt.Println("リアルタイムで最新のニュースや情報を検索します。")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  go run main.go \"検索クエリ\"")
	fmt.Println("  go run main.go [オプション]")
	fmt.Println("")
	fmt.Println("例:")
	fmt.Println("  go run main.go \"今日の経済ニュース\"")
	fmt.Println("  go run main.go \"最新のAI技術動向\"")
	fmt.Println("  go run main.go \"円安ドル高の最新状況\"")
	fmt.Println("")
	fmt.Println("オプション:")
	fmt.Println("  -h, --help    このヘルプメッセージを表示")
	fmt.Println("")
	fmt.Println("機能:")
	fmt.Println("  ✅ リアルタイムWeb検索")
	fmt.Println("  ✅ 最新情報の自動取得")
	fmt.Println("  ✅ 日本語での要約表示")
	fmt.Println("  ✅ 情報源URL付きの結果")
	fmt.Println("")
	fmt.Println("注意: OPENAI_API_KEY環境変数の設定が必要です")
}

func main() {
	// 使用方法を表示する関数
	showUsage := showHelp

	// コマンドライン引数をチェック
	if len(os.Args) < 2 {
		showUsage()
		fmt.Println("❌ エラー: 検索クエリが指定されていません")
		os.Exit(1)
	}

	// ヘルプオプション
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		showUsage()
		os.Exit(0)
	}

	// 検索クエリを取得（複数の引数を結合）
	query := strings.Join(os.Args[1:], " ")
	if strings.TrimSpace(query) == "" {
		showUsage()
		fmt.Println("❌ エラー: 空の検索クエリです")
		os.Exit(1)
	}

	// 設定を読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ 設定エラー: %v\n", err)
		fmt.Println()
		fmt.Println("💡 ヒント: OPENAI_API_KEY環境変数を設定してください")
		fmt.Println("   export OPENAI_API_KEY=\"your-api-key-here\"")
		os.Exit(1)
	}

	// OpenAIクライアントを初期化
	openaiClient := client.NewOpenAIClient(cfg)

	// 検索ハンドラーを初期化
	searchHandler := handlers.NewSearchHandler(openaiClient)

	// 検索を実行
	if err := searchHandler.HandleSearch(query); err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ 検索が完了しました！")
}
