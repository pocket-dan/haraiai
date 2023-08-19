package handler

import (
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
	"github.com/raahii/haraiai/pkg/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestHandlePostbackData_saveLiquidationDate(t *testing.T) {
	group := newTestGroup(
		GROUP_ID,
		store.GROUP_STARTED,
		[]*store.User{newTaroUser(1000), newHanakoUser(1000)},
	)

	cases := []struct {
		name     string
		data     string
		date     string
		expected *store.DateRange
	}{
		{
			name: "update start date",
			data: "liquidationStartDate",
			date: "2023-03-24",
			expected: &store.DateRange{
				Start: time.Date(2023, time.March, 24, 0, 0, 0, 0, timeutil.JST),
			},
		},
		{
			name: "update end date",
			data: "liquidationEndDate",
			date: "2023-04-01",
			expected: &store.DateRange{
				End: time.Date(2023, time.April, 2, 0, 0, 0, 0, timeutil.JST), // set 00:00 on next day
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
				GetGroup(group.ID).
				Return(group, nil).
				Times(1)

			liq := store.Liquidation{}
			s.
				EXPECT().
				GetLiquidation(group.ID).
				Return(&liq, nil).
				Times(1)

			// Check reply message.
			b.
				EXPECT().
				ReplyMessage(REPLY_TOKEN, gomock.Any()).
				Times(1)

			// Check liquidation update
			s.
				EXPECT().
				UpdateLiquidation(group.ID, gomock.Any()).
				Times(1).
				Do(func(_ string, actual *store.Liquidation) {
					assert.Equal(t, tt.expected, actual.Period)
				})

			// Call target method
			event := newTestPostbackEvent(
				linebot.EventSourceTypeGroup,
				group.ID,
				tt.data,
				tt.date,
			)

			err := target.handlePostbackData(event)
			assert.Nil(t, err)
		})
	}
}

func TestHandlePostbackData_IgnoreEventsNotInGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, _, _, s, target := initializeMocksAndHandler(ctrl)

	// These events are not from group
	eventSourceTypes := []linebot.EventSourceType{
		linebot.EventSourceTypeUser,
		linebot.EventSourceTypeRoom,
	}

	for _, sourceType := range eventSourceTypes {
		t.Run("should ignore eventSourceType: "+string(sourceType), func(t *testing.T) {
			event := newTestPostbackEvent(
				sourceType,
				GROUP_ID,
				"liquidationStartDate",
				"2023-03-24",
			)

			err := target.handlePostbackData(event)
			assert.Nil(t, err)

			s.
				EXPECT().
				GetLiquidation(gomock.Any()).
				Times(0)
		})
	}
}
