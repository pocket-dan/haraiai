package flexmessage

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/messaging/config"
)

const (
	LIQ_SELECT_MODE int = iota
	LIQ_INPUT_PERIOD
	LIQ_CONFIRMATION
)

type FlexMessageBuilderImpl struct {
	templates map[int]*Template
}

func ProvideFlexMessageBuilder(bc config.BotConfig) FlexMessageBuilder {
	templates := map[int]*Template{
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

	for _, v := range templates {
		p := path.Join(bc.GetFlexTemplateDir(), v.FileName)
		t, err := template.ParseFiles(p)
		if err != nil {
			panic(fmt.Errorf("failed to load message template (path=%s): %w", p, err))
		}
		v.Engine = t
	}

	return &FlexMessageBuilderImpl{
		templates: templates,
	}
}

type Template struct {
	FileName string
	Engine   *template.Template
}

func (b *FlexMessageBuilderImpl) BuildLiquidationModeSelectionMessage(params LiquidationModeParams) (*linebot.FlexMessage, error) {
	title := "清算方法を選択するメッセージ"
	return b.buildFlexMessage(title, LIQ_SELECT_MODE, params)
}

func (b *FlexMessageBuilderImpl) BuildLiquidationPeriodInputMessage(params LiquidationInputPeriodParams) (*linebot.FlexMessage, error) {
	title := "清算期間を選択できるメッセージ"
	return b.buildFlexMessage(title, LIQ_INPUT_PERIOD, params)
}

func (b *FlexMessageBuilderImpl) BuildLiquidationConfirmationMessage(params LiquidationConfirmationParams) (*linebot.FlexMessage, error) {
	title := "清算が完了したか確認するメッセージ"
	return b.buildFlexMessage(title, LIQ_CONFIRMATION, params)
}

func (b *FlexMessageBuilderImpl) buildFlexMessage(title string, key int, params interface{}) (*linebot.FlexMessage, error) {
	buf := new(bytes.Buffer)
	t := b.templates[key]
	t.Engine.Execute(buf, params)

	flexContents, err := linebot.UnmarshalFlexMessageJSON(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to build flexmessage(file=%s, params=%+v): %w", t.FileName, params, err)
	}

	return linebot.NewFlexMessage(title, flexContents), nil
}
