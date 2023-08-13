package handler

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
	"github.com/raahii/haraiai/pkg/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestHandleTextMessage_addNewMember_firstPerson_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, b, _, s, target := initializeMocksAndHandler(ctrl)

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_CREATED,
		[]*store.User{},
	)

	// Expect to reply text message.
	expectedMessage := "太郎さんだね！👍"
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
			assert.Equal(t, "太郎", newUser.Name)
			assert.Equal(t, int64(0), newUser.PayAmount)
			assert.Equal(t, int64(0), newUser.InitialPayAmount)
		})

		// Test handler.handleTextMessage call.
	event := newTestMessageEvent(REPLY_TOKEN, linebot.EventSourceTypeGroup, GROUP_ID, SENDER_ID)
	err := target.addNewMember(event, group, "太郎だよ")

	assert.Nil(t, err)
}

func TestHandleTextMessage_addNewMember_secondPerson_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, b, _, s, target := initializeMocksAndHandler(ctrl)

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_CREATED,
		[]*store.User{newTaroUser(0)},
	)

	// Expect to reply text message.
	expectedMessage := "花子さんだね！👍"
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
			assert.Equal(t, "花子", newUser.Name)
			assert.Equal(t, int64(0), newUser.PayAmount)
			assert.Equal(t, int64(0), newUser.InitialPayAmount)
		})

		// Test handler.handleTextMessage call.
	event := newTestMessageEvent(REPLY_TOKEN, linebot.EventSourceTypeGroup, GROUP_ID, SENDER_ID)
	err := target.addNewMember(event, group, "花子だよ")

	assert.Nil(t, err)
}

func TestHandleTextMessage_totalUp_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, b, _, s, target := initializeMocksAndHandler(ctrl)

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(1000), newHanakoUser(5000)},
	)

	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	expectedMessage := "支払った総額は...\n太郎: 1000円\n花子: 5000円\n\n花子さんが 2000 円多く払っているよ。次は太郎さんが払うと距離が縮まるね🤝"
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

	_, b, _, s, target := initializeMocksAndHandler(ctrl)

	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(1000), newHanakoUser(5000)},
	)

	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	s.
		EXPECT().
		CreatePayment(group.ID, gomock.Any()).
		Do(func(_ string, payment *store.Payment) {
			assert.Equal(t, "スタバ", payment.Name)
			assert.Equal(t, int64(1000), payment.Amount)
			assert.Equal(t, store.PAYMENT_TYPE_DEFAULT, payment.Type)
			assert.Equal(t, TARO_ID, payment.PayerID)
		}).
		Times(1)

	s.
		EXPECT().
		SaveGroup(gomock.Any()).
		Do(func(newGroup *store.Group) {
			assert.Equal(t, group.ID, newGroup.ID)
			assert.Len(t, newGroup.Members, 2)
			assert.Equal(t, TIME_GROUP_CREATED, group.CreatedAt)

			expectedTaro := newTaroUser(2000)
			actual := newGroup.Members[TARO_ID]
			assert.Equal(t, expectedTaro.ID, actual.ID)
			assert.Equal(t, expectedTaro.Name, actual.Name)
			assert.Equal(t, expectedTaro.PayAmount, actual.PayAmount)
			assert.Equal(t, expectedTaro.InitialPayAmount, actual.InitialPayAmount)

			assert.Equal(t, group.Members[HANAKO_ID], newGroup.Members[HANAKO_ID])
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
		TARO_ID,
	)
	message := newTextMessage("スタバ\n1000円")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleTextMessage_startLiquidation_success(t *testing.T) {
	inputTexts := []string{
		"清算",
		"清算したい",
		"精算",
		"精算したい",
		"リセット",
		"リセットしたい",
	}

	for _, text := range inputTexts {
		caseName := fmt.Sprintf("input text: %s", text)

		t.Run(caseName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, b, _, s, target := initializeMocksAndHandler(ctrl)

			// Mock and check GetGroup method call.
			group := newTestGroup(
				GROUP_ID,
				store.GROUP_STARTED,
				[]*store.User{newTaroUser(1000), newHanakoUser(5000)},
			)
			s.
				EXPECT().
				GetGroup(group.ID).
				Return(group, nil).
				Times(1)

			// Mock and check CreateLiquidation method call.
			s.
				EXPECT().
				CreateLiquidation(group.ID, store.Liquidation{}).
				Return(nil).
				Times(1)

			b.
				EXPECT().
				ReplyMessage(REPLY_TOKEN, gomock.Any()).
				Times(1).
				Do(func(_ string, messages ...linebot.SendingMessage) {
					assert.Len(t, messages, 1)

					// reply message is flex message, so ignore content
				})

			event := newTestMessageEvent(
				REPLY_TOKEN,
				linebot.EventSourceTypeGroup,
				group.ID,
				TARO_ID,
			)
			message := newTextMessage(text)
			err := target.handleTextMessage(event, message)

			assert.Nil(t, err)
		})
	}
}

func TestHandleTextMessage_calculateLiquidationAmount_whole(t *testing.T) {
	cases := []struct {
		name           string
		group          *store.Group
		whoPayLessName string
		whoPayALotName string
		expected       store.Liquidation
	}{
		{
			name: "taro's payAmount is greater than hanako's",
			group: newTestGroup(
				GROUP_ID,
				store.GROUP_STARTED,
				[]*store.User{newTaroUser(500), newHanakoUser(0)},
			),
			whoPayLessName: HANAKO_NAME,
			whoPayALotName: TARO_NAME,
			expected: store.Liquidation{
				PayerID: HANAKO_ID,
				Amount:  250,
			},
		},
		{
			name: "hanako's payAmount is greater than taro's",
			group: newTestGroup(
				GROUP_ID,
				store.GROUP_STARTED,
				[]*store.User{newTaroUser(1000), newHanakoUser(4000)},
			),
			whoPayLessName: TARO_NAME,
			whoPayALotName: HANAKO_NAME,
			expected: store.Liquidation{
				PayerID: TARO_ID,
				Amount:  1500,
			},
		},
		{
			name: "pay amount is negative value",
			group: newTestGroup(
				GROUP_ID,
				store.GROUP_STARTED,
				[]*store.User{newTaroUser(-100), newHanakoUser(100)},
			),
			whoPayLessName: TARO_NAME,
			whoPayALotName: HANAKO_NAME,
			expected: store.Liquidation{
				PayerID: TARO_ID,
				Amount:  100,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, b, _, s, target := initializeMocksAndHandler(ctrl)

			// Mock and check GetGroup method call.
			s.
				EXPECT().
				GetGroup(tt.group.ID).
				Return(tt.group, nil).
				Times(1)

			// Mock and check GetLiquidation method call.
			s.
				EXPECT().
				GetLiquidation(tt.group.ID).
				Return(&store.Liquidation{}, nil).
				Times(1)

			// Check liquidation update
			s.
				EXPECT().
				UpdateLiquidation(tt.group.ID, gomock.Any()).
				Do(func(_ string, liq *store.Liquidation) {
					assert.Nil(t, liq.Period)
					assert.Equal(t, tt.expected.Amount, liq.Amount)
					assert.Equal(t, tt.expected.PayerID, liq.PayerID)
				}).
				Times(1)

			// Check reply message.
			expectedMessage := linebot.NewTextMessage(
				fmt.Sprintf("%sさんは%sさんに %d 円渡してね🙏", tt.whoPayLessName, tt.whoPayALotName, tt.expected.Amount),
			)
			b.
				EXPECT().
				ReplyMessage(REPLY_TOKEN, gomock.Any()).
				Times(1).
				Do(func(_ string, messages ...linebot.SendingMessage) {
					assert.Len(t, messages, 2)
					assert.Equal(t, expectedMessage, messages[0])
				})

			event := newTestMessageEvent(
				REPLY_TOKEN,
				linebot.EventSourceTypeGroup,
				tt.group.ID,
				HANAKO_ID,
			)
			message := newTextMessage("清算額を計算")
			err := target.handleTextMessage(event, message)

			assert.Nil(t, err)
		})
	}
}

func TestHandleTextMessage_calculateLiquidationAmount_partial(t *testing.T) {
	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(500), newHanakoUser(0)},
	)

	liq := store.Liquidation{
		Period: &store.DateRange{
			Start: time.Date(2023, time.May, 3, 0, 0, 0, 0, timeutil.JST),
			End:   time.Date(2023, time.May, 8, 0, 0, 0, 0, timeutil.JST),
		},
	}

	cases := []struct {
		name           string
		payAmountMap   map[string]int64
		whoPayLessName string
		whoPayALotName string
		expected       store.Liquidation
	}{
		{
			name:           "taro's payAmount is greater than hanako's",
			payAmountMap:   map[string]int64{TARO_ID: 500, HANAKO_ID: 0},
			whoPayLessName: HANAKO_NAME,
			whoPayALotName: TARO_NAME,
			expected: store.Liquidation{
				PayerID: HANAKO_ID,
				Amount:  250,
			},
		},
		{
			name:           "hanako's payAmount is greater than taro's",
			payAmountMap:   map[string]int64{TARO_ID: 1000, HANAKO_ID: 4000},
			whoPayLessName: TARO_NAME,
			whoPayALotName: HANAKO_NAME,
			expected: store.Liquidation{
				PayerID: TARO_ID,
				Amount:  1500,
			},
		},
		{
			name:           "pay amount is negative value",
			payAmountMap:   map[string]int64{TARO_ID: -100, HANAKO_ID: 100},
			whoPayLessName: TARO_NAME,
			whoPayALotName: HANAKO_NAME,
			expected: store.Liquidation{
				PayerID: TARO_ID,
				Amount:  100,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, b, _, s, target := initializeMocksAndHandler(ctrl)

			// Mock and check GetGroup method call.
			s.
				EXPECT().
				GetGroup(group.ID).
				Return(group, nil).
				Times(1)

			// Mock and check GetLiquidation method call.
			s.
				EXPECT().
				GetLiquidation(group.ID).
				Return(&liq, nil).
				Times(1)

			// Mock and check BuildPayAmountMapBetweenCreatedAt method call.
			s.
				EXPECT().
				BuildPayAmountMapBetweenCreatedAt(group.ID, liq.Period).
				Return(tt.payAmountMap, nil).
				Times(1)

			// Check liquidation update
			s.
				EXPECT().
				UpdateLiquidation(group.ID, gomock.Any()).
				Do(func(_ string, actual *store.Liquidation) {
					assert.Equal(t, liq.Period, actual.Period)
					assert.Equal(t, tt.expected.Amount, actual.Amount)
					assert.Equal(t, tt.expected.PayerID, actual.PayerID)
				}).
				Times(1)

			// Check reply message.
			expectedMessage := linebot.NewTextMessage(
				fmt.Sprintf("%sさんは%sさんに %d 円渡してね🙏", tt.whoPayLessName, tt.whoPayALotName, tt.expected.Amount),
			)
			b.
				EXPECT().
				ReplyMessage(REPLY_TOKEN, gomock.Any()).
				Times(1).
				Do(func(_ string, messages ...linebot.SendingMessage) {
					assert.Len(t, messages, 2)
					assert.Equal(t, expectedMessage, messages[0])
				})

			event := newTestMessageEvent(
				REPLY_TOKEN,
				linebot.EventSourceTypeGroup,
				group.ID,
				HANAKO_ID,
			)
			message := newTextMessage("清算額を計算")
			err := target.handleTextMessage(event, message)

			assert.Nil(t, err)
		})
	}
}

func TestHandleTextMessage_startPartialLiquidationSetting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, b, _, s, target := initializeMocksAndHandler(ctrl)

	// Mock and check GetGroup method call.
	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(1000), newHanakoUser(5000)},
	)
	s.
		EXPECT().
		GetGroup(group.ID).
		Return(group, nil).
		Times(1)

	// Mock and check CreateLiquidation method call.
	s.
		EXPECT().
		CreateLiquidation(group.ID, store.Liquidation{}).
		Return(nil).
		Times(1)

	b.
		EXPECT().
		ReplyMessage(REPLY_TOKEN, gomock.Any()).
		Times(1).
		Do(func(_ string, messages ...linebot.SendingMessage) {
			assert.Len(t, messages, 1)

			// reply message is flex message, so ignore content
		})

		// Call test target
	event := newTestMessageEvent(
		REPLY_TOKEN,
		linebot.EventSourceTypeGroup,
		group.ID,
		TARO_ID,
	)
	message := newTextMessage("特定期間のみ清算する")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleTextMessage_calculateLiquidationAmount_noDifference(t *testing.T) {
	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(1000), newHanakoUser(1000)},
	)

	cases := []struct {
		name         string
		liq          store.Liquidation
		payAmountMap map[string]int64
	}{
		{
			name:         "target period is nothing",
			liq:          store.Liquidation{},
			payAmountMap: map[string]int64{},
		},
		{
			name: "target period exists but incomplete map",
			liq: store.Liquidation{
				Period: &store.DateRange{
					Start: time.Date(2023, time.May, 3, 0, 0, 0, 0, timeutil.JST),
					End:   time.Date(2023, time.May, 8, 0, 0, 0, 0, timeutil.JST),
				},
			},
			payAmountMap: map[string]int64{},
		},
		{
			name: "target period exists and same amount in map",
			liq: store.Liquidation{
				Period: &store.DateRange{
					Start: time.Date(2023, time.May, 3, 0, 0, 0, 0, timeutil.JST),
					End:   time.Date(2023, time.May, 8, 0, 0, 0, 0, timeutil.JST),
				},
			},
			payAmountMap: map[string]int64{
				TARO_ID:   10000,
				HANAKO_ID: 10000,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, b, _, s, target := initializeMocksAndHandler(ctrl)

			// Mock and check GetGroup method call.
			s.
				EXPECT().
				GetGroup(group.ID).
				Return(group, nil).
				Times(1)

			// Mock and check GetLiquidation method call.
			s.
				EXPECT().
				GetLiquidation(group.ID).
				Return(&tt.liq, nil).
				Times(1)

			// Mock and check BuildPayAmountMapBetweenCreatedAt method call.
			s.
				EXPECT().
				BuildPayAmountMapBetweenCreatedAt(group.ID, tt.liq.Period).
				Return(tt.payAmountMap, nil).
				MaxTimes(1)

			// Check reply message.
			expectedMessage := linebot.NewTextMessage("払った額は同じ！清算の必要はないよ")
			b.
				EXPECT().
				ReplyMessage(REPLY_TOKEN, gomock.Any()).
				Times(1).
				Do(func(_ string, messages ...linebot.SendingMessage) {
					assert.Len(t, messages, 1)
					assert.Equal(t, expectedMessage, messages[0])
				})

			// Check liquidation deleted
			s.
				EXPECT().
				DeleteLiquidation(group.ID).
				Return(nil).
				Times(1)

			event := newTestMessageEvent(
				REPLY_TOKEN,
				linebot.EventSourceTypeGroup,
				group.ID,
				TARO_ID,
			)
			message := newTextMessage("清算額を計算")
			err := target.handleTextMessage(event, message)

			assert.Nil(t, err)
		})
	}
}

func TestHandleTextMessage_completeLiquidation(t *testing.T) {
	cases := []struct {
		name                 string
		group                *store.Group
		liquidation          *store.Liquidation
		expectedPayAmountMap map[string]int64
	}{
		{
			name: "Add 1500 yen to taro",
			group: newTestGroup(
				GROUP_ID,
				store.GROUP_STARTED,
				[]*store.User{newTaroUser(1000), newHanakoUser(4000)},
			),
			liquidation: &store.Liquidation{
				PayerID: TARO_ID,
				Amount:  1500,
			},
			expectedPayAmountMap: map[string]int64{
				TARO_ID:   4000,
				HANAKO_ID: 4000,
			},
		},
		{
			name: "Add 1000 yen to taro regardless current pay amount",
			group: newTestGroup(
				GROUP_ID,
				store.GROUP_STARTED,
				[]*store.User{newTaroUser(1000), newHanakoUser(1000)},
			),
			liquidation: &store.Liquidation{
				PayerID: TARO_ID,
				Amount:  1000,
			},
			expectedPayAmountMap: map[string]int64{
				TARO_ID:   3000,
				HANAKO_ID: 1000,
			},
		},
		{
			name: "Do liduidation regardless date range",
			group: newTestGroup(
				GROUP_ID,
				store.GROUP_STARTED,
				[]*store.User{newTaroUser(1000), newHanakoUser(1000)},
			),
			liquidation: &store.Liquidation{
				Period:  &store.DateRange{},
				PayerID: TARO_ID,
				Amount:  100,
			},
			expectedPayAmountMap: map[string]int64{
				TARO_ID:   1200,
				HANAKO_ID: 1000,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, b, _, s, target := initializeMocksAndHandler(ctrl)

			s.
				EXPECT().
				GetGroup(tt.group.ID).
				Return(tt.group, nil).
				Times(1)

			s.
				EXPECT().
				GetLiquidation(tt.group.ID).
				Return(tt.liquidation, nil).
				Times(1)

			s.
				EXPECT().
				CreatePayment(tt.group.ID, gomock.Any()).
				Do(func(_ string, payment *store.Payment) {
					assert.Equal(t, "清算", payment.Name)
					assert.Equal(t, store.PAYMENT_TYPE_LIQUIDATION, payment.Type)
					assert.Equal(t, tt.liquidation.Amount, payment.Amount)
					assert.Equal(t, tt.liquidation.PayerID, payment.PayerID)
				}).
				Times(1)

			// Check updated group
			s.
				EXPECT().
				SaveGroup(gomock.Any()).
				Times(1).
				Do(func(newGroup *store.Group) {
					assert.Equal(t, tt.group.ID, newGroup.ID)
					assert.Equal(t, store.GROUP_STARTED, newGroup.Status)
					assert.Len(t, newGroup.Members, 2)

					hanako, exists := newGroup.Members[HANAKO_ID]
					assert.True(t, exists)
					assert.Equal(t, tt.expectedPayAmountMap[HANAKO_ID], hanako.PayAmount)

					taro, exists := newGroup.Members[TARO_ID]
					assert.True(t, exists)
					assert.Equal(t, tt.expectedPayAmountMap[TARO_ID], taro.PayAmount)
				})

			s.
				EXPECT().
				DeleteLiquidation(tt.group.ID).
				Times(1)

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
				tt.group.ID,
				TARO_ID,
			)
			message := newTextMessage("清算完了")
			err := target.handleTextMessage(event, message)

			assert.Nil(t, err)
		})
	}
}

func TestHandleTextMessage_completeLiquidation_invalidLiquidation(t *testing.T) {
	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(1000), newHanakoUser(4000)},
	)

	cases := []struct {
		name                 string
		liquidation          *store.Liquidation
		getLiquidationResult error
	}{
		{
			name:                 "should ignore if liquidation doesn't exist",
			liquidation:          nil,
			getLiquidationResult: errors.New("test"),
		},
		{
			name: "should ignore if liquidation amount is zero",
			liquidation: &store.Liquidation{
				PayerID: TARO_ID,
				Amount:  0,
			},
			getLiquidationResult: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, _, _, s, target := initializeMocksAndHandler(ctrl)

			s.
				EXPECT().
				GetGroup(group.ID).
				Return(group, nil).
				Times(1)

			s.
				EXPECT().
				GetLiquidation(group.ID).
				Return(tt.liquidation, tt.getLiquidationResult).
				Times(1)

			s.
				EXPECT().
				DeleteLiquidation(group.ID).
				Times(1)

			event := newTestMessageEvent(
				REPLY_TOKEN,
				linebot.EventSourceTypeGroup,
				group.ID,
				TARO_ID,
			)
			message := newTextMessage("清算完了")
			err := target.handleTextMessage(event, message)

			assert.Nil(t, err)
		})
	}
}

func TestHandleHelpMessage_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c, b, _, s, target := initializeMocksAndHandler(ctrl)

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
		TARO_ID,
	)
	message := newTextMessage("ヘルプ")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleMessageForNameChangeGuide(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, b, _, s, target := initializeMocksAndHandler(ctrl)

	// Mock and check GetGroup method call.
	s.
		EXPECT().
		GetGroup(GROUP_ID).
		Return(DEFAULT_GROUP, nil).
		Times(1)

	// Check reply message.
	expectedMessage := linebot.NewTextMessage(
		"名前を変更したいときは\n「名前を○○に変更」\nと言ってね！",
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
		TARO_ID,
	)
	message := newTextMessage("名前を変更")
	err := target.handleTextMessage(event, message)

	assert.Nil(t, err)
}

func TestHandleNameChange(t *testing.T) {
	cases := []struct {
		message string
		newName string
	}{
		{"名前をほげに変更", "ほげ"},
		{"名前をテスト 太郎に変更", "テスト 太郎"},
		{"名前を    テスト太郎\nに変更", "テスト太郎"},
	}

	for _, tt := range cases {
		t.Run(tt.message, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, b, _, s, target := initializeMocksAndHandler(ctrl)

			// Mock and check GetGroup method call.
			s.
				EXPECT().
				GetGroup(GROUP_ID).
				Return(DEFAULT_GROUP, nil).
				Times(1)

			s.
				EXPECT().
				SaveGroup(gomock.Any()).
				Times(1).
				Do(func(newGroup *store.Group) {
					assert.Equal(t, GROUP_ID, newGroup.ID)
					assert.Len(t, newGroup.Members, 2)

					taro, exists := newGroup.Members[TARO_ID]
					assert.True(t, exists)

					// should be changed
					assert.Equal(t, tt.newName, taro.Name)

					// should not be changed
					assert.Equal(t, int64(0), taro.PayAmount)
					assert.Equal(t, int64(200), taro.InitialPayAmount)

					hanako, exists := newGroup.Members[HANAKO_ID]
					assert.Equal(t, hanako, newHanakoUser(0))
				})

			// Check reply message.
			expectedMessage := linebot.NewTextMessage(
				fmt.Sprintf("名前を「%s」に変更しました👍", tt.newName),
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
				TARO_ID,
			)
			message := newTextMessage(tt.message)
			err := target.handleTextMessage(event, message)

			assert.Nil(t, err)
		})
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

				_, b, _, _, target := initializeMocksAndHandler(ctrl)

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
