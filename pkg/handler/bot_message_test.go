package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/mock"
	"github.com/raahii/haraiai/pkg/store"
	"github.com/stretchr/testify/assert"
)

const (
	REPLY_TOKEN string = "replyToken"

	SENDER_ID string = "uid1"
	GROUP_ID  string = "gid1"

	WARIO_ID  string = "wario ID"
	WARIKO_ID string = "wariko ID"
)

var (
	JST                = time.FixedZone("Asia/Tokyo", 9*60*60)
	TIME_GROUP_CREATED = time.Date(2020, time.January, 1, 1, 0, 0, 0, JST)
	// TIME_NOW           = time.Date(2022, time.August, 1, 1, 0, 0, 0, JST)

	DEFAULT_GROUP = newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]store.User{newWarioUser(0), newWarikoUser(0)},
	)
)

func TestHandleTextMessage_addNewMember_firstPerson_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_CREATED,
		[]store.User{},
	)

	// Expect to reply text message.
	expectedMessage := "わり夫さんだね！👍"
	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 1)
			assert.Equal(t, linebot.NewTextMessage(expectedMessage), messages[0])
		})

	// Expect to save new group.
	s.
		EXPECT().
		SaveGroup(gomock.Any()).
		Times(1).
		Do(func(newGroup *store.Group) {
			assert.Equal(t, GROUP_ID, newGroup.ID)
			assert.Equal(t, store.GROUP_CREATED, newGroup.Status)
			assert.Len(t, newGroup.Members, 1)

			newUser, exists := newGroup.Members[SENDER_ID]
			assert.True(t, exists)
			assert.Equal(t, SENDER_ID, newUser.ID)
			assert.Equal(t, "わり夫", newUser.Name)
			assert.Equal(t, int64(0), newUser.PayAmount)
		})

		// Test handler.handleTextMessage call.
	event := newTestMessageEvent(REPLY_TOKEN, linebot.EventSourceTypeGroup, GROUP_ID, SENDER_ID)
	err := target.addNewMember(event, group, "わり夫だよ")

	assert.Nil(t, err)
}

func TestHandleTextMessage_addNewMember_secondPerson_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_CREATED,
		[]store.User{newWarioUser(0)},
	)

	// Expect to reply text message.
	expectedMessage := "わり子さんだね！👍"
	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 2)
			assert.Equal(t, linebot.NewTextMessage(expectedMessage), messages[0])
			assert.Equal(t, READY_TO_START_MESSAGES[0], messages[1])
			assert.Equal(t, TIME_GROUP_CREATED, group.CreatedAt)
		})

	// Expect to save new group.
	s.
		EXPECT().
		SaveGroup(gomock.Any()).
		Times(1).
		Do(func(newGroup *store.Group) {
			assert.Equal(t, GROUP_ID, newGroup.ID)
			assert.Equal(t, store.GROUP_STARTED, newGroup.Status)
			assert.Len(t, newGroup.Members, 2)

			newUser, exists := newGroup.Members[SENDER_ID]
			assert.True(t, exists)
			assert.Equal(t, SENDER_ID, newUser.ID)
			assert.Equal(t, "わり子", newUser.Name)
			assert.Equal(t, int64(0), newUser.PayAmount)
		})

		// Test handler.handleTextMessage call.
	event := newTestMessageEvent(REPLY_TOKEN, linebot.EventSourceTypeGroup, GROUP_ID, SENDER_ID)
	err := target.addNewMember(event, group, "わり子だよ")

	assert.Nil(t, err)
}

func TestHandleTextMessage_totalUp_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]store.User{newWarioUser(1000), newWarikoUser(5000)},
	)

	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	expectedMessage := "支払った総額は...\nわり夫: 1000円\nわり子: 5000円\n\nわり子さんが 2000 円多く払っているよ。次はわり夫さんが払うと距離が縮まるね🤝"
	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Equal(t, linebot.NewTextMessage(expectedMessage), messages[0])
		})

	event := newTestMessageEvent(
		REPLY_TOKEN,
		linebot.EventSourceTypeGroup,
		group.ID,
		SENDER_ID,
	)
	message := newTextMessage("集計")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleTextMessage_addNewPayment_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]store.User{newWarioUser(1000), newWarikoUser(5000)},
	)

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
			assert.Equal(t, TIME_GROUP_CREATED, group.CreatedAt)

			expectedWario := newWarioUser(2000)
			actual := newGroup.Members[WARIO_ID]
			assert.Equal(t, expectedWario.ID, actual.ID)
			assert.Equal(t, expectedWario.Name, actual.Name)
			assert.Equal(t, expectedWario.PayAmount, actual.PayAmount)

			assert.Equal(t, group.Members[WARIKO_ID], newGroup.Members[WARIKO_ID])
		}).
		Times(1)

	expectedMessage := "👍"
	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 1)
			assert.Equal(t, linebot.NewTextMessage(expectedMessage), messages[0])
		})

	event := newTestMessageEvent(
		REPLY_TOKEN,
		linebot.EventSourceTypeGroup,
		group.ID,
		WARIO_ID,
	)
	message := newTextMessage("スタバ\n1000円")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleTextMessage_evenUpConfirmation_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]store.User{newWarioUser(1000), newWarikoUser(5000)},
	)

	// Mock and check GetGroup method call.
	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	// Check reply message.
	expectedTextMessage := linebot.NewTextMessage(
		"わり夫さんはわり子さんに 2000 円渡してね🙏",
	)

	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 2)
			assert.Equal(t, expectedTextMessage, messages[0])

			// Omit flex type message verification
			// assert.Equal(t, expectedConfirmationMessage, messages[1])
		})

	event := newTestMessageEvent(
		REPLY_TOKEN,
		linebot.EventSourceTypeGroup,
		group.ID,
		WARIO_ID,
	)
	message := newTextMessage("精算")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleTextMessage_evenUpConfirmation_noNeed_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]store.User{newWarioUser(1000), newWarikoUser(1000)},
	)

	// Mock and check GetGroup method call.
	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

		// Check reply message.
	expectedMessage := linebot.NewTextMessage("払った額は同じ！精算の必要はないよ")
	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 1)
			assert.Equal(t, expectedMessage, messages[0])
		})

	event := newTestMessageEvent(
		REPLY_TOKEN,
		linebot.EventSourceTypeGroup,
		group.ID,
		WARIO_ID,
	)
	message := newTextMessage("精算")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleTextMessage_evenUpComplete_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]store.User{newWarioUser(1000), newWarikoUser(4000)},
	)

	// Mock and check GetGroup method call.
	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	// Check updated group
	s.
		EXPECT().
		SaveGroup(gomock.Any()).
		Times(1).
		Do(func(newGroup *store.Group) {
			assert.Equal(t, group.ID, newGroup.ID)
			assert.Equal(t, store.GROUP_STARTED, newGroup.Status)
			assert.Len(t, newGroup.Members, 2)

			wario, exists := newGroup.Members[WARIKO_ID]
			assert.True(t, exists)
			assert.Equal(t, int64(4000), wario.PayAmount)
		})

		// Check reply message.
	expectedMessage := linebot.NewTextMessage("👍")
	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 1)
			assert.Equal(t, expectedMessage, messages[0])
		})

	event := newTestMessageEvent(
		REPLY_TOKEN,
		linebot.EventSourceTypeGroup,
		group.ID,
		WARIO_ID,
	)
	message := newTextMessage("精算完了")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleHelpMessage_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewMockBotConfig(ctrl)
	b := mock.NewMockBotClient(ctrl)
	s := mock.NewMockStore(ctrl)
	target := BotHandlerImpl{config: c, bot: b, store: s}

	// Mock config
	helpPageURL := "https://test.com/help"
	c.
		EXPECT().
		GetHelpPageURL().
		Return(helpPageURL).
		Times(1)

	// Mock and check GetGroup method call.
	s.
		EXPECT().
		GetGroup(GROUP_ID).
		Return(DEFAULT_GROUP, nil).
		Times(1)

	// Check reply message.
	expectedMessage := linebot.NewTextMessage(
		"ヘルプページはこちら:\n" + helpPageURL,
	)
	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 1)
			assert.Equal(t, expectedMessage, messages[0])
		})

	event := newTestMessageEvent(
		REPLY_TOKEN,
		linebot.EventSourceTypeGroup,
		GROUP_ID,
		WARIO_ID,
	)
	message := newTextMessage("ヘルプ")
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

				c := mock.NewMockBotConfig(ctrl)
				b := mock.NewMockBotClient(ctrl)
				s := mock.NewMockStore(ctrl)
				target := BotHandlerImpl{config: c, bot: b, store: s}

				b.
					EXPECT().
					ReplyTextMessage(REPLY_TOKEN, gomock.Any()).
					Times(0)

				event := newTestMessageEvent(REPLY_TOKEN, eventSourceType, "dummy", "dummy")
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

func newTestGroup(ID string, status store.GroupStatus, members []store.User) *store.Group {
	g := new(store.Group)
	g.ID = ID
	g.Status = status
	g.IsTutorial = false

	g.Members = make(map[string]store.User, len(members))
	for _, u := range members {
		g.Members[u.ID] = u
	}

	g.CreatedAt = TIME_GROUP_CREATED
	g.UpdatedAt = TIME_GROUP_CREATED

	return g
}

// TODO: It's not easy to distinguish 'Wario' and 'Wariko', shold be renamed.
func newWarioUser(payAmount int64) store.User {
	return store.User{
		ID:        WARIO_ID,
		Name:      "わり夫",
		PayAmount: payAmount,
		CreatedAt: TIME_GROUP_CREATED,
		UpdatedAt: TIME_GROUP_CREATED,
	}
}

func newWarikoUser(payAmount int64) store.User {
	return store.User{
		ID:        WARIKO_ID,
		Name:      "わり子",
		PayAmount: payAmount,
		CreatedAt: TIME_GROUP_CREATED,
		UpdatedAt: TIME_GROUP_CREATED,
	}
}
