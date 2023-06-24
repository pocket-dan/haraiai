package store

import (
	"context"
	"fmt"
	"log"

	"github.com/raahii/haraiai/pkg/timeutil"
)

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
	group.UpdatedAt = timeutil.Now()

	ctx := context.Background()

	doc := s.client.Collection(GROUP_COLLECTION_ID).Doc(group.ID)
	_, err := doc.Set(ctx, group)
	if err != nil {
		log.Printf("failed to save group %+v: %v\n", group, err)
		return fmt.Errorf("failed to save group: %w", err)
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
