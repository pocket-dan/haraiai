//go:generate mockgen -source=$GOFILE -destination=../mock/handler_$GOFILE -package=mock
package handler

import (
	"fmt"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/client"
	"github.com/raahii/haraiai/pkg/config"
	"github.com/raahii/haraiai/pkg/log"
	"github.com/raahii/haraiai/pkg/store"
)

type BotHandler interface {
	HandleWebhook(http.ResponseWriter, *http.Request)
}

type BotHandlerImpl struct {
	config config.BotConfig
	bot    client.BotClient
	store  store.Store
}

func NewBotHandler() (*BotHandlerImpl, error) {
	c, err := config.NewBotConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bot config: %w", err)
	}

	bc, err := client.NewBotClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bot client: %w", err)
	}

	s, err := store.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize store: %w", err)
	}

	handler := &BotHandlerImpl{config: c, bot: bc, store: s}
	err = handler.createRichMenu()
	if err != nil {
		return nil, fmt.Errorf("failed to create rich menu: %w", err)
	}

	return handler, nil
}

func (bh *BotHandlerImpl) HandleWebhook(w http.ResponseWriter, req *http.Request) {
	logger := log.NewLogger(req)

	// Allow POST request only.
	if req.Method != http.MethodPost {
		logger.Warnf("unsuppoed request method: %s", req.Method)
		http.Error(w, "Method Not Allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Check X-Line-Signature header value.
	events, err := bh.bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			logger.Warnf("failed to parse webhook request: %s", err)
			http.Error(w, "Bad Request.", http.StatusBadRequest)
		} else {
			logger.Errorf("failed to parse webhook request: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		var err error

		switch event.Type {
		// Bot follow
		case linebot.EventTypeFollow:
			err = bh.handleBotFollow(event)
		// Group join / leave
		case linebot.EventTypeJoin:
			err = bh.handleBotJoin(event)
		case linebot.EventTypeLeave:
			err = bh.handleBotLeave(event)
		// Message
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				err = bh.handleTextMessage(event, message)
			}
		// Postback Action
		case linebot.EventTypePostback:
			err = bh.handlePostbackData(event)
		default:
			logger.Debugf("%s event type is not supported, skip.", event.Type)
		}

		if err != nil {
			logger.Errorf("failed to handle webhook: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Debugf("Event successfully handled. Type: %s, Event ID: %s, User ID: %s\n",
			event.Type, event.WebhookEventID, event.Source.UserID)
	}
}
