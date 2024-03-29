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

type StoreImpl struct {
	client *firestore.Client
}

func ProvideStore() (Store, error) {
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
