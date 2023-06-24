package flexmessage

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Template struct {
	FileName string
	Engine   *template.Template
}

const (
	LIQ_SELECT_MODE int = iota
	LIQ_INPUT_PERIOD
	LIQ_CONFIRMATION
)

var templates = map[int]*Template{
	LIQ_SELECT_MODE: {
		FileName: "liquidation/select-mode.json.tmpl",
	},
	LIQ_INPUT_PERIOD: {
		FileName: "liquidation/input-period.json.tmpl",
	},
	LIQ_CONFIRMATION: {
		FileName: "liquidation/confirmation.json.tmpl",
	},
}

func init() {
	packageBasePath := os.Getenv("PACKAGE_BASE_PATH")
	if packageBasePath == "" {
		panic(errors.New("$PACKAGE_BASE_PATH required"))
	}

	templateRootDir := path.Join(packageBasePath, "flexmessage/templates")

	for _, v := range templates {
		p := path.Join(templateRootDir, v.FileName)
		t, err := template.ParseFiles(p)
		if err != nil {
			panic(fmt.Errorf("failed to load message template (path=%s): %w", p, err))
		}
		v.Engine = t
	}
}

// 清算モードを選択するときのメッセージ
type LiquidationModeParams struct {
	SelectWholeModeText   string
	SelectPartialModeText string
}

func BuildLiquidationModeSelectionMessage(params LiquidationModeParams) (*linebot.FlexMessage, error) {
	title := "清算方法を選択するメッセージ"
	return buildFlexMessage(title, LIQ_SELECT_MODE, params)
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

func BuildLiquidationPeriodInputMessage(params LiquidationInputPeriodParams) (*linebot.FlexMessage, error) {
	title := "清算期間を選択できるメッセージ"
	return buildFlexMessage(title, LIQ_INPUT_PERIOD, params)
}

// 清算が終わったか確認するときのメッセージ
type LiquidationConfirmationParams struct {
	OkMessageText string
}

func BuildLiquidationConfirmationMessage(params LiquidationConfirmationParams) (*linebot.FlexMessage, error) {
	title := "清算が完了したか確認するメッセージ"
	return buildFlexMessage(title, LIQ_CONFIRMATION, params)
}

func buildFlexMessage(title string, key int, params interface{}) (*linebot.FlexMessage, error) {
	buf := new(bytes.Buffer)
	t := templates[key]
	t.Engine.Execute(buf, params)

	flexContents, err := linebot.UnmarshalFlexMessageJSON(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to build flexmessage(file=%s, params=%+v): %w", t.FileName, params, err)
	}

	return linebot.NewFlexMessage(title, flexContents), nil
}
