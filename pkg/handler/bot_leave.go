package handler

import "github.com/line/line-bot-sdk-go/v7/linebot"

const (

// GOOD_BYE_MESSAGE_TEXT string = `このグループに関する割り勘データを削除しました。また何かあったら呼んでね`
)

var (
// GOOD_BYE_MESSAGE linebot.SendingMessage = linebot.NewTextMessage(GOOD_BYE_MESSAGE_TEXT)
)

func (bh *BotHandlerImpl) handleBotLeave(event *linebot.Event) error {
	if event.Source.Type != linebot.EventSourceTypeGroup {
		return nil
	}

	// Delete group data.
	groupID := event.Source.GroupID
	if err := bh.store.DeleteGroup(groupID); err != nil {
		return err
	}

	// // Send good by message.
	// if _, err := h.bot.PushMessage(event.Source.GroupID, goodByeMessage).Do(); err != nil {
	// 	return err
	// }

	return nil
}
