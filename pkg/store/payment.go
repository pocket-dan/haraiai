package store

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/oklog/ulid/v2"
	"github.com/raahii/haraiai/pkg/timeutil"
	"google.golang.org/api/iterator"
)

const (
	CreatedAtField = "CreatedAt"
)

// Payment
type Payment struct {
	ID        string
	Name      string
	Amount    int64
	Type      PaymentType
	PayerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentType string

const (
	PAYMENT_TYPE_DEFAULT     PaymentType = "DEFAULT"     // 通常の支払い
	PAYMENT_TYPE_LIQUIDATION PaymentType = "LIQUIDATION" // 清算
)

// CreatePayment create a payment in the group.
func (s *StoreImpl) CreatePayment(groupID string, payment *Payment) error {
	payment.ID = ulid.Make().String()

	now := timeutil.Now()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	ctx := context.Background()

	doc := s.client.
		Collection(GROUP_COLLECTION_ID).Doc(groupID).
		Collection(PAYMENT_COLLECTION_ID).Doc(payment.ID)

	_, err := doc.Set(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to add payment to group(id=%s): %w", groupID, err)
	}

	return nil
}

// BuildPayAmountMapBetweenCreatedAt select payments and build pay amount map for group member.
func (s *StoreImpl) BuildPayAmountMapBetweenCreatedAt(groupID string, period *DateRange) (map[string]int64, error) {
	// storeにロジックが入っていて気持ち悪いが、handler側でselectPaymentsBetweenCreatedAtを呼ぶと、
	// DocumentIteratorのモックが必要になって面倒なのでやむを得ず
	payAmountMap := map[string]int64{}
	iter, err := s.selectPaymentsBetweenCreatedAt(groupID, period)
	if err != nil {
		return payAmountMap, err
	}

	for {
		docsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return payAmountMap, err
		}

		payment := new(Payment)
		if err := docsnap.DataTo(payment); err != nil {
			return payAmountMap, fmt.Errorf("failed to unmarshal payment data to struct: %w", err)
		}
		if payment.Type == PAYMENT_TYPE_LIQUIDATION {
			continue
		}

		payAmountMap[payment.PayerID] += payment.Amount
	}

	return payAmountMap, nil
}

func (s *StoreImpl) selectPaymentsBetweenCreatedAt(groupID string, period *DateRange) (*firestore.DocumentIterator, error) {
	ctx := context.Background()

	docs := s.client.
		Collection(GROUP_COLLECTION_ID).
		Doc(groupID).
		Collection(PAYMENT_COLLECTION_ID).
		Where(CreatedAtField, ">=", period.Start).
		Where(CreatedAtField, "<=", period.End).
		Documents(ctx)

	return docs, nil
}
