//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_config.go -package=mock
package config

type BotConfig interface {
	// FE URL
	GetAboutPageURL() string
	GetHelpPageURL() string
	GetInquiryPageURL() string

	// richmenu image path
	GetRichMenuImagePath() string

	// flex message template path
	GetFlexTemplateDir() string
}
