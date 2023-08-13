package handler

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (bh *BotHandlerImpl) createRichMenu() error {
	richMenu := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 800, Height: 540},
		Selected:    true,
		Name:        "Default Menu",
		ChatBarText: "haraiai について",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 800, Height: 270},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  bh.config.GetAboutPageURL(),
					Text: "サービス概要ページへ飛ぶ",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 270, Width: 266, Height: 270},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  bh.config.GetHelpPageURL(),
					Text: "よくある質問のページへ飛ぶ",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 267, Y: 270, Width: 267, Height: 270},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  bh.config.GetInquiryPageURL(),
					Text: "お問い合わせページへ飛ぶ",
				},
			},
		},
	}

	richMenuID, err := bh.bot.CreateRichMenu(richMenu)
	if err != nil {
		return err
	}

	if err := bh.bot.UploadRichMenuImage(richMenuID, bh.config.GetRichMenuImagePath()); err != nil {
		return err
	}

	if err := bh.bot.SetDefaultRichMenu(richMenuID); err != nil {
		return err
	}

	return nil
}
