package config

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	RICH_MENU_IMAGE_PATH = "messaging/static/images/richmenu.png"
	FLEX_TEMPLATE_DIR    = "messaging/static/flexmessages"
)

type BotConfigImpl struct {
	frontBaseURL    string
	packageBasePath string
}

func ProvideBotConfig() (BotConfig, error) {
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

func (c *BotConfigImpl) GetFlexTemplateDir() string {
	return filepath.Join(c.packageBasePath, FLEX_TEMPLATE_DIR)
}
