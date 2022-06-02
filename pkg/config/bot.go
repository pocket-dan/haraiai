//go:generate mockgen -source=$GOFILE -destination=../mock/config_$GOFILE -package=mock
package config

import (
	"errors"
	"os"
)

type BotConfig interface {
	GetHelpPageURL() string
}

type BotConfigImpl struct {
	frontBaseURL string
}

func NewBotConfig() (*BotConfigImpl, error) {
	frontBaseURL := os.Getenv("FE_BASE_URL")
	if frontBaseURL == "" {
		return nil, errors.New("$FE_BASE_URL required.")
	}

	return &BotConfigImpl{
		frontBaseURL: frontBaseURL,
	}, nil
}

func (c *BotConfigImpl) GetHelpPageURL() string {
	return c.frontBaseURL + "/help"
}
