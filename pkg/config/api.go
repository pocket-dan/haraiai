//go:generate mockgen -source=$GOFILE -destination=../mock/config_$GOFILE -package=mock
package config

import (
	"errors"
	"os"
)

type ApiConfig interface {
	GetFrontOrigin() string
}

type ApiConfigImpl struct {
	frontBaseURL string
}

func NewApiConfig() (*BotConfigImpl, error) {
	frontBaseURL := os.Getenv("FE_BASE_URL")
	if frontBaseURL == "" {
		return nil, errors.New("$FE_BASE_URL required.")
	}

	return &BotConfigImpl{
		frontBaseURL: frontBaseURL,
	}, nil
}

func (c *BotConfigImpl) GetFrontOrigin() string {
	return c.frontBaseURL
}
