package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebhook_405_invalidHTTPMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, b, _, _, target := initializeMocksAndHandler(ctrl)

	req, _ := http.NewRequest("GET", "https://example.com", nil)

	// Do not call bot methods
	b.
		EXPECT().
		ParseRequest(gomock.Any()).
		Times(0)

	b.
		EXPECT().
		ReplyMessage(gomock.Any(), gomock.Any()).
		Times(0)

	recorder := httptest.NewRecorder()
	target.HandleWebhook(recorder, req)

	// Return 405 Method Not Allowed.
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func TestHandleWebhook_200_NotTextMessage(t *testing.T) {
	unsupportedEventTypes := []linebot.EventType{
		// linebot.EventTypeMessage,
		// linebot.EventTypeFollow,
		linebot.EventTypeUnfollow,
		// linebot.EventTypeJoin,
		// linebot.EventTypeLeave,
		linebot.EventTypeMemberJoined,
		linebot.EventTypeMemberLeft,
		// linebot.EventTypePostback,
		linebot.EventTypeBeacon,
		linebot.EventTypeAccountLink,
		linebot.EventTypeThings,
		linebot.EventTypeUnsend,
		linebot.EventTypeVideoPlayComplete,
	}

	for _, eventType := range unsupportedEventTypes {
		t.Run(
			fmt.Sprintf("eventType: %s", eventType),
			func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				_, b, _, _, target := initializeMocksAndHandler(ctrl)

				req, _ := http.NewRequest("POST", "https://example.com", nil)

				events := []*linebot.Event{newTestEvent(eventType)}

				b.
					EXPECT().
					ParseRequest(gomock.Any()).
					Return(events, nil).
					Times(1)

				b.
					EXPECT().
					ReplyMessage(gomock.Any(), gomock.Any()).
					Times(0)

				recorder := httptest.NewRecorder()
				target.HandleWebhook(recorder, req)

				assert.Equal(t, http.StatusOK, recorder.Code)
			})
	}
}
