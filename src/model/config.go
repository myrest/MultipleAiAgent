package model

type SystemConfig struct {
	OPENAI_API_KEY     string // OpenAI API key
	OPENAI_BASE_URL    string // OpenAI API base URL
	BotSettingFileName string // 機器人設定檔
	AiSource           string // AI來源(openai, ollama, openaistory, ollamastory)
	ModelName          string // 機器人模型 (llama3.1, Openai-4o-mini)
	Sequential         string // 回答模式(sequential, random)
	ConfigRoot         string // 設定檔根目錄
	EnableVoice        bool   // 是否啟用語音
	UseEdgeTTS         bool   // 是否使用Edge TTS
	MaxResponseLength  int    // AI回答最大長度
	LogFileName        string // 記錄檔名
}

type ContextVar string

const (
	ContextMaxResponseLength ContextVar = "MaxResponseLength"
)

type AnsiColor int

const (
	ColorPrompt  AnsiColor = 36
	ColorError   AnsiColor = 31
	ColorMessage AnsiColor = 32
)
