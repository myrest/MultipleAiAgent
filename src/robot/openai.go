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

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func OpenAIStart(ctx context.Context) {
	config := config.GetConfig()
	opts := []openai.Option{
		openai.WithModel(config.ModelName),
		openai.WithToken(config.OPENAI_API_KEY),
		openai.WithBaseURL(config.OPENAI_BASE_URL),
	}

	llm, err := openai.New(opts...)
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
			OpenAIRollingChat(ctx, llm, config.EnableVoice, skipround)
		} else {
			OpenAIRandomChat(ctx, llm, config.EnableVoice, skipround)
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
					OpenAIConclusion(ctx, llm)
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

func OpenAIConclusion(ctx context.Context, llm *openai.LLM) {
	maxLength := ctx.Value(model.ContextMaxResponseLength).(int)
	bot := AllBots["總結"]
	if bot != nil {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", 22, "總結：")
		_, err := llm.GenerateContent(ctx, bot.History,
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
	} else {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", 22, "未找到總結角色")
	}
}

func OpenAIRollingChat(ctx context.Context, llm *openai.LLM, enableVoide bool, startnum int) {
	player := mp3player.NewMP3Player()
	//依序取出Bot們
	i := 1
	for _, botsetName := range JsonBotsSetting {
		message := createOpenAIResponse(ctx, llm, botsetName.Name)
		if enableVoide {
			filename := fmt.Sprintf("Voice-%d.mp3", (startnum*10)+i)
			i += 1
			err := voicebuilder.ConvertToMp3(message, botsetName.Voice, filename)
			if err != nil {
				panic("轉MP3錯誤。")
			}
			player.Add(filename, nil, nil)
		}
	}
}

func OpenAIRandomChat(ctx context.Context, llm *openai.LLM, enableVoide bool, startnum int) {
	selector := NewRandomBotNameSelector(JsonBotsSetting, enableVoide)
	player := mp3player.NewMP3Player()
	i := 1
	for {
		botName, hasRecord := selector.GetRandom()
		if botName != "" {
			message := createOpenAIResponse(ctx, llm, botName)
			if enableVoide {
				filename := fmt.Sprintf("Voice-%d.mp3", (startnum*10)+i)
				i += 1
				err := voicebuilder.ConvertToMp3(message, AllBots[botName].Voice, filename)
				if err != nil {
					panic("轉MP3錯誤。")
				}
				player.Add(filename, nil, nil)
			}
		}
		if !hasRecord {
			return
		}
	}
}

func createOpenAIResponse(ctx context.Context, llm *openai.LLM, BotName string) string {
	maxLength := ctx.Value(model.ContextMaxResponseLength).(int)
	botset := AllBots[BotName]
	if botset.Name == "總結" { //排除掉特殊角色
		return ""
	}
	//開始Gen answer
	fmt.Printf("\x1b[%dm%s%s\x1b[0m\n", 22, botset.Name, " 發言：")
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
	fmt.Println("\n--------------------------------------------------") //換人，要換行
	return completion.Choices[0].Content
}
