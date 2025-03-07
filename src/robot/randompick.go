package robot

import (
	"RestChatBot/src/model"
	"time"

	"math/rand"
)

type RandomNameSelector struct {
	Items       []string
	UsedIndexes map[int]struct{}
}

func NewRandomBotNameSelector(objArr []model.BotSetting) *RandomNameSelector {
	var items []string
	for _, bot := range JsonBotsSetting {
		items = append(items, bot.Name)
	}
	return &RandomNameSelector{
		Items:       items,
		UsedIndexes: make(map[int]struct{}),
	}
}

func (rs *RandomNameSelector) GetRandom() (string, bool) {
	if len(rs.UsedIndexes) == len(rs.Items) {
		// 如果所有項目都已使用，重置
		rs.UsedIndexes = make(map[int]struct{})
		return "", false
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		index := rand.Intn(len(rs.Items))
		if _, exists := rs.UsedIndexes[index]; !exists {
			rs.UsedIndexes[index] = struct{}{}
			return rs.Items[index], true
		}
	}
}
