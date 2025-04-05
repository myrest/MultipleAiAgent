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
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var systemConfig = config.GetConfig()
var player = mp3player.NewMP3Player()

func OllamaStart(ctx context.Context) {
	opts := []ollama.Option{
		ollama.WithModel(systemConfig.ModelName),
	}
	llm, err := ollama.New(opts...)
	if err != nil {
		panic("New robot faield. EX0001.")
	}

	var subject string //使用者先提一個問題
	subject, shouldExist := util.GetUserInput("您要討論的主題是？")
	util.PutLog_File(fmt.Sprintln(subject))
	if shouldExist {
		return
	}
	//先將題目sync給所有的bots
	for _, botset := range AllBots {
		botset.Append("", subject) //因為Name為空，比對不到，所以都會變成HumanMessage
	}

	skipround := 0
	for {
		//開始回答
		if strings.ToLower(systemConfig.Sequential) == "sequential" {
			OllamaRollingChat(ctx, llm)
		} else {
			OllamaRandomChat(ctx, llm)
		}
		if skipround > 1 {
			//不人為介入，故跳過。
			skipround--
			continue
		}

		input, shouldExist := util.GetUserInput("以下為支援的命令\n 結束程式執行(Exit/Enter) | 自動討論回合數(n) | 由[總結]機器人進行總結(總結)")
		if shouldExist {
			return
		} else if strings.ToLower(input) == "exit" || len(input) < 1 {
			return
		} else {
			num, err := strconv.Atoi(input)
			if err != nil {
				if input == "" {
					continue //輸入為空，直接跳過
				}
				util.PutLog_File(fmt.Sprintln("輸入：", input))
				if input == "總結" {
					OllamaConclusion(ctx, llm)
					return
				} else {
					//將介入的部份sync給所有的bots
					for _, botset := range AllBots {
						botset.Append("", input) //因為Name為空，比對不到，所以都會變成HumanMessage
					}
				}
			} else {
				if num < 1 {
					num = 1
				}
				skipround = num
				fmt.Printf("\x1b[%dm%s%d%s\x1b[0m\n", model.ColorMessage, "--跳過", num, "輪--")
			}
		}
	}
}

func OllamaConclusion(ctx context.Context, llm *ollama.LLM) {
	//maxLength := ctx.Value(model.ContextDefaultMaxResponseLength).(int)
	bot := AllBots["總結"]
	if bot != nil {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", model.ColorPrompt, "總結：")
		completion, err := llm.GenerateContent(ctx, bot.History,
			llms.WithTemperature(0.8),
			//llms.WithMaxLength(maxLength),
			llms.WithMaxTokens(bot.BotMaxResponseLength),
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
	} else {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", model.ColorError, "未找到總結角色")
	}
}

func OllamaRollingChat(ctx context.Context, llm *ollama.LLM) {
	//依序取出Bot們
	for _, botsetName := range JsonBotsSetting {
		createOllamaResponse(ctx, llm, botsetName.Name)
	}
}

func OllamaRandomChat(ctx context.Context, llm *ollama.LLM) {
	selector := NewRandomBotNameSelector(JsonBotsSetting)

	for {
		botName, hasRecord := selector.GetRandom()
		if botName != "" {
			createOllamaResponse(ctx, llm, botName)
		}
		if !hasRecord {
			return
		}
	}
}

func createOllamaResponse(ctx context.Context, llm *ollama.LLM, BotName string) {
	//maxLength := ctx.Value(model.ContextDefaultMaxResponseLength).(int)
	botset := AllBots[BotName]
	if botset.Name == "總結" { //排除掉特殊角色
		return
	}
	//開始Gen answer
	util.PutLog(fmt.Sprintf("\x1b[%dm%s%s\x1b[0m\n", model.ColorPrompt, botset.Name, " 發言："))
	completion, err := llm.GenerateContent(ctx, botset.History,
		llms.WithTemperature(0.8),
		//llms.WithMaxLength(maxLength),
		llms.WithMaxTokens(botset.BotMaxResponseLength),
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

	if systemConfig.EnableVoice {
		filename := fmt.Sprintf("%d-%s.mp3", time.Now().UnixNano(), botset.Name)
		msg := util.RemoveThinkTags(completion.Choices[0].Content)
		err := voicebuilder.ConvertToMp3(msg, botset.Voice, filename)
		if err != nil {
			panic(fmt.Sprintln("\nVoice:[]", botset.Voice, "] 轉MP3錯誤。\n", err.Error()))
		}
		player.Add(filename, nil, nil)
	}

	//將botset的回答，sync到所有的bots
	for _, makeupBotAnswer := range AllBots {
		makeupBotAnswer.Append(botset.Name, util.RemoveThinkTags(completion.Choices[0].Content))
	}
	util.PutLog(fmt.Sprintln("\n--------------------------------------------------")) //換人，要換行
}
