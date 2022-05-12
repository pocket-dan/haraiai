//go:generate mockgen -source=$GOFILE -destination=../mock/client_$GOFILE -package=mock
package client

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// This client is just a wrapper for line-bot-sdk-go client
// to enable mocking/stubbing behavior.

type BotClient interface {
	ParseRequest(*http.Request) ([]*linebot.Event, error)
	ReplyTextMessage(string, ...string) error
	ReplyMessage(string, ...linebot.SendingMessage) error
}

type BotClientImpl struct {
	client *linebot.Client
}

func NewBotClient() (*BotClientImpl, error) {
	channelSecret := os.Getenv("CHANNEL_SECRET")
	if channelSecret == "" {
		return nil, errors.New("$CHANNEL_SECRET required.")
	}

	channelAccessToken := os.Getenv("CHANNEL_ACCESS_TOKEN")
	if channelAccessToken == "" {
		return nil, errors.New("$CHANNEL_ACCESS_TOKEN required.")
	}

	// Initialize LINE Bot client.
	b, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LINE Bot: %w", err)
	}

	return &BotClientImpl{
		client: b,
	}, nil
}

func (bc *BotClientImpl) ParseRequest(r *http.Request) ([]*linebot.Event, error) {
	return bc.client.ParseRequest(r)
}

func (bc *BotClientImpl) ReplyTextMessage(replyToken string, textMessages ...string) error {
	messages := make([]linebot.SendingMessage, 0, len(textMessages))
	for _, m := range textMessages {
		messages = append(messages, linebot.NewTextMessage(m))
	}

	_, err := bc.client.ReplyMessage(replyToken, messages...).Do()
	return err
}

func (bc *BotClientImpl) ReplyMessage(replyToken string, messages ...linebot.SendingMessage) error {
	_, err := bc.client.ReplyMessage(replyToken, messages...).Do()
	return err
}
