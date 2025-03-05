package robot

import (
	"RestChatBot/src/config"
	"RestChatBot/src/model"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/tmc/langchaingo/llms"
)

var AllBots = make(map[string]*model.BotSetting)
var JsonBotsSetting = []model.BotSetting{} //用來放順序而已。

func init() {
	init_pross()
}

func init_pross() {
	AllBots = make(map[string]*model.BotSetting)
	JsonBotsSetting = []model.BotSetting{}
	settings := config.GetConfig()

	if settings.BotSettingFileName == "" {
		panic("Config setting incorrect: BotSettingFileName")
	}
	settingfilename := filepath.Join(settings.ConfigRoot, settings.BotSettingFileName)
	// 讀取檔案內容
	bytes, err := os.ReadFile(settingfilename)
	if err != nil {
		panic("Not a bot setting file")
	}

	// 解析 JSON
	err = json.Unmarshal(bytes, &JsonBotsSetting)
	if err != nil {
		panic("Can't unmarshal setting file")
	}

	// 預設好角色設定
	for _, bot := range JsonBotsSetting {
		bot.History = append(bot.History, llms.TextParts(llms.ChatMessageTypeSystem, bot.SystemPrompt))
		AllBots[bot.Name] = &bot
	}
}
