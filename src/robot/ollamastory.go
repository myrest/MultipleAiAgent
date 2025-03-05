package robot

import (
	"RestChatBot/src/config"
	"RestChatBot/src/model"
	"RestChatBot/src/mp3player"
	"RestChatBot/src/util"
	"RestChatBot/src/voicebuilder"
	"context"
	"fmt"
	"strconv"

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
	if shouldExist {
		return
	}

	StoryRound, shouldExist := util.GetUserInput("要說幾輪次呢？")
	if shouldExist {
		return
	}
	if StoryRound == "" {
		StoryRound = "1"
	}

	//先將題目sync給所有的bots
	for _, botset := range AllBots {
		botset.Append("", subject) //因為Name為空，比對不到，所以都會變成HumanMessage
	}

	totalRound, err := strconv.Atoi(StoryRound)
	if err != nil {
		totalRound = 1
	}
	CurrentRound := 1
	for {
		//開始回答
		OllamaStory(ctx, llm, config.EnableVoice, CurrentRound)
		if CurrentRound < totalRound {
			init_pross()
			CurrentRound++
			continue
		}
	}
}

func OllamaStory(ctx context.Context, llm *ollama.LLM, enableVoide bool, startnum int) {
	//依序取出Bot們
	player := mp3player.NewMP3Player()
	defer player.Close()
	maxLength := ctx.Value(model.ContextMaxResponseLength).(int)
	for i, botsetName := range JsonBotsSetting {
		message := createOllamaResponseForStory(ctx, llm, botsetName.Name, maxLength)
		if enableVoide {
			filename := fmt.Sprintf("故事-%d.mp3", (startnum*10)+i)
			i += 1
			err := voicebuilder.ConvertToMp3(message, botsetName.Voice, filename)
			if err != nil {
				panic("轉MP3錯誤。")
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
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		panic(err.Error())
	}
	//將botset的回答，sync到所有的bots
	for _, makeupBotAnswer := range AllBots {
		makeupBotAnswer.Append(botset.Name, completion.Choices[0].Content)
	}
	return completion.Choices[0].Content
}
