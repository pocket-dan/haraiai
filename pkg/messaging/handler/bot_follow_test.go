package handler

import (
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
)

func TestHandleBotFollow_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, b, _, _, target := initializeMocksAndHandler(ctrl)

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
