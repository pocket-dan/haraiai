package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestHandleBotFollow_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	replyToken := "replyToken"

	b.
		EXPECT().
		ReplyMessage(replyToken, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 2)
		})

	event := newTestFollowEvent(replyToken)
	err := target.handleBotFollow(event)

	assert.Nil(t, err)
}
