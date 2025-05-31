# 📰 News Reporter

OpenAI Responses APIのweb_search機能を使用したニュース検索アプリケーション

## ✨ 機能

- OpenAI Responses APIのWeb検索機能を使用
- リアルタイムでのWeb検索とAI要約
- コマンドライン インターフェース
- ストリーミングレスポンス対応
- 日本語対応

## 🚀 セットアップ

### 1. リポジトリのクローン
```bash
git clone <repository-url>
cd news_reporter
```

### 2. 依存関係のインストール
```bash
go mod tidy
```

### 3. 環境変数の設定
```bash
export OPENAI_API_KEY="your-openai-api-key-here"
```

または、`.env`ファイルを作成：
```bash
# .env
OPENAI_API_KEY=your-openai-api-key-here
OPENAI_BASE_URL=https://api.openai.com/v1  # オプション
```

## 📖 使用方法

### 基本的な使用法
```bash
go run main.go "検索クエリ"
```

### 使用例
```bash
# 最新のAI技術ニュースを検索
go run main.go "最新のAI技術ニュース"

# 日本の経済動向を検索
go run main.go "日本の経済動向 2024"

# 複数のキーワードで検索
go run main.go "Python プログラミング 最新トレンド"
```

### ヘルプの表示
```bash
go run main.go --help
```

## 🏗️ アーキテクチャ

```
news_reporter/
├── main.go           # メインアプリケーション
├── config/
│   └── config.go     # 設定管理
├── client/
│   └── openai.go     # OpenAI API クライアント
├── handlers/
│   └── search.go     # 検索ハンドラー
├── models/
│   └── response.go   # データ構造体
├── go.mod
└── go.sum
```

## 🔧 技術仕様

- **言語**: Go 1.21+
- **API**: OpenAI Responses API
- **機能**: Web Search, Streaming Response
- **依存関係**: 
  - `github.com/joho/godotenv` - 環境変数管理

## 📋 出力例

```
🔍 検索中: 最新のAI技術ニュース
--------------------------------------------------
📊 検索結果 (2024-01-15 10:30:45)
==================================================

🌐 Web検索結果 (5件):
------------------------------

1. 2024年のAI技術トレンドと展望
   🔗 https://example.com/ai-trends-2024
   📄 2024年のAI技術における最新動向と今後の展望について...

2. ChatGPTの新機能発表
   🔗 https://example.com/chatgpt-new-features
   📄 OpenAIが発表したChatGPTの新機能について詳細解説...

🤖 AI要約:
------------------------------
2024年のAI技術は大きな進歩を遂げており、特に自然言語処理と
画像生成の分野で注目すべき発展が見られます...

==================================================
✅ 検索が完了しました！
```

## 🔐 環境変数

| 変数名 | 必須 | 説明 | デフォルト |
|--------|------|------|------------|
| `OPENAI_API_KEY` | ✅ | OpenAI APIキー | - |
| `OPENAI_BASE_URL` | ❌ | OpenAI API Base URL | `https://api.openai.com/v1` |

## 🛠️ 今後の拡張予定

- [ ] Webインターフェース（REST API）
- [ ] 検索結果の保存機能（JSON/CSV出力）
- [ ] フィルタリング機能（日付、ソース等）
- [ ] 設定ファイル対応
- [ ] ログ機能
- [ ] バッチ検索モード

## 📝 ライセンス

MIT License

## 🤝 コントリビューション

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request