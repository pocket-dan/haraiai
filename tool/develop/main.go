package main

import (
	"log"
	"net/http"

	"github.com/raahii/haraiai/pkg/wire"
	"github.com/rs/cors"
)

// Run for local development.
func main() {
	bot, err := wire.BuildBotHandler()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/BotWebhookHandler", bot.HandleWebhook)

	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
