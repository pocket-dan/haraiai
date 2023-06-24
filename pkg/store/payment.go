package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/oklog/ulid/v2"
	"github.com/raahii/haraiai/pkg/timeutil"
)

const (
	CreatedAtField = "CreatedAt"
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

// SelectPaymentsBetweenCreatedAt create a payment in the group.
func (s *StoreImpl) SelectPaymentsBetweenCreatedAt(groupID string, period DateRange) (*firestore.DocumentIterator, error) {
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
