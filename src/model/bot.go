package model

import (
	"github.com/tmc/langchaingo/llms"
)

type BotSetting struct {
	Name                 string
	SystemPrompt         string
	Voice                string
	BotMaxResponseLength int
	History              []llms.MessageContent
}

func (mh *BotSetting) Append(name, message string) {
	var chatMsg llms.MessageContent
	if name == mh.Name {
		chatMsg = llms.TextParts(llms.ChatMessageTypeAI, message)
	} else {
		chatMsg = llms.TextParts(llms.ChatMessageTypeHuman, message) //讓AI知道其它的的觀點。
	}
	if len(mh.History) >= 50 {
		mh.History = mh.History[1:]
	}
	mh.History = append(mh.History, chatMsg)
}
