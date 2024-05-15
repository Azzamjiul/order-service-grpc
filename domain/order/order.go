package order

type Order struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type OrderRepository interface {
	Create(order *Order) error
	FindByID(id uint) (*Order, error)
}
