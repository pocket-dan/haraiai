package handler

import (
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
)

func TestHandleBotLeave_success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, _, _, s, target := initializeMocksAndHandler(ctrl)

	groupID := "groupID"

	s.
		EXPECT().
		DeleteGroup(groupID).
		Times(1)

	event := newTestLeaveEvent(linebot.EventSourceTypeGroup, groupID)
	err := target.handleBotLeave(event)

	assert.Nil(t, err)
}

func TestHandleBotLeave_unsupportedSourceType(t *testing.T) {
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

				_, _, _, s, target := initializeMocksAndHandler(ctrl)

				s.
					EXPECT().
					DeleteGroup(gomock.Any()).
					Times(0)

				event := newTestLeaveEvent(eventSourceType, "dummy group")
				err := target.handleBotLeave(event)

				assert.Nil(t, err)
			})
	}
}
