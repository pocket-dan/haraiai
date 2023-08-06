//go:generate mockgen -source=$GOFILE -destination=../mock/store_$GOFILE -package=mock
package store

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

const (
	GROUP_COLLECTION_ID   = "groups"
	PAYMENT_COLLECTION_ID = "payments"
)

type Store interface {
	// Group
	GetGroup(string) (*Group, error)
	SaveGroup(*Group) error
	DeleteGroup(string) error
	// Payment
	CreatePayment(string, *Payment) error
	BuildPayAmountMapBetweenCreatedAt(string, *DateRange) (map[string]int64, error)
	// Liquidation
	GetLiquidation(string) (*Liquidation, error)
	CreateLiquidation(string, Liquidation) error
	UpdateLiquidation(string, *Liquidation) error
	DeleteLiquidation(string) error
}

type StoreImpl struct {
	client *firestore.Client
}

func New() (*StoreImpl, error) {
	// Initialize Cloud Firestore.
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		return nil, errors.New("$PROJECT_ID required.")
	}

	appPhase := os.Getenv("PHASE")
	if appPhase == "" {
		return nil, errors.New("$PHASE required")
	}

	var client *firestore.Client
	var err error
	ctx := context.Background()
	if appPhase == "local" {
		credentialPath := os.Getenv("FIRESTORE_CREDENTIALS")
		if credentialPath == "" {
			return nil, errors.New("$FIRESTORE_CREDENTIALS requried for local development")
		}

		opt := option.WithCredentialsFile(credentialPath)
		client, err = firestore.NewClient(ctx, projectID, opt)
	} else {
		client, err = firestore.NewClient(ctx, projectID)
	}

	if err != nil {
		return nil, err
	}

	return &StoreImpl{client}, nil
}
