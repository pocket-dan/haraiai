package store

import (
	"context"
	"fmt"
	"time"

	"github.com/raahii/haraiai/pkg/timeutil"
)

const (
	LIQUIDATION_COLLECTION_ID = "liquidations"
)

type DateRange struct {
	Start time.Time
	End   time.Time
}

type Liquidation struct {
	Period    *DateRange `firestore:"period,omitempty"`
	PayerID   string     `firestore:"payer_id"`
	Amount    int64      `firestore:"amount"`
	CreatedAt time.Time  `firestore:"created_at"`
	UpdatedAt time.Time  `firestore:"updated_at"`
}

func (l *Liquidation) IsValidLiquidationPeriod() bool {
	if l.Period == nil {
		return false
	}

	start := l.Period.Start
	end := l.Period.End

	if end.Before(start) || end.Sub(start) >= timeutil.DURATION_MONTH {
		return false
	}

	return true
}

// GetLiquidation find a liquidation by ID.
func (s *StoreImpl) GetLiquidation(groupID string) (*Liquidation, error) {
	ctx := context.Background()

	doc := s.client.Collection(LIQUIDATION_COLLECTION_ID).Doc(groupID)
	docsnap, err := doc.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get liquidation doc ref: %w", err)
	}

	liquidation := new(Liquidation)
	if err := docsnap.DataTo(liquidation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal liquidation data to struct: %w", err)
	}

	return liquidation, nil
}

// CreateLiquidation create a liquidation.
func (s *StoreImpl) CreateLiquidation(groupID string, liquidation Liquidation) error {
	now := timeutil.Now()
	liquidation.CreatedAt = now
	liquidation.UpdatedAt = now

	err := s.saveLiquidation(groupID, &liquidation)
	if err != nil {
		return fmt.Errorf("failed to create liquidation: %w", err)
	}

	return nil
}

// UpdateLiquidation update a liquidation.
func (s *StoreImpl) UpdateLiquidation(groupID string, liquidation *Liquidation) error {
	now := timeutil.Now()
	liquidation.UpdatedAt = now

	err := s.saveLiquidation(groupID, liquidation)
	if err != nil {
		return fmt.Errorf("failed to update liquidation: %w", err)
	}

	return nil
}

func (s *StoreImpl) saveLiquidation(groupID string, liquidation *Liquidation) error {
	ctx := context.Background()

	doc := s.client.Collection(LIQUIDATION_COLLECTION_ID).Doc(groupID)
	_, err := doc.Set(ctx, liquidation)
	return err
}

// DeleteLiquidation delete a liquidation.
func (s *StoreImpl) DeleteLiquidation(groupID string) error {
	ctx := context.Background()

	doc := s.client.Collection(LIQUIDATION_COLLECTION_ID).Doc(groupID)
	_, err := doc.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete liquidation (groupID=%s): %w", groupID, err)
	}

	return nil
}
