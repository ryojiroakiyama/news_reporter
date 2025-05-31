package audio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"

	"news_reporter/config"
)

type TTSClient struct {
	config     *config.Config
	httpClient *http.Client
}

type TTSRequest struct {
	Model  string `json:"model"`
	Input  string `json:"input"`
	Voice  string `json:"voice"`
	Format string `json:"response_format"`
}

// NewTTSClient 新しいTTSクライアントを作成
func NewTTSClient(cfg *config.Config) *TTSClient {
	return &TTSClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // 音声生成は時間がかかる場合があるため長めに設定
		},
	}
}

// SynthesizeAndPlay テキストを音声に変換して再生
func (t *TTSClient) SynthesizeAndPlay(text string) error {
	fmt.Println("🎵 音声を生成中...")

	// 音声データを生成
	audioData, err := t.synthesize(text)
	if err != nil {
		return fmt.Errorf("音声生成に失敗しました: %w", err)
	}

	fmt.Println("🔊 音声を再生中...")

	// 音声を再生
	if err := t.playAudio(audioData); err != nil {
		return fmt.Errorf("音声再生に失敗しました: %w", err)
	}

	return nil
}

// synthesize OpenAI TTS APIを使用してテキストを音声に変換
func (t *TTSClient) synthesize(text string) ([]byte, error) {
	// リクエストボディを構築
	request := TTSRequest{
		Model:  "tts-1", // 高速モデル（tts-1-hdもあります）
		Input:  text,
		Voice:  "alloy", // 利用可能な声: alloy, echo, fable, onyx, nova, shimmer
		Format: "mp3",
	}

	// JSONエンコード
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", t.config.BaseURL+"/audio/speech", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.config.OpenAIAPIKey)

	// リクエストを送信
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 音声データを読み取り
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	return audioData, nil
}

// playAudio MP3音声データを再生
func (t *TTSClient) playAudio(audioData []byte) error {
	// MP3デコーダーを作成
	decoder, err := mp3.NewDecoder(bytes.NewReader(audioData))
	if err != nil {
		return fmt.Errorf("failed to create MP3 decoder: %w", err)
	}

	// オーディオコンテキストを初期化
	ctx, ready, err := oto.NewContext(decoder.SampleRate(), 2, 2)
	if err != nil {
		return fmt.Errorf("failed to create audio context: %w", err)
	}
	<-ready

	// オーディオプレイヤーを作成
	player := ctx.NewPlayer(decoder)
	defer player.Close()

	// 再生開始
	player.Play()

	// 再生完了まで待機
	for player.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// SaveToFile 音声データをファイルに保存（オプション機能）
func (t *TTSClient) SaveToFile(text, filename string) error {
	fmt.Printf("🎵 音声ファイルを生成中: %s\n", filename)

	// 音声データを生成
	audioData, err := t.synthesize(text)
	if err != nil {
		return fmt.Errorf("音声生成に失敗しました: %w", err)
	}

	// ファイルに保存
	if err := os.WriteFile(filename, audioData, 0644); err != nil {
		return fmt.Errorf("ファイル保存に失敗しました: %w", err)
	}

	fmt.Printf("✅ 音声ファイルを保存しました: %s\n", filename)
	return nil
}
