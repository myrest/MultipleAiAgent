package robot

import (
	"github.com/tmc/langchaingo/llms"
)

type MessageHistory struct {
	history []*llms.MessageContent
}

func (mh *MessageHistory) GetHistory() []llms.MessageContent {
	output := make([]llms.MessageContent, len(mh.history))
	for i, item := range mh.history {
		output[i] = *item
	}
	return output
}

func (mh *MessageHistory) Append(message string) {
	newMessage := llms.TextParts(llms.ChatMessageTypeSystem, message)
	if len(mh.history) >= 50 {
		// 移除第一筆資料（FIFO）
		mh.history = mh.history[1:]
	}
	mh.history = append(mh.history, &newMessage)
}

func (mh *MessageHistory) ReplaceSystemRole(message string) {
	newMessage := llms.TextParts(llms.ChatMessageTypeSystem, message)
	if len(mh.history) > 0 {
		mh.history[0] = &newMessage
	}
}
