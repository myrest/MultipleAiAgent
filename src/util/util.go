package util

import (
	"RestChatBot/src/config"
	"RestChatBot/src/model"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/acarl005/stripansi"
)

var outputlog string
var logfilename string

func init() {
	config := config.GetConfig()
	logfilename = config.LogFileName
}

func GetUserInput(prompts ...string) (string, bool) {
	reader := bufio.NewReader(os.Stdin)

	if len(prompts) > 0 {
		fmt.Printf("\x1b[1;%dm%s\x1b[0m\n", model.ColorPrompt, prompts[0])
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("讀取輸入時發生錯誤：", err)
		return "", false
	}

	input = strings.TrimSpace(input)

	if input == "exit" {
		fmt.Println("程式結束")
		return "", true
	}

	return input, false
}

func RemoveThinkTags(text string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	return re.ReplaceAllString(text, "")
}

func PutLog(msg string) {
	PutLog_Console(msg)
	PutLog_File(msg)
}

// 因為是用在Console輸出，所以不需要針對tag: think移除
func PutLog_Console(msg string) {
	fmt.Print(msg)
}

// 因為是用在批次寫入，所以需要針對tag: think及ANSI移除
// Todo: 要移除字串的前後空白
func PutLog_File(msg string) {
	//去除think
	msg = RemoveThinkTags(msg)
	//去除ANSI
	msg = stripansi.Strip(msg)
	//移除前後空白
	msg = strings.TrimSpace(msg)
	outputlog += fmt.Sprintln(msg)
}

func SaveLog() {
	if len(outputlog) > 0 {
		file, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("檔案創建失敗")
			return
		}
		defer file.Close()
		file.WriteString(outputlog)
		outputlog = "" //寫入檔案後清空
	}
}
