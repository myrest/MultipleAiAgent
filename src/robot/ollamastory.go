package robot

import (
	"RestChatBot/src/config"
	"RestChatBot/src/mp3player"
	"RestChatBot/src/util"
	"RestChatBot/src/voicebuilder"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func OllamaStoryStart(ctx context.Context) {
	config := config.GetConfig()
	opts := []ollama.Option{
		ollama.WithModel(config.ModelName),
	}

	llm, err := ollama.New(opts...)
	if err != nil {
		panic("New robot faield. EX0001.")
	}

	var subject string //使用者先提一個問題
	subject, shouldExist := util.GetUserInput("您的故事主題是？")
	util.PutLog_File(fmt.Sprintln(subject))
	if shouldExist {
		return
	}

	//先將題目sync給所有的bots
	for _, botset := range AllBots {
		botset.Append("", subject) //因為Name為空，比對不到，所以都會變成HumanMessage
	}

	//開始回答
	OllamaStory(ctx, llm, config.EnableVoice)
}

func OllamaStory(ctx context.Context, llm *ollama.LLM, enableVoide bool) {
	//依序取出Bot們
	player := mp3player.NewMP3Player()
	defer player.Close()
	//maxLength := ctx.Value(model.ContextDefaultMaxResponseLength).(int)
	for i, botsetName := range JsonBotsSetting {
		message := createOllamaResponseForStory(ctx, llm, botsetName.Name, botsetName.BotMaxResponseLength)
		if enableVoide {
			filename := fmt.Sprintf("故事-%d.mp3", i)
			i += 1
			err := voicebuilder.ConvertToMp3(util.RemoveThinkTags(message), botsetName.Voice, filename)
			if err != nil {
				panic(fmt.Sprintln("\nVoice:[]", botsetName.Voice, "] 轉MP3錯誤。\n", err.Error()))
			}
			player.Add(filename, nil, nil)
		}
	}
}

func createOllamaResponseForStory(ctx context.Context, llm *ollama.LLM, BotName string, maxLength int) string {
	botset := AllBots[BotName]
	//開始Gen answer
	completion, err := llm.GenerateContent(ctx, botset.History,
		llms.WithTemperature(0.8),
		llms.WithMaxTokens(maxLength),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			util.PutLog_Console(string(chunk))
			return nil
		}),
	)
	if err != nil {
		panic(err.Error())
	}
	util.PutLog_File(completion.Choices[0].Content)
	util.SaveLog()
	//將botset的回答，sync到所有的bots
	for _, makeupBotAnswer := range AllBots {
		makeupBotAnswer.Append(botset.Name, util.RemoveThinkTags(completion.Choices[0].Content))
	}
	return completion.Choices[0].Content
}
