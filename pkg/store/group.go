package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/raahii/haraiai/pkg/timeutil"
)

// Group
type Group struct {
	ID         string           `json:"id"`
	Members    map[string]*User `json:"members"`
	Status     GroupStatus      `json:"status"`
	IsTutorial bool             `json:"is_tutorial"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type GroupStatus string

const (
	GROUP_CREATED GroupStatus = "GROUP_CREATED"
	GROUP_STARTED GroupStatus = "GROUP_STARTED"
)

func NewGroup(ID string, status GroupStatus) *Group {
	g := new(Group)
	g.ID = ID
	g.Status = status
	g.IsTutorial = false

	g.Members = map[string]*User{}

	now := timeutil.Now()
	g.CreatedAt = now
	g.UpdatedAt = now

	return g
}

// User
type User struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	PayAmount        int64     `json:"pay_amount"`
	InitialPayAmount int64     `json:"initial_pay_amount"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewUser(ID, name string, payAmount int64) *User {
	u := new(User)
	u.ID = ID
	u.Name = name
	u.PayAmount = payAmount

	now := timeutil.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	return u
}

func (u *User) Touch() {
	u.UpdatedAt = timeutil.Now()
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
