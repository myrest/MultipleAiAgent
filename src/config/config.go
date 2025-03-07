package config

import (
	"RestChatBot/src/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var configFolder string
var sysConfg model.SystemConfig
var defaultConfigFolder = "/Users/roy_tai/LangChain/src/"

func init() {
	//取得存放設定檔目錄。
	settingfilename, ok := os.LookupEnv("CONFIGFOLDER")
	if !ok {
		os.Setenv("CONFIGFOLDER", defaultConfigFolder)
		configFolder = defaultConfigFolder
	} else {
		configFolder = settingfilename
	}
}

func GetConfig() model.SystemConfig {
	if sysConfg.BotSettingFileName != "" {
		return sysConfg
	}
	settingfilename := filepath.Join(configFolder, "systemconfig.json")
	// 讀取檔案內容
	_, err := os.Stat(settingfilename)
	if err != nil {
		if os.IsNotExist(err) && defaultConfigFolder == configFolder {
			panic("Environment variable not found: CONFIGFOLDER")
		} else {
			panic("Not a bot setting file")
		}
	}

	bytes, err := os.ReadFile(settingfilename)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	err = json.Unmarshal(bytes, &sysConfg)
	if err != nil {
		panic("Can't unmarshal setting file")
	}
	sysConfg.ConfigRoot = configFolder
	filename := fmt.Sprintf("%d-%d-%d-%d-%d.log", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute())
	sysConfg.LogFileName = filename
	return sysConfg
}
