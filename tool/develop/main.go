package main

import (
	"log"
	"net/http"

	"github.com/raahii/haraiai/pkg/handler"
	"github.com/rs/cors"
)

// Run for local development.
func main() {
	bot, err := handler.NewBotHandler()
	if err != nil {
		panic(err)
	}

	api, err := handler.NewApiHandler()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/BotWebhookHandler", bot.HandleWebhook)
	mux.HandleFunc("/NotifyInquiry", api.NotifyInquiry)

	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
