//go:generate mockgen -source=$GOFILE -destination=../mock/store_$GOFILE -package=mock
package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

const (
	GROUP_COLLECTION_ID = "groups"
)

type Store interface {
	GetGroup(string) (*Group, error)
	SaveGroup(*Group) error
	DeleteGroup(string) error
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

// GetGroup find a group from ID.
func (s *StoreImpl) GetGroup(groupID string) (*Group, error) {
	ctx := context.Background()

	doc := s.client.Collection(GROUP_COLLECTION_ID).Doc(groupID)
	docsnap, err := doc.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get group doc ref: %w", err)
	}

	group := Group{}
	if err := docsnap.DataTo(&group); err != nil {
		return nil, fmt.Errorf("failed to unmarshal group data to struct: %w", err)
	}

	return &group, nil
}

// SaveGroup update a group.
func (s *StoreImpl) SaveGroup(group *Group) error {
	ctx := context.Background()

	doc := s.client.Collection(GROUP_COLLECTION_ID).Doc(group.ID)
	_, err := doc.Set(ctx, group)
	if err != nil {
		log.Printf("failed to save group %+v: %v\n", group, err)
		return fmt.Errorf("failed to create group: %w", err)
	}

	return nil
}

// DeleteGroup delete a group.
func (s *StoreImpl) DeleteGroup(groupID string) error {
	ctx := context.Background()

	doc := s.client.Collection(GROUP_COLLECTION_ID).Doc(groupID)
	_, err := doc.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}
