package handler

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	FOLLOW_REPLYS []linebot.SendingMessage = []linebot.SendingMessage{
		linebot.NewTextMessage("友だち追加ありがとうございます。haraiai は二人の折半をサポートするアプリです🤝"),
		linebot.NewTextMessage("はじめるには、一緒に使う相手と haraiai の3人のグループを作ってね！"),
	}
)

func (bh *BotHandlerImpl) handleBotFollow(event *linebot.Event) error {
	// Send explanation message.
	err := bh.bot.ReplyMessage(event.ReplyToken, FOLLOW_REPLYS...)
	if err != nil {
		return err
	}

	return nil
}
