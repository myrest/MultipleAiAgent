package model

type VoicePayload struct {
	Model           string  `json:"model"`
	Input           string  `json:"input"`
	Voice           string  `json:"voice"` //voices = ["nova", "shimmer", "echo", "fable", "alloy", "onyx"]
	Response_format string  `json:"response_format"`
	Speed           float32 `json:"speed"`
}
