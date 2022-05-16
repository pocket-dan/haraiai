package handler

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
)

const (
	NOT_SUPPORTED_MESSAGE string = `すみません、このトークで払い合いをお使いいただけません。
  お手数ですが、グループを作成していただき、再度追加してください 🙇`

	GREETING_MESSAGE string = `招待ありがとう！haraiai が二人の割り勘をサポートするよ🤝

まずは2人のニックネームを教えてね。短いときれいに表示できるよ！

○○だよ

と答えてね。`
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
		[]store.User{},
	)

	err := bh.store.SaveGroup(group)
	if err != nil {
		return err
	}

	// Send greeting message.
	err = bh.bot.ReplyTextMessage(event.ReplyToken, GREETING_MESSAGE)
	if err != nil {
		return err
	}

	return nil
}
