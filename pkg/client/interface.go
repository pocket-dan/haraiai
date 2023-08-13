//go:generate mockgen -source=$GOFILE -destination=../mock/mock_client.go -package=mock
package client

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// This package is just a wrapper for line-bot-sdk-go client
// to enable mocking/stubbing behavior.

type BotClient interface {
	ParseRequest(*http.Request) ([]*linebot.Event, error)
	ReplyTextMessage(string, ...string) error
	ReplyMessage(string, ...linebot.SendingMessage) error
	CreateRichMenu(linebot.RichMenu) (string, error)
	UploadRichMenuImage(string, string) error
	SetDefaultRichMenu(string) error
}
