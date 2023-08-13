package bot

import (
	"net/http"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/raahii/haraiai/pkg/messaging/handler"
	"github.com/raahii/haraiai/pkg/wire"
)

var botHandler handler.BotHandler

func init() {
	var err error

	botHandler, err = wire.BuildBotHandler()
	if err != nil {
		panic(err)
	}
}

func HandleWebhook(w http.ResponseWriter, req *http.Request) {
	botHandler.HandleWebhook(w, req)
}
