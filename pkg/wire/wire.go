//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/raahii/haraiai/pkg/client"
	"github.com/raahii/haraiai/pkg/messaging/config"
	"github.com/raahii/haraiai/pkg/messaging/flexmessage"
	"github.com/raahii/haraiai/pkg/messaging/handler"
	"github.com/raahii/haraiai/pkg/store"
)

var mainSet = wire.NewSet(
	config.ProvideBotConfig,
	client.ProvideBotClient,
	store.ProvideStore,
	flexmessage.ProvideFlexMessageBuilder,
	handler.ProvideBotHandler,
)

func BuildBotHandler() (handler.BotHandler, error) {
	wire.Build(mainSet)
	return nil, nil
}
