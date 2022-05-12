package handler

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/raahii/haraiai/pkg/store"
)

const (
	NOT_SUPPORTED_MESSAGE string = `すみません、このトークで払い合いをお使いいただけません。
  お手数ですが、グループを作成していただき、再度追加してください 🙇`

	// GREETING_MESSAGE string = `グループへの追加ありがとう！僕は皆さんの割り勘をサポートします 🤝
	//
	// 【はじめに】
	// 割り勘したいメンバーをまず全員このグループに招待してから始めてね。
	//
	// 【使い方① 支払いを記録】
	// 使い方はとっても簡単！まずは、まとめて支払いをした人が次のように2行のメッセージを送ってね。
	//
	// 夜ご飯代
	// 2000
	//
	// 1行目にお好きなタイトル、2行目に金額を入力してね！このルールから外れたメッセージは無視されてしまうので注意！
	//
	// 【使い方② 集計結果を確認】
	// そして今だれがいくら払ったかを見たいときは次のように一言メッセージを送ってね！
	//
	// 集計
	//
	// 「集計」と送ると、皆さんがそれぞれ払った金額が表示されるよ。
	// 集計結果を見ながら支払いが少ない人が次は払うよう調整してね。`

	GREETING_MESSAGE string = `グループへの追加ありがとう！haraiai が二人の割り勘をサポートするよ🤝
  始めるにあたって2人の名前を教えてね。短い呼び名のほうがきれいに表示できるよ！

  ○○だよ

  と答えてね。`
)

func (bh *BotHandlerImpl) handleBotJoin(event *linebot.Event) error {
	if event.Source.Type != linebot.EventSourceTypeGroup {
		// Currently, support group talk only.
		err := bh.bot.ReplyTextMessage(event.ReplyToken, NOT_SUPPORTED_MESSAGE)
		if err != nil {
			return err
		}
		return nil
	}

	// Initialize group data.
	group := &store.Group{
		ID:      event.Source.GroupID,
		Members: map[string]store.User{},
		Status:  store.CREATED,
	}

	err := bh.store.SaveGroup(group)
	if err != nil {
		return err
	}

	// Send greeting message.
	err = bh.bot.ReplyTextMessage(event.ReplyToken, GREETING_MESSAGE)
	if err != nil {
		return err
	}

	return nil
}

// func (bh *BotHandlerImpl) newGroup(groupID string) (*store.Group, error) {
// 	members := make(map[string]store.User, len(memberIDs))
// 	for i := 0; i < len(memberIDs); i++ {
// 		members[memberIDs[i]] = store.User{
// 			ID:        memberIDs[i],
// 			Name:      memberNames[i],
// 			PayAmount: 0,
// 		}
// 	}
//
// 	group := &store.Group{
// 		ID:          groupID,
// 		MemberCount: len(members),
// 		Members:     members,
// 	}
//
// 	return group, nil
// }
