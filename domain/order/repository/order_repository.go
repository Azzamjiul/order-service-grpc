package repository

import (
	"order-service/domain/order"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (ur *OrderRepository) Create(order *order.Order) error {
	return ur.db.Create(order).Error
}

func (ur *OrderRepository) FindByID(id uint) (*order.Order, error) {
	var order order.Order
	if err := ur.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
