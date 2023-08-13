package client

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type BotClientImpl struct {
	client *linebot.Client
}

func ProvideBotClient() (BotClient, error) {
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

	return &BotClientImpl{client: b}, nil
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

func (bc *BotClientImpl) CreateRichMenu(richMenu linebot.RichMenu) (string, error) {
	resp, err := bc.client.CreateRichMenu(richMenu).Do()
	if err != nil {
		return "", fmt.Errorf("SDK CreateRichMenu method returns error: %w", err)
	}

	return resp.RichMenuID, nil
}

func (bc *BotClientImpl) UploadRichMenuImage(richMenuID, imgPath string) error {
	_, err := bc.client.UploadRichMenuImage(richMenuID, imgPath).Do()
	return err
}

func (bc *BotClientImpl) SetDefaultRichMenu(richMenuID string) error {
	_, err := bc.client.SetDefaultRichMenu(richMenuID).Do()
	return err
}
