//go:generate mockgen -source=$GOFILE -destination=../mock/config_$GOFILE -package=mock
package config

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	RICH_MENU_IMAGE_PATH = "images/richmenu.png"
)

type BotConfig interface {
	// FE URL
	GetAboutPageURL() string
	GetHelpPageURL() string
	GetInquiryPageURL() string

	// richmenu image path
	GetRichMenuImagePath() string
}

type BotConfigImpl struct {
	frontBaseURL    string
	packageBasePath string
}

func NewBotConfig() (*BotConfigImpl, error) {
	frontBaseURL := os.Getenv("FE_BASE_URL")
	if frontBaseURL == "" {
		return nil, errors.New("$FE_BASE_URL required.")
	}

	packageBasePath := os.Getenv("PACKAGE_BASE_PATH")
	if packageBasePath == "" {
		return nil, errors.New("$PACKAGE_BASE_PATH required")
	}

	return &BotConfigImpl{
		frontBaseURL:    frontBaseURL,
		packageBasePath: packageBasePath,
	}, nil
}

func (c *BotConfigImpl) GetAboutPageURL() string {
	return c.frontBaseURL
}

func (c *BotConfigImpl) GetHelpPageURL() string {
	return c.frontBaseURL + "/help"
}

func (c *BotConfigImpl) GetInquiryPageURL() string {
	return c.frontBaseURL + "/inquiry"
}

func (c *BotConfigImpl) GetRichMenuImagePath() string {
	return filepath.Join(c.packageBasePath, RICH_MENU_IMAGE_PATH)
}
