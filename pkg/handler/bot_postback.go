package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
	"github.com/raahii/haraiai/pkg/timeutil"
)

const (
	POSTBACK_LIQUIDATION_START_DATE = "liquidationStartDate"
	POSTBACK_LIQUIDATION_END_DATE   = "liquidationEndDate"
)

func (bh *BotHandlerImpl) handlePostbackData(event *linebot.Event) error {
	if event.Source.Type != linebot.EventSourceTypeGroup {
		return nil
	}

	groupID := event.Source.GroupID
	group, err := bh.store.GetGroup(groupID)
	if err != nil {
		return err
	}

	switch event.Postback.Data {
	case POSTBACK_LIQUIDATION_START_DATE:
		err = bh.updateLiquidationStartDate(event, group)
	case POSTBACK_LIQUIDATION_END_DATE:
		err = bh.updateLiquidationEndDate(event, group)
	default:
		log.Printf("unhandled postback data found (data=%s, params=%s)\n", event.Postback.Data, event.Postback.Params)
	}

	return err
}

func (bh *BotHandlerImpl) updateLiquidationStartDate(event *linebot.Event, group *store.Group) error {
	date := event.Postback.Params.Date
	liquidation, err := bh.store.GetLiquidation(group.ID)
	if err != nil {
		return fmt.Errorf("%s postback data received but liquidation is not initialized: %w", POSTBACK_LIQUIDATION_START_DATE, err)
	}

	t, err := time.Parse(FULL_DATE_FORMAT, date)
	if err != nil {
		return fmt.Errorf("failed to parse postback data for liquidation startDate (data=%s): %w", date, err)
	}

	// Just set time to 00:00:00.000 in JST for startDate
	startDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, timeutil.JST)

	if liquidation.Period == nil {
		liquidation.Period = new(store.DateRange)
	}
	liquidation.Period.Start = startDate

	err = bh.store.UpdateLiquidation(group.ID, liquidation)
	if err != nil {
		return fmt.Errorf("failed to set startDate to liqdatioin (groupId=%s, startDate=%v): %w", group.ID, startDate, err)
	}

	// Notify date result
	message := linebot.NewTextMessage("開始日: " + date)
	err = bh.bot.ReplyMessage(event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("failed to reply message: %w", err)
	}

	return nil
}

func (bh *BotHandlerImpl) updateLiquidationEndDate(event *linebot.Event, group *store.Group) error {
	date := event.Postback.Params.Date
	liquidation, err := bh.store.GetLiquidation(group.ID)
	if err != nil {
		return fmt.Errorf("%s postback data received but liquidation is not initialized: %w", POSTBACK_LIQUIDATION_START_DATE, err)
	}

	t, err := time.Parse(FULL_DATE_FORMAT, date)
	if err != nil {
		return fmt.Errorf("failed to parse postback data for liquidation endDate (data=%s): %w", date, err)
	}

	// Plus 1day and set time to 00:00:00.000 in JST for endTime
	endDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, timeutil.JST).AddDate(0, 0, 1)

	if liquidation.Period == nil {
		liquidation.Period = new(store.DateRange)
	}
	liquidation.Period.End = endDate

	err = bh.store.UpdateLiquidation(group.ID, liquidation)
	if err != nil {
		return fmt.Errorf("failed to set endDate to liquidation (groupId=%s, endDate=%v): %w", group.ID, endDate, err)
	}

	// Notify date result
	message := linebot.NewTextMessage("終了日: " + date)
	err = bh.bot.ReplyMessage(event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("failed to reply message: %w", err)
	}

	return nil
}
