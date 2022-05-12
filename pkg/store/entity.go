package store

type Group struct {
	ID      string          `json:"id"`
	Members map[string]User `json:"members"`
	Status  GroupStatus     `json:"status"`
}

type GroupStatus string

const (
	CREATED GroupStatus = "CREATED"
	STARTED GroupStatus = "STARTED"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PayAmount int64  `json:"pay_amount"`
}
