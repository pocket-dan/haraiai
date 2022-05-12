package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
)

const (
	TOTAL_UP_TEXT   = "集計"
	TOTAL_UP_PREFIX = "支払った総額は..."

	JOIN_SUFFIX = "だよ"

	START_MESSAGE = "2人の名前を登録できたよ、ありがとう！\n是非試しに「集計」と言ってみてね。"
)

func (bh *BotHandlerImpl) handleTextMessage(event *linebot.Event, message *linebot.TextMessage) error {
	if event.Source.Type != linebot.EventSourceTypeGroup {
		return nil
	}

	groupID := event.Source.GroupID
	group, err := bh.store.GetGroup(groupID)
	if err != nil {
		return err
	}

	if group.Status == store.CREATED {
		if strings.HasSuffix(message.Text, JOIN_SUFFIX) {
			if err := bh.addNewMember(event, group, message.Text); err != nil {
				return err
			}
			return nil
		}
	} else {
		// Total up payment amount for each member.
		if message.Text == TOTAL_UP_TEXT {
			if err := bh.totalUpPayments(event, group); err != nil {
				return err
			}
			return nil
		}

		// Save a new payment if it's valid message.
		if payAmount, err := extractPayAmount(message.Text); err == nil {
			if err := bh.addNewPayment(event, group, payAmount); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (bh *BotHandlerImpl) addNewMember(event *linebot.Event, group *store.Group, text string) error {
	// FIXME: Need to consider multiple users are added to the group simultaneously.

	memberName := strings.TrimSuffix(text, JOIN_SUFFIX)
	memberName = strings.Trim(memberName, " \n")

	senderID := event.Source.UserID
	group.Members[senderID] = store.User{
		ID:   senderID,
		Name: memberName,
	}

	if len(group.Members) == 2 {
		group.Status = store.STARTED
	}

	err := bh.store.SaveGroup(group)
	if err != nil {
		return err
	}

	replyTexts := []string{memberName + "さんだね！👍"}
	if len(group.Members) == 2 {
		replyTexts = append(replyTexts, START_MESSAGE)
	}

	if err = bh.bot.ReplyTextMessage(event.ReplyToken, replyTexts...); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) totalUpPayments(event *linebot.Event, group *store.Group) error {
	replyText := createPayAmountResultMessage(mapToList(group.Members))
	if err := bh.bot.ReplyTextMessage(event.ReplyToken, replyText); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) addNewPayment(event *linebot.Event, group *store.Group, payAmount int) error {
	senderID := event.Source.UserID
	sender, ok := group.Members[senderID]
	if !ok {
		return fmt.Errorf("sender is not found in group (ID=%s)", group.ID)
	}

	sender.PayAmount += int64(payAmount)

	group.Members[sender.ID] = sender

	if err := bh.store.SaveGroup(group); err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	replyText := "👍"
	if err := bh.bot.ReplyTextMessage(event.ReplyToken, replyText); err != nil {
		return err
	}

	return nil
}

func extractPayAmount(text string) (int, error) {
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	if len(lines) != 2 {
		return 0, errors.New("Not supported text due to not 2 lines.")
	}

	const REMOVE_CHARS = " \n\\¥円"
	trimmed := strings.Trim(lines[1], REMOVE_CHARS)

	value, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, fmt.Errorf("Not supported text due to 2nd line is not number: %w", err)
	}

	return value, nil
}

func createPayAmountResultMessage(members []store.User) string {
	lines := []string{TOTAL_UP_PREFIX}
	for _, u := range members {
		lines = append(lines, fmt.Sprintf("%s: %d円", u.Name, u.PayAmount))
	}
	lines = append(lines, "")

	var whoPayALot *store.User
	var whoPayLess *store.User
	if members[0].PayAmount > members[1].PayAmount {
		whoPayALot = &members[0]
		whoPayLess = &members[1]
	} else {
		whoPayALot = &members[1]
		whoPayLess = &members[0]
	}

	var text string
	if whoPayALot.PayAmount == whoPayLess.PayAmount {
		text = "2人とも支払った額は同じだよ！仲良し〜！"
	} else {
		d := (whoPayALot.PayAmount - whoPayLess.PayAmount) / 2
		text = fmt.Sprintf("%sさんが%d円多く払っているよ！", whoPayALot.Name, d)
	}

	lines = append(lines, text)

	return strings.Join(lines, "\n")
}

// TODO: rewrite to generics method and move to common package
func mapToList(m map[string]store.User) []store.User {
	l := make([]store.User, 0, len(m))
	for _, e := range m {
		l = append(l, e)
	}

	return l
}
