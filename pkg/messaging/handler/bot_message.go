package handler

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/messaging/flexmessage"
	"github.com/raahii/haraiai/pkg/store"
	"github.com/raahii/haraiai/pkg/timeutil"
	"github.com/samber/lo"
)

const (
	JOIN_MESSAGE_SUFFIX = "だよ"

	START_TUTORIAL_MESSAGE          = "使い方を教えて"
	TUTORIAL_PAYMENT_MESSAGE        = "例: お昼ごはん代\n3000"
	TUTORIAL_PAYMENT_CANCEL_MESSAGE = "例: お昼ごはん代\n-3000"

	TOTAL_UP_MESSAGE = "集計"
	HELP_MESSAGE     = "ヘルプ"
	TOTAL_UP_PREFIX  = "支払った総額は..."

	LIQUIDATION_PARTIAL_MESSAGE        = "特定期間のみ清算する"
	LIQUIDATION_CALC_MESSAGE           = "清算額を計算"
	LIQUIDATION_DONE_MESSAGE           = "清算完了"
	LIQUIDATION_PERIOD_INVALID_MESSAGE = "期間が正しく選択されていないか、半年をこえているよ😢\nもう一度期間を選択してね！"

	CHANGE_NAME_MESSAGE_PREFIX = "名前を"
	CHANGE_NAME_MESSAGE_SUFFIX = "に変更"

	DONE_REPLY_MESSAGE = "👍"

	FULL_DATE_FORMAT = "2006-01-02"
)

var (
	MESSAGES_FOR_NAME_CHANGE_GUIDE = []string{
		"名前変更",
		"ニックネーム変更",
		"名前を変更",
		"ニックネームを変更",
		"名前変えて",
		"ニックネーム変えて",
		"名前を変えて",
		"ニックネームを変えて",
		"名前を変えたい",
		"ニックネームを変えたい",
	}

	MESSAGES_FOR_LIQUIDATION = []string{
		"清算",
		"清算したい",
		"精算",
		"精算したい",
		"リセット",
		"リセットしたい",
	}

	READY_TO_START_MESSAGES = []linebot.SendingMessage{
		linebot.NewTextMessage("2人の名前を登録したよ、ありがとう！折半をはじめられるよ。").
			WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("使い方を聞きたい場合はタップ", START_TUTORIAL_MESSAGE),
				),
			)),
	}

	TUTORIAL_REPLYS_1 = []linebot.SendingMessage{
		linebot.NewTextMessage("使い方を説明するよ！\n" +
			"折半したい支払いを記録するときは、まとめて支払った人が「タイトル」と「金額」の2行のメッセージを送ってね！"),
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
	}
)

// Entry point of handing text type webhook event
func (bh *BotHandlerImpl) handleTextMessage(event *linebot.Event, message *linebot.TextMessage) error {
	if event.Source.Type != linebot.EventSourceTypeGroup {
		return nil
	}

	// FIXME: Don't access database every time
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

		// There's no supported commands when group status is GROUP_CREATED.
		return nil
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

	// Liquidation start.
	if lo.Contains(MESSAGES_FOR_LIQUIDATION, message.Text) {
		if err := bh.replyLiquidationStart(event, group); err != nil {
			return err
		}
		return nil
	}

	// Liquidate against payments in a specific period.
	if message.Text == LIQUIDATION_PARTIAL_MESSAGE {
		if err := bh.replyPartialLiquidationStart(event, group); err != nil {
			return err
		}
		return nil
	}

	// Calculate liquidation amount.
	if message.Text == LIQUIDATION_CALC_MESSAGE {
		if err := bh.replyLiquidationAmount(event, group); err != nil {
			return err
		}
		return nil
	}

	// Complete liquidation.
	if message.Text == LIQUIDATION_DONE_MESSAGE {
		if err := bh.replyLiquidationComplete(event, group); err != nil {
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

	// Show guide for name change.
	if lo.Contains(MESSAGES_FOR_NAME_CHANGE_GUIDE, message.Text) {
		if err := bh.replyGuideMessageForNameChange(event); err != nil {
			return err
		}
		return nil
	}

	// Change name.
	if strings.HasPrefix(message.Text, CHANGE_NAME_MESSAGE_PREFIX) && strings.HasSuffix(message.Text, CHANGE_NAME_MESSAGE_SUFFIX) {
		if err := bh.updateMemberName(event, group, message.Text); err != nil {
			return err
		}
		return nil
	}

	// Save a new payment if it's valid message.
	if title, payAmount, err := parsePaymentText(message.Text); err == nil {
		if err := bh.addNewPayment(event, group, title, payAmount); err != nil {
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
	// FIXME: Need to consider nickname (user) validation.

	memberName := strings.TrimSuffix(text, JOIN_MESSAGE_SUFFIX)
	memberName = strings.Trim(memberName, " \n")

	senderID := event.Source.UserID
	group.Members[senderID] = store.NewUser(senderID, memberName, 0)

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
		createPayAmountResultMessage(listMembers(group.Members)),
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

func (bh *BotHandlerImpl) replyLiquidationStart(
	event *linebot.Event,
	group *store.Group,
) error {
	err := bh.store.CreateLiquidation(group.ID, store.Liquidation{}) // upsert actually
	if err != nil {
		return fmt.Errorf("failed to create liquidation(groupId=%s): %w", group.ID, err)
	}

	params := flexmessage.LiquidationModeParams{
		SelectWholeModeText:   LIQUIDATION_CALC_MESSAGE,
		SelectPartialModeText: LIQUIDATION_PARTIAL_MESSAGE,
	}

	replyMessage, err := bh.fs.BuildLiquidationModeSelectionMessage(params)
	if err != nil {
		return err
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessage); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) replyPartialLiquidationStart(
	event *linebot.Event,
	group *store.Group,
) error {
	// Re-initialize liquidation to make liquidation sequence simple.
	err := bh.store.CreateLiquidation(group.ID, store.Liquidation{})
	if err != nil {
		return fmt.Errorf("failed to create liquidation(groupId=%s): %w", group.ID, err)
	}

	serviceStartAt := timeutil.NewDate(2023, 6, 23)

	minDate := timeutil.Max(serviceStartAt, group.CreatedAt).Format(FULL_DATE_FORMAT)
	today := timeutil.Now().Format(FULL_DATE_FORMAT)
	yesterday := timeutil.Now().AddDate(0, 0, -1).Format(FULL_DATE_FORMAT)

	params := flexmessage.LiquidationInputPeriodParams{
		DoneMessageText: LIQUIDATION_CALC_MESSAGE,
		StartDate: flexmessage.LiqudationSelectDateParams{
			Data:        POSTBACK_LIQUIDATION_START_DATE,
			InitialDate: yesterday,
			MinDate:     minDate,
			MaxDate:     today,
		},
		EndDate: flexmessage.LiqudationSelectDateParams{
			Data:        POSTBACK_LIQUIDATION_END_DATE,
			InitialDate: today,
			MinDate:     minDate,
			MaxDate:     today,
		},
	}

	flexMessage, err := bh.fs.BuildLiquidationPeriodInputMessage(params)
	if err != nil {
		return err
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, flexMessage); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) replyLiquidationAmount(
	event *linebot.Event,
	group *store.Group,
) error {
	liquidation, err := bh.store.GetLiquidation(group.ID)
	if err != nil {
		return err
	}

	var liquidationAmount int64
	var whoPayALot *store.User
	var whoPayLess *store.User
	if liquidation.Period == nil {
		liquidationAmount, whoPayALot, whoPayLess = bh.calcWholeLiquidationAmount(group)
	} else {
		var err error
		if !liquidation.IsValidLiquidationPeriod() {
			textMessage := linebot.NewTextMessage(LIQUIDATION_PERIOD_INVALID_MESSAGE)
			if err := bh.bot.ReplyMessage(event.ReplyToken, textMessage); err != nil {
				return err
			}
			return nil
		}
		liquidationAmount, whoPayALot, whoPayLess, err = bh.calcPartialLiquidationAmount(group, liquidation.Period)
		if err != nil {
			return fmt.Errorf("failed to calculate liquidation amount in the period: %w", err)
		}
	}

	if liquidationAmount == 0 {
		bh.store.DeleteLiquidation(group.ID)

		textMessage := linebot.NewTextMessage("払った額は同じ！清算の必要はないよ")

		if err := bh.bot.ReplyMessage(event.ReplyToken, textMessage); err != nil {
			return err
		}
		return nil
	}

	liquidation.Amount = liquidationAmount
	liquidation.PayerID = whoPayLess.ID
	err = bh.store.UpdateLiquidation(group.ID, liquidation)
	if err != nil {
		return fmt.Errorf("failed to update partial liquidation (groupId=%s): %w", group.ID, err)
	}

	params := flexmessage.LiquidationConfirmationParams{
		OkMessageText: LIQUIDATION_DONE_MESSAGE,
	}
	confirmationMessage, err := bh.fs.BuildLiquidationConfirmationMessage(params)
	if err != nil {
		return err
	}

	text := fmt.Sprintf("%sさんは%sさんに %d 円渡してね🙏", whoPayLess.Name, whoPayALot.Name, liquidationAmount)
	replyMessages := []linebot.SendingMessage{
		linebot.NewTextMessage(text),
		confirmationMessage,
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessages...); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) calcWholeLiquidationAmount(
	group *store.Group,
) (int64, *store.User, *store.User) {
	members := listMembers(group.Members)
	sortUsersByPayAmountDesc(members)
	whoPayALot, whoPayLess := members[0], members[1]
	liquidationAmount := (whoPayALot.PayAmount - whoPayLess.PayAmount) / 2
	return liquidationAmount, whoPayALot, whoPayLess
}

func (bh *BotHandlerImpl) calcPartialLiquidationAmount(
	group *store.Group,
	period *store.DateRange,
) (int64, *store.User, *store.User, error) {
	members := listMembers(group.Members)

	payAmountMap, err := bh.store.BuildPayAmountMapBetweenCreatedAt(group.ID, period)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to select and build pay amount map (groupID=%s): %w", group.ID, err)
	}

	userA := members[0]
	amountA := payAmountMap[userA.ID]

	userB := members[1]
	amountB := payAmountMap[userB.ID]

	if amountA > amountB {
		d := (amountA - amountB) / 2
		return d, userA, userB, nil
	} else {
		d := (amountB - amountA) / 2
		return d, userB, userA, nil
	}
}

func (bh *BotHandlerImpl) replyLiquidationComplete(
	event *linebot.Event,
	group *store.Group,
) error {
	// Get liquidation
	liquidation, err := bh.store.GetLiquidation(group.ID)
	if err != nil || liquidation.Amount <= 0 {
		bh.store.DeleteLiquidation(group.ID)
		return nil
	}

	// FIXME: operate group and payment using a transaction

	// Liquidation
	payer := group.Members[liquidation.PayerID]
	payer.PayAmount += liquidation.Amount * 2
	payer.Touch()
	if err := bh.store.SaveGroup(group); err != nil {
		return err
	}

	// Record payment
	payment := new(store.Payment)
	payment.Name = "清算"
	payment.Amount = liquidation.Amount
	payment.PayerID = liquidation.PayerID
	payment.Type = store.PAYMENT_TYPE_LIQUIDATION

	if err := bh.store.CreatePayment(group.ID, payment); err != nil {
		return err
	}

	bh.store.DeleteLiquidation(group.ID)

	replyMessage := linebot.NewTextMessage(DONE_REPLY_MESSAGE)
	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessage); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) addNewPayment(event *linebot.Event, group *store.Group, title string, amount int) error {
	senderID := event.Source.UserID
	sender, ok := group.Members[senderID]
	if !ok {
		return fmt.Errorf("sender is not found in group (ID=%s)", group.ID)
	}

	// FIXME: operate group and payment using a transaction

	// Record payment
	payment := new(store.Payment)
	payment.Name = title
	payment.Amount = int64(amount)
	payment.Type = store.PAYMENT_TYPE_DEFAULT
	payment.PayerID = senderID

	if err := bh.store.CreatePayment(group.ID, payment); err != nil {
		return err
	}

	// Plus payAmount
	sender.PayAmount += int64(amount)
	sender.Touch()

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

func (bh *BotHandlerImpl) replyGuideMessageForNameChange(event *linebot.Event) error {
	message := fmt.Sprintf(
		"名前を変更したいときは\n「%s○○%s」\nと言ってね！",
		CHANGE_NAME_MESSAGE_PREFIX, CHANGE_NAME_MESSAGE_SUFFIX,
	)
	replyMessage := []linebot.SendingMessage{
		linebot.NewTextMessage(message),
	}

	if err := bh.bot.ReplyMessage(event.ReplyToken, replyMessage...); err != nil {
		return err
	}

	return nil
}

func (bh *BotHandlerImpl) updateMemberName(
	event *linebot.Event,
	group *store.Group,
	messageText string,
) error {
	senderID := event.Source.UserID
	sender, ok := group.Members[senderID]
	if !ok {
		return fmt.Errorf("sender is not found in group (ID=%s)", group.ID)
	}

	sender.Name = extractNewName(messageText)
	sender.Touch()

	err := bh.store.SaveGroup(group)
	if err != nil {
		return err
	}

	replyMessage := []linebot.SendingMessage{
		linebot.NewTextMessage(fmt.Sprintf("名前を「%s」に変更しました👍", sender.Name)),
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

func parsePaymentText(text string) (string, int, error) {
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	if len(lines) != 2 {
		return "", 0, errors.New("Not supported text due to not 2 lines.")
	}

	title := parsePaymentTitle(lines[0])
	amount, err := parsePayAmount(lines[1])
	if err != nil {
		return "", 0, fmt.Errorf("2nd line text is not number: %w", err)
	}

	return title, amount, nil
}

func parsePaymentTitle(text string) string {
	return strings.Trim(text, "\n ")
}

func parsePayAmount(text string) (int, error) {
	trimmed := strings.Trim(text, " \n\\¥円")

	value, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse '%s' as integer: %w", text, err)
	}

	return value, nil
}

func createPayAmountResultMessage(members []*store.User) string {
	lines := []string{TOTAL_UP_PREFIX}
	for _, u := range members {
		lines = append(lines, fmt.Sprintf("%s: %d円", u.Name, u.PayAmount))
	}
	lines = append(lines, "")

	sortUsersByPayAmountDesc(members)
	whoPayALot, whoPayLess := members[0], members[1]

	var text string
	if whoPayALot.PayAmount == whoPayLess.PayAmount {
		text = "2人とも支払った額は同じだよ！仲良し〜！"
	} else {
		d := (whoPayALot.PayAmount - whoPayLess.PayAmount) / 2
		text = fmt.Sprintf("%sさんが %d 円多く払っているよ。", whoPayALot.Name, d)
		text += fmt.Sprintf("次は%sさんが払うと距離が縮まるね🤝", whoPayLess.Name)
	}

	lines = append(lines, text)

	return strings.Join(lines, "\n")
}

func sortUsersByPayAmountDesc(users []*store.User) {
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].PayAmount > users[j].PayAmount
	})
}

func listMembers(m map[string]*store.User) []*store.User {
	v := make([]*store.User, 0, len(m))
	for _, e := range m {
		v = append(v, e)
	}

	sort.Slice(v, func(i, j int) bool { return v[i].Name < v[j].Name })
	return v
}

func extractNewName(text string) string {
	textRunes := []rune(text)
	prefixCharCounts := utf8.RuneCountInString(CHANGE_NAME_MESSAGE_PREFIX)
	suffixCharCounts := utf8.RuneCountInString(CHANGE_NAME_MESSAGE_SUFFIX)

	start := prefixCharCounts
	end := len(textRunes) - suffixCharCounts

	extracted := string(textRunes[start:end])
	extracted = strings.Trim(extracted, " \n")

	return extracted
}
