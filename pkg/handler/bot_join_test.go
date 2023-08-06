package handler

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/mock"
	"github.com/raahii/haraiai/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestHandleBotJoin_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	replyToken := "replyToken"
	groupID := "groupID"

	s.
		EXPECT().
		SaveGroup(gomock.Any()).
		Do(func(group *store.Group) {
			assert.Equal(t, groupID, group.ID)
			assert.Len(t, group.Members, 0)
		}).
		Times(1)

	b.
		EXPECT().
		ReplyMessage(replyToken, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 3)
		})

	event := newTestJoinEvent(replyToken, linebot.EventSourceTypeGroup, groupID)
	err := target.handleBotJoin(event)

	assert.Nil(t, err)
}

func TestHandleBotJoin_unsupportedSourceType(t *testing.T) {
	unsupportedEventSourceTypes := []linebot.EventSourceType{
		linebot.EventSourceTypeRoom,
		linebot.EventSourceTypeUser,
	}

	for _, eventSourceType := range unsupportedEventSourceTypes {
		t.Run(
			fmt.Sprintf("eventSourteType: %s", eventSourceType),
			func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				c := mock.NewMockBotConfig(ctrl)
				b := mock.NewMockBotClient(ctrl)
				s := mock.NewMockStore(ctrl)
				target := BotHandlerImpl{config: c, bot: b, store: s}

				replyToken := "replyToken"
				groupID := "groupID"

				b.
					EXPECT().
					ReplyTextMessage(replyToken, NOT_SUPPORTED_MESSAGE).
					Times(1)

				event := newTestJoinEvent(replyToken, eventSourceType, groupID)
				err := target.handleBotJoin(event)

				assert.Nil(t, err)
			})
	}
}
