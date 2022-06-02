package handler

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
)

const (
	JOIN_MESSAGE_SUFFIX = "だよ"

	START_TUTORIAL_MESSAGE          = "使い方を教えて"
	TUTORIAL_PAYMENT_MESSAGE        = "例: お昼ごはん代\n3000"
	TUTORIAL_PAYMENT_CANCEL_MESSAGE = "例: お昼ごはん代\n-3000"

	TOTAL_UP_MESSAGE         = "集計"
	EVEN_UP_MESSAGE          = "精算"
	EVEN_UP_COMPLETE_MESSAGE = "精算完了"
	HELP_MESSAGE             = "ヘルプ"
	TOTAL_UP_PREFIX          = "支払った総額は..."

	DONE_REPLY_MESSAGE = "👍"
)

var (
	READY_TO_START_MESSAGES = []linebot.SendingMessage{
		linebot.NewTextMessage("2人の名前を登録したよ、ありがとう！割り勘をはじめられるよ。").
			WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("使い方を聞きたい場合はタップ", START_TUTORIAL_MESSAGE),
				),
			)),
	}

	TUTORIAL_REPLYS_1 = []linebot.SendingMessage{
		linebot.NewTextMessage("使い方を説明するよ！\n" +
			"割り勘したいときは、まとめて支払った人が「タイトル」と「金額」の2行のメッセージを送ってね！"),
		linebot.NewTextMessage(TUTORIAL_PAYMENT_MESSAGE).
			WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("例のメッセージを送ってみる", TUTORIAL_PAYMENT_MESSAGE),
				),
			)),
	}

	TUTORIAL_REPLYS_2 = []linebot.SendingMessage{
		linebot.NewTextMessage("支払い状況を確認したい場合は「集計」とメッセージを送ってみてね。").
			WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("集計 と送る", TOTAL_UP_MESSAGE),
				),
			)),
	}

	TUTORIAL_REPLYS_3 = []linebot.SendingMessage{
		linebot.NewTextMessage("もし過去の支払いを取り消したい場合は、下の例のようにマイナスで打ち消してね。"),
		linebot.NewTextMessage(TUTORIAL_PAYMENT_CANCEL_MESSAGE).
			WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("取り消しの例のメッセージを送る", TUTORIAL_PAYMENT_CANCEL_MESSAGE),
				),
			)),
	}

	TUTORIAL_REPLYS_4 = []linebot.SendingMessage{
		linebot.NewTextMessage("お疲れさまでした！使い方の説明はおしまいです！😄"),
		linebot.NewTextMessage("わからないことがあったら ヘルプ と声をかけてね"),
		// linebot.NewTextMessage("最後に haraiai には支払いを精算してリセットする機能はないよ。定期的な精算をするよりも、支払いが少ない側が次回多めに払うことで支払い額のバランスを保つようにしよう！"),
	}
)

// Entry point of handing text type webhook event
func (bh *BotHandlerImpl) handleTextMessage(event *linebot.Event, message *linebot.TextMessage) error {
	if event.Source.Type != linebot.EventSourceTypeGroup {
		return nil
	}

	groupID := event.Source.GroupID
	group, err := bh.store.GetGroup(groupID)
	if err != nil {
		return err
	}

	if group.Status == store.GROUP_CREATED {
		// This group is under setting up as it doesn't have sufficient members.
		// Handle messages to join to group only.
		if strings.HasSuffix(message.Text, JOIN_MESSAGE_SUFFIX) {
			if err := bh.addNewMember(event, group, message.Text); err != nil {
				return err
			}
			return nil
		}
	}

	// Tutorial
	if message.Text == START_TUTORIAL_MESSAGE {
		group.IsTutorial = true
		if err := bh.store.SaveGroup(group); err != nil {
			return err
		}

		if err := bh.bot.ReplyMessage(event.ReplyToken, TUTORIAL_REPLYS_1...); err != nil {
			return err
		}
		return nil
	}

	// Total up payment amount for each member.
	if message.Text == TOTAL_UP_MESSAGE {
		if err := bh.replyTotalUpResult(event, group); err != nil {
			return err
		}
		return nil
	}

	// Even up payment amount.
	if message.Text == EVEN_UP_MESSAGE {
		if err := bh.replyEvenUpConfirmation(event, group); err != nil {
			return err
		}
		return nil
	}

	// Complete even up payment amount.
	if message.Text == EVEN_UP_COMPLETE_MESSAGE {
		if err := bh.replyEvenUpComplete(event, group); err != nil {
			return err
		}
		return nil
	}

	// Show guide for help.
	if message.Text == HELP_MESSAGE {
		if err := bh.replyHelpMessage(event); err != nil {
			return err
		}
		return nil
	}

	// Save a new payment if it's valid message.
	if payAmount, err := extractPayAmount(message.Text); err == nil {
		if err := bh.addNewPayment(event, group, payAmount); err != nil {
			return err
		}

		if err := bh.replyToNewPayment(event, message.Text); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (bh *BotHandlerImpl) addNewMember(event *linebot.Event, group *store.Group, text string) error {
	// FIXME: Need to consider multiple users are added to the group simultaneously.

	memberName := strings.TrimSuffix(text, JOIN_MESSAGE_SUFFIX)
	memberName = strings.Trim(memberName, " \n")

	senderID := event.Source.UserID
	group.Members[senderID] = store.User{
		ID:   senderID,
		Name: memberName,
	}

	if len(group.Members) == 2 {
		group.Status = store.GROUP_STARTED
	}

	err := bh.store.SaveGroup(group)
	if err != nil {
		return err
	}

	replyMessages := []linebot.SendingMessage{
		linebot.NewTextMessage(memberName + "さんだね！👍"),
	}

	if len(group.Members) == 2 {
		replyMessages = append(replyMessages, READY_TO_START_MESSAGES...)
	}

	if err = bh.bot.ReplyMessage(event.ReplyToken, replyMessages...); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) replyTotalUpResult(
	event *linebot.Event,
	group *store.Group,
) error {
	replyMessages := []linebot.SendingMessage{}

	replyMessages = append(replyMessages, linebot.NewTextMessage(
		createPayAmountResultMessage(mapToList(group.Members)),
	))

	if group.IsTutorial {
		group.IsTutorial = false
		err := bh.store.SaveGroup(group)
		if err != nil {
			return err
		}

		replyMessages = append(replyMessages, TUTORIAL_REPLYS_3...)
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessages...); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) replyEvenUpConfirmation(
	event *linebot.Event,
	group *store.Group,
) error {
	members := mapToList(group.Members)
	sortUsersByPayAmountDesc(members)

	whoPayALot := &members[0]
	whoPayLess := &members[1]

	var replyMessage linebot.SendingMessage
	if whoPayALot.PayAmount == whoPayLess.PayAmount {
		replyMessage = linebot.NewTextMessage("払った額は同じ！精算の必要はないよ")
	} else {
		d := (whoPayALot.PayAmount - whoPayLess.PayAmount) / 2
		text := fmt.Sprintf("%s は %s に %d 円払うと精算完了です。精算しましたか？", whoPayLess.Name, whoPayALot.Name, d)
		replyMessage = linebot.NewTextMessage(text).WithQuickReplies(
			linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("はい", EVEN_UP_COMPLETE_MESSAGE),
				),
			),
		)
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessage); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) replyEvenUpComplete(
	event *linebot.Event,
	group *store.Group,
) error {
	members := mapToList(group.Members)
	sortUsersByPayAmountDesc(members)
	whoPayALot := &members[0]
	whoPayLess := &members[1]

	if whoPayALot.PayAmount == whoPayLess.PayAmount {
		return nil
	}

	whoPayLess.PayAmount = whoPayALot.PayAmount
	group.Members[whoPayLess.ID] = *whoPayLess

	if err := bh.store.SaveGroup(group); err != nil {
		return err
	}

	replyMessage := linebot.NewTextMessage(DONE_REPLY_MESSAGE)
	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessage); err != nil {
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

	return nil
}

func (bh *BotHandlerImpl) replyHelpMessage(event *linebot.Event) error {
	replyMessage := []linebot.SendingMessage{
		linebot.NewTextMessage("ヘルプページはこちら:\n" + bh.config.GetHelpPageURL()),
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessage...); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) replyToNewPayment(event *linebot.Event, text string) error {
	replyMessages := []linebot.SendingMessage{
		linebot.NewTextMessage(DONE_REPLY_MESSAGE),
	}

	// For tutorial.
	if text == TUTORIAL_PAYMENT_MESSAGE {
		replyMessages = append(replyMessages, TUTORIAL_REPLYS_2...)
	} else if text == TUTORIAL_PAYMENT_CANCEL_MESSAGE {
		replyMessages = append(replyMessages, TUTORIAL_REPLYS_4...)
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessages...); err != nil {
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

	sortUsersByPayAmountDesc(members)
	whoPayALot := &members[0]
	whoPayLess := &members[1]

	var text string
	if whoPayALot.PayAmount == whoPayLess.PayAmount {
		text = "2人とも支払った額は同じだよ！仲良し〜！"
	} else {
		d := whoPayALot.PayAmount - whoPayLess.PayAmount
		text = fmt.Sprintf("%s は今度 %d 円分支払うと追いつくよ🙌", whoPayLess.Name, d)
	}

	lines = append(lines, text)

	return strings.Join(lines, "\n")
}

func sortUsersByPayAmountDesc(users []store.User) {
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].PayAmount > users[j].PayAmount
	})
}

func mapToList(m map[string]store.User) []store.User {
	// TODO: rewrite to generics method and move to common package

	v := make([]store.User, 0, len(m))
	for _, e := range m {
		v = append(v, e)
	}

	sort.Slice(v, func(i, j int) bool { return v[i].Name < v[j].Name })
	return v
}
