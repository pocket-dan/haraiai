//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_handler.go -package=mock
package handler

import "net/http"

type BotHandler interface {
	HandleWebhook(http.ResponseWriter, *http.Request)
}
