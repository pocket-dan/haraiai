package store

type Group struct {
	ID         string          `json:"id"`
	Members    map[string]User `json:"members"`
	Status     GroupStatus     `json:"status"`
	IsTutorial bool            `json:"is_tutorial"`
}

type GroupStatus string

const (
	GROUP_CREATED GroupStatus = "GROUP_CREATED"
	GROUP_STARTED GroupStatus = "GROUP_STARTED"
)

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
	ID        string `json:"id"`
	Name      string `json:"name"`
	PayAmount int64  `json:"pay_amount"`
}

func NewUser(ID, Name string, PayAmount int64) *User {
	u := new(User)
	u.ID = ID
	u.Name = Name
	u.PayAmount = PayAmount

	return u
}
