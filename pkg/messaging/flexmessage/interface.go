//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_flexmessage.go -package=mock

package flexmessage

import "github.com/line/line-bot-sdk-go/v7/linebot"

type FlexMessageBuilder interface {
	BuildLiquidationModeSelectionMessage(LiquidationModeParams) (*linebot.FlexMessage, error)
	BuildLiquidationPeriodInputMessage(LiquidationInputPeriodParams) (*linebot.FlexMessage, error)
	BuildLiquidationConfirmationMessage(LiquidationConfirmationParams) (*linebot.FlexMessage, error)
}

// 清算モードを選択するときのメッセージ
type LiquidationModeParams struct {
	SelectWholeModeText   string
	SelectPartialModeText string
}

// 清算が終わったか確認するときのメッセージ
type LiquidationConfirmationParams struct {
	OkMessageText string
}

// 特定期間清算の場合の期間を選択するときのメッセージ
type LiqudationSelectDateParams struct {
	Data        string
	InitialDate string
	MaxDate     string
	MinDate     string
}

type LiquidationInputPeriodParams struct {
	DoneMessageText string
	StartDate       LiqudationSelectDateParams
	EndDate         LiqudationSelectDateParams
}
