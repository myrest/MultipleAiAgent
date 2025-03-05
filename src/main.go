package main

import (
	"RestChatBot/src/config"
	"RestChatBot/src/model"
	"RestChatBot/src/robot"
	"context"
	"fmt"
	"strings"
)

func main() {
	config := config.GetConfig()
	ctx := context.Background()
	ctx = context.WithValue(ctx, model.ContextMaxResponseLength, config.MaxResponseLength)

	switch strings.ToLower(config.AiSource) {
	case "openai":
		robot.OpenAIStart(ctx)
	case "ollama":
		robot.OllamaStart(ctx)
	case "openaistory":
		robot.OpenAIStoryStart(ctx)
	case "ollamastory":
		robot.OllamaStoryStart(ctx)
	default:
		fmt.Println("Unknown AiSource:", config.AiSource)
	}
}
