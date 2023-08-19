package handler

import (
	"time"

	"go.uber.org/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/raahii/haraiai/pkg/mock"
	"github.com/raahii/haraiai/pkg/store"
)

const (
	REPLY_TOKEN string = "replyToken"

	SENDER_ID string = "uid1"
	GROUP_ID  string = "gid1"

	TARO_ID     string = "taro"
	TARO_NAME   string = "太郎"
	HANAKO_ID   string = "hanako"
	HANAKO_NAME string = "花子"
)

var (
	JST                = time.FixedZone("Asia/Tokyo", 9*60*60)
	TIME_GROUP_CREATED = time.Date(2020, time.January, 1, 1, 0, 0, 0, JST)

	DEFAULT_GROUP = newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(0), newHanakoUser(0)},
	)
)

// helper functions to build test linebot.Event

func initializeMocksAndHandler(ctrl *gomock.Controller) (
	*mock.MockBotConfig,
	*mock.MockBotClient,
	*mock.MockFlexMessageBuilder,
	*mock.MockStore,
	*BotHandlerImpl,
) {
	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	f := mock.NewMockFlexMessageBuilder(ctrl)
	s := mock.NewMockStore(ctrl)

	h := &BotHandlerImpl{config: c, bot: b, fs: f, store: s}

	// mock calls for FlexMessageBuilder by default
	f.
		EXPECT().
		BuildLiquidationModeSelectionMessage(gomock.Any()).
		Return(&linebot.FlexMessage{}, nil).
		MaxTimes(1)

	f.
		EXPECT().
		BuildLiquidationPeriodInputMessage(gomock.Any()).
		Return(&linebot.FlexMessage{}, nil).
		MaxTimes(1)

	f.
		EXPECT().
		BuildLiquidationConfirmationMessage(gomock.Any()).
		Return(&linebot.FlexMessage{}, nil).
		MaxTimes(1)

	return c, b, f, s, h
}

func newTestEvent(eventType linebot.EventType) *linebot.Event {
	return &linebot.Event{
		Type:           eventType,
		WebhookEventID: "event ID",
		Source: &linebot.EventSource{
			Type:    linebot.EventSourceTypeGroup,
			UserID:  "user A",
			GroupID: "group A",
		},
	}
}

func newTestJoinEvent(
	replyToken string,
	eventSourceType linebot.EventSourceType,
	groupID string,
) *linebot.Event {
	return &linebot.Event{
		Type:       linebot.EventTypeJoin,
		ReplyToken: replyToken,
		Source: &linebot.EventSource{
			Type:    eventSourceType,
			GroupID: groupID,
		},
	}
}

func newTestLeaveEvent(
	eventSourceType linebot.EventSourceType,
	groupID string,
) *linebot.Event {
	return &linebot.Event{
		Type: linebot.EventTypeLeave,
		Source: &linebot.EventSource{
			Type:    eventSourceType,
			GroupID: groupID,
		},
	}
}

func newTestFollowEvent(
	replyToken string,
) *linebot.Event {
	return &linebot.Event{
		Type:       linebot.EventTypeFollow,
		ReplyToken: replyToken,
		Source: &linebot.EventSource{
			Type:   linebot.EventSourceTypeUser,
			UserID: "dummy ID",
		},
	}
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

func newTestPostbackEvent(
	sourceType linebot.EventSourceType,
	groupID string,
	data string,
	date string,
) *linebot.Event {
	return &linebot.Event{
		Type:       linebot.EventTypePostback,
		ReplyToken: REPLY_TOKEN,
		Source: &linebot.EventSource{
			Type:    sourceType,
			GroupID: groupID,
			UserID:  SENDER_ID,
		},
		Postback: &linebot.Postback{
			Data: data,
			Params: &linebot.Params{
				Date: date,
			},
		},
	}
}

func newTextMessage(message string) *linebot.TextMessage {
	return &linebot.TextMessage{
		ID:     "text message ID",
		Text:   message,
		Emojis: []*linebot.Emoji{},
	}
}

// helper functions to build test entities

func newTestGroup(ID string, status store.GroupStatus, members []*store.User) *store.Group {
	g := new(store.Group)
	g.ID = ID
	g.Status = status
	g.IsTutorial = false

	g.Members = make(map[string]*store.User, len(members))
	for _, u := range members {
		g.Members[u.ID] = u
	}

	g.CreatedAt = TIME_GROUP_CREATED
	g.UpdatedAt = TIME_GROUP_CREATED

	return g
}

func newTaroUser(payAmount int64) *store.User {
	return &store.User{
		ID:               TARO_ID,
		Name:             TARO_NAME,
		PayAmount:        payAmount,
		InitialPayAmount: 200,
		CreatedAt:        TIME_GROUP_CREATED,
		UpdatedAt:        TIME_GROUP_CREATED,
	}
}

func newHanakoUser(payAmount int64) *store.User {
	return &store.User{
		ID:               HANAKO_ID,
		Name:             HANAKO_NAME,
		PayAmount:        payAmount,
		InitialPayAmount: 0,
		CreatedAt:        TIME_GROUP_CREATED,
		UpdatedAt:        TIME_GROUP_CREATED,
	}
}
