package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetUserInput(prompts ...string) (string, bool) {
	reader := bufio.NewReader(os.Stdin)

	if len(prompts) > 0 {
		fmt.Printf("\x1b[%dm%s\x1b[0m\n", 21, prompts[0])
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
