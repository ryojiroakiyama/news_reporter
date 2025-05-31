package main

import (
	"fmt"
	"os"
	"strings"

	"news_reporter/audio"
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
	fmt.Println("  go run main.go [オプション] \"検索クエリ\"")
	fmt.Println("")
	fmt.Println("例:")
	fmt.Println("  go run main.go \"今日の経済ニュース\"")
	fmt.Println("  go run main.go \"最新のAI技術動向\"")
	fmt.Println("  go run main.go \"円安ドル高の最新状況\"")
	fmt.Println("  go run main.go --audio \"今日のニュース\"")
	fmt.Println("  go run main.go --save summary.mp3 \"AIニュース\"")
	fmt.Println("")
	fmt.Println("オプション:")
	fmt.Println("  -h, --help                このヘルプメッセージを表示")
	fmt.Println("  -a, --audio               音声再生機能付きで実行")
	fmt.Println("  -s, --save <filename>     要約を音声ファイルに保存")
	fmt.Println("")
	fmt.Println("機能:")
	fmt.Println("  ✅ リアルタイムWeb検索")
	fmt.Println("  ✅ 最新情報の自動取得")
	fmt.Println("  ✅ 日本語での要約表示")
	fmt.Println("  ✅ 情報源URL付きの結果")
	fmt.Println("  🎵 音声読み上げ機能")
	fmt.Println("  💾 音声ファイル保存機能")
	fmt.Println("")
	fmt.Println("音声機能について:")
	fmt.Println("  • OpenAI TTSを使用した高品質な音声合成")
	fmt.Println("  • 日本語要約の自動読み上げ")
	fmt.Println("  • MP3形式での音声ファイル保存")
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

	// フラグとクエリの解析
	var audioMode bool
	var saveMode bool
	var saveFilename string
	var query string
	var args []string

	// 引数を解析
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "--help", "-h":
			showUsage()
			os.Exit(0)
		case "--audio", "-a":
			audioMode = true
		case "--save", "-s":
			saveMode = true
			if i+1 < len(os.Args) {
				saveFilename = os.Args[i+1]
				i++ // 次の引数をスキップ
			} else {
				fmt.Println("❌ エラー: --save オプションにはファイル名が必要です")
				os.Exit(1)
			}
		default:
			args = append(args, arg)
		}
	}

	// クエリを結合
	query = strings.Join(args, " ")
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

	// TTSクライアントを初期化
	ttsClient := audio.NewTTSClient(cfg)

	// 検索ハンドラーを初期化
	searchHandler := handlers.NewSearchHandler(openaiClient, ttsClient)

	// モードに応じて実行
	if saveMode {
		// 音声ファイル保存モード
		if err := searchHandler.SaveAudioSummary(query, saveFilename); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	} else if audioMode {
		// 音声再生モード
		if err := searchHandler.HandleSearchWithAudio(query); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	} else {
		// 通常の検索モード
		if err := searchHandler.HandleSearch(query); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n✅ 処理が完了しました！")
}
