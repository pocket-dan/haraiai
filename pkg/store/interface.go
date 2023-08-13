//go:generate mockgen -source=$GOFILE -destination=../mock/mock_store.go -package=mock
package store

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
