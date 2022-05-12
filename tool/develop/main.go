package main

import (
	"log"
	"net/http"

	"github.com/raahii/haraiai/pkg/handler"
)

// Run for local development.
func main() {
	bh, err := handler.NewBotHandler()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/bot/webhook", bh.HandleWebhook)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
