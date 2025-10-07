package interfaces

import (
	dbm "backend/internal/db_model"
	"context"

	"github.com/google/uuid"
)

// FoodOrderRepository คือสัญญาที่ usecase ต้องพึ่ง
type FoodOrderRepository interface {
	CreateOrderWithItems(ctx context.Context,
		order dbm.FoodOrder,
		items []dbm.FoodOrderItem,
		options []dbm.FoodOrderItemOption,
	) (*dbm.FoodOrder, error)

	GetOrderWithItems(ctx context.Context, orderID uuid.UUID) (*dbm.FoodOrder, error)
}
