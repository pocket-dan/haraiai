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

func TestHandleTextMessage_totalUp_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{
		bot:   b,
		store: s,
	}

	replyToken := "replyToken"
	group := &store.Group{
		ID:     "group ID",
		Status: store.STARTED,
		Members: map[string]store.User{
			"uid1": store.User{ID: "uid1", Name: "わり夫", PayAmount: 1000},
			"uid2": store.User{ID: "uid2", Name: "わり子", PayAmount: 5000},
		},
	}

	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	expectedMessage := "支払った総額は...\nわり夫: 1000円\nわり子: 5000円\n\nわり子さんが2000円多く払っているよ！"
	b.
		EXPECT().
		ReplyTextMessage(replyToken, gomock.Any()).
		Times(1).
		Do(func(_, message string) {
			assert.Equal(t, expectedMessage, message)
		})

	event := newTestMessageEvent(replyToken, linebot.EventSourceTypeGroup, group.ID, "uid1")
	message := newTextMessage("集計")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleTextMessage_addNewPayment_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{
		bot:   b,
		store: s,
	}

	replyToken := "replyToken"
	warioID := "uid1"
	warikoID := "uid2"
	group := &store.Group{
		ID: "group ID",
		Members: map[string]store.User{
			warioID:  store.User{ID: warioID, Name: "わり夫", PayAmount: 1000},
			warikoID: store.User{ID: warikoID, Name: "わり子", PayAmount: 5000},
		},
	}

	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	s.
		EXPECT().
		SaveGroup(gomock.Any()).
		Do(func(newGroup *store.Group) {
			assert.Equal(t, group.ID, newGroup.ID)

			assert.Len(t, newGroup.Members, 2)

			expectedWario := store.User{ID: warioID, Name: "わり夫", PayAmount: 2000}
			assert.Equal(t, expectedWario, newGroup.Members[warioID])

			assert.Equal(t, group.Members[warikoID], newGroup.Members[warikoID])
		}).
		Times(1)

	expectedMessage := "👍"
	b.
		EXPECT().
		ReplyTextMessage(replyToken, expectedMessage).
		Times(1)

	event := newTestMessageEvent(replyToken, linebot.EventSourceTypeGroup, group.ID, warioID)
	message := newTextMessage("スタバ\n1000円")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func newTestMessageEvent(
	replyToken string,
	eventSourceType linebot.EventSourceType,
	groupID string,
	senderID string,
) *linebot.Event {
	return &linebot.Event{
		Type:       linebot.EventTypeMessage,
		ReplyToken: replyToken,
		Source: &linebot.EventSource{
			Type:    eventSourceType,
			GroupID: groupID,
			UserID:  senderID,
		},
	}
}

func TestHandleTextMessage_unsupportedSourceType(t *testing.T) {
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

				b := mock.NewMockBotClient(ctrl)
				s := mock.NewMockStore(ctrl)
				target := BotHandlerImpl{
					bot:   b,
					store: s,
				}

				b.
					EXPECT().
					ReplyTextMessage(gomock.Any(), gomock.Any()).
					Times(0)

				event := newTestMessageEvent("replyToken", eventSourceType, "dummy", "dummy")
				message := newTextMessage("おーい")
				err := target.handleTextMessage(event, message)

				assert.Nil(t, err)
			})
	}
}

func newTextMessage(message string) *linebot.TextMessage {
	return &linebot.TextMessage{
		ID:     "text message ID",
		Text:   message,
		Emojis: []*linebot.Emoji{},
	}
}
