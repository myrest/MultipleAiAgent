package config

import (
	"RestChatBot/src/model"
	"encoding/json"
	"os"
	"path/filepath"
)

var configFolder string
var sysConfg model.SystemConfig

func init() {
	//取得存放設定檔目錄。
	settingfilename := os.Getenv("CONFIGFOLDER")
	if settingfilename == "" {
		panic("No Environment variable: CONFIGFOLDER")
	}
	configFolder = settingfilename
}

func GetConfig() model.SystemConfig {
	if sysConfg.BotSettingFileName != "" {
		return sysConfg
	}
	settingfilename := filepath.Join(configFolder, "systemconfig.json")
	// 讀取檔案內容
	bytes, err := os.ReadFile(settingfilename)
	if err != nil {
		panic("Not a bot setting file")
	}

	// 解析 JSON
	err = json.Unmarshal(bytes, &sysConfg)
	if err != nil {
		panic("Can't unmarshal setting file")
	}
	sysConfg.ConfigRoot = configFolder
	return sysConfg
}
