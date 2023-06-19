package store

import (
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
	PAYMENT_TYPE_DEFAULT PaymentType = "DEFAULT" // 通常の支払い
	PAYMENT_TYPE_EVEN_UP PaymentType = "EVEN_UP" // 清算
)

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
