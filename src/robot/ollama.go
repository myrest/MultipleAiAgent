package robot

import (
	"RestChatBot/src/config"
	"RestChatBot/src/model"
	"RestChatBot/src/util"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func OllamaStart(ctx context.Context) {
	config := config.GetConfig()
	opts := []ollama.Option{
		ollama.WithModel(config.ModelName),
	}
	llm, err := ollama.New(opts...)
	if err != nil {
		panic("New robot faield. EX0001.")
	}

	var subject string //使用者先提一個問題
	subject, shouldExist := util.GetUserInput("您要討論的主題是？")
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
		if strings.ToLower(config.Sequential) == "sequential" {
			OllamaRollingChat(ctx, llm)
		} else {
			OllamaRandomChat(ctx, llm)
		}
		if skipround > 1 {
			//不人為介入，故跳過。
			skipround--
			continue
		}

		input, shouldExist := util.GetUserInput("以下為支援的命令\n 不介入，由AI繼續討論(Enter) | 結束程式執行(Exit) | 自動討論回合數(n) | 由[總結]機器人進行總結(總結)")
		if shouldExist {
			return
		} else if strings.ToLower(input) == "exit" {
			return
		} else {
			num, err := strconv.Atoi(input)
			if err != nil {
				if input == "" {
					continue //輸入為空，直接跳過
				}
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
				fmt.Printf("\x1b[%dm%s%d%s\x1b[0m\n", 21, "--跳過", num, "輪--")
			}
		}
	}
}

func OllamaConclusion(ctx context.Context, llm *ollama.LLM) {
	maxLength := ctx.Value(model.ContextMaxResponseLength).(int)
	bot := AllBots["總結"]
	if bot != nil {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", 22, "總結：")
		_, err := llm.GenerateContent(ctx, bot.History,
			llms.WithTemperature(0.8),
			//llms.WithMaxLength(maxLength),
			llms.WithMaxTokens(maxLength),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fmt.Print(string(chunk))
				return nil
			}),
		)
		if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", 22, "未找到總結角色")
	}
}

func OllamaRollingChat(ctx context.Context, llm *ollama.LLM) {
	//依序取出Bot們
	for _, botsetName := range JsonBotsSetting {
		createOllamaResponse(ctx, llm, botsetName.Name)
	}
}

func OllamaRandomChat(ctx context.Context, llm *ollama.LLM) {
	selector := NewRandomBotNameSelector(JsonBotsSetting, false)

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
	maxLength := ctx.Value(model.ContextMaxResponseLength).(int)
	botset := AllBots[BotName]
	if botset.Name == "總結" { //排除掉特殊角色
		return
	}
	//開始Gen answer
	fmt.Printf("\x1b[%dm%s%s\x1b[0m\n", 22, botset.Name, " 發言：")
	completion, err := llm.GenerateContent(ctx, botset.History,
		llms.WithTemperature(0.8),
		//llms.WithMaxLength(maxLength),
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
	fmt.Println("\n--------------------------------------------------") //換人，要換行
}
