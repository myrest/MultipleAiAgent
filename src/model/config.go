package model

type SystemConfig struct {
	OPENAI_API_KEY     string
	OPENAI_BASE_URL    string
	BotSettingFileName string
	AiSource           string
	ModelName          string
	Sequential         string
	ConfigRoot         string
	EnableVoice        bool
	UseEdgeTTS         bool
	MaxResponseLength  int
}

type ContextVar string

const (
	ContextMaxResponseLength ContextVar = "MaxResponseLength"
)
