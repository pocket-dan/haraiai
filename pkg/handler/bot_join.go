package handler

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
)

const (
	NOT_SUPPORTED_MESSAGE string = `すみません、このトークでは haraiai をお使いいただけません。
  お手数ですが、グループを作成していただき、再度追加してください 🙇`
)

var (
	JOIN_REPLYS []linebot.SendingMessage = []linebot.SendingMessage{
		linebot.NewTextMessage(
			"招待ありがとう！haraiai が二人の折半をサポートするよ🤝\n\n" +
				"まずは2人のニックネームを教えてね。短いときれいに表示できるよ！",
		),
		linebot.NewTextMessage("○○だよ"),
		linebot.NewTextMessage("こんなふうに答えてね！"),
	}
)

func (bh *BotHandlerImpl) handleBotJoin(event *linebot.Event) error {
	if event.Source.Type != linebot.EventSourceTypeGroup {
		// Currently, support group talk only.
		err := bh.bot.ReplyTextMessage(event.ReplyToken, NOT_SUPPORTED_MESSAGE)
		if err != nil {
			return err
		}
		return nil
	}

	// Initialize group data.
	group := store.NewGroup(
		event.Source.GroupID,
		store.GROUP_CREATED,
	)

	err := bh.store.SaveGroup(group)
	if err != nil {
		return err
	}

	// Send greeting message.
	err = bh.bot.ReplyMessage(event.ReplyToken, JOIN_REPLYS...)
	if err != nil {
		return err
	}

	return nil
}
