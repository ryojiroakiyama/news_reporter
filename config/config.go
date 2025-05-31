package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
	BaseURL      string
}

// LoadConfig 環境変数から設定を読み込む
func LoadConfig() (*Config, error) {
	// .envファイルが存在する場合は読み込む（オプション）
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &Config{
		OpenAIAPIKey: apiKey,
		BaseURL:      baseURL,
	}, nil
}
