package store

import "time"

var JST *time.Location

func init() {
	var err error
	JST, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

type Group struct {
	ID         string          `json:"id"`
	Members    map[string]User `json:"members"` // TODO: Change type to map[string]*User.
	Status     GroupStatus     `json:"status"`
	IsTutorial bool            `json:"is_tutorial"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type GroupStatus string

const (
	GROUP_CREATED GroupStatus = "GROUP_CREATED"
	GROUP_STARTED GroupStatus = "GROUP_STARTED"
)

// TODO: Members parameter is not necessary.
func NewGroup(ID string, status GroupStatus, members []User) *Group {
	g := new(Group)
	g.ID = ID
	g.Status = status
	g.IsTutorial = false

	g.Members = make(map[string]User, len(members))
	for _, u := range members {
		g.Members[u.ID] = u
	}

	return g
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PayAmount int64     `json:"pay_amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(ID, name string, payAmount int64) *User {
	u := new(User)
	u.ID = ID
	u.Name = name
	u.PayAmount = payAmount

	u.CreatedAt = nowInJST()
	u.UpdatedAt = nowInJST()

	return u
}

func (u *User) Touch() {
	u.UpdatedAt = nowInJST()
}

func nowInJST() time.Time {
	// TZ environment variable is set, but also set in code.
	return time.Now().In(JST)
}
