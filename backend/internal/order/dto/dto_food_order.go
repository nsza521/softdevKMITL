package dto

import (
	"time"

	"github.com/google/uuid"
)

type FoodOrderItemOptionDTO struct {
	ID         uuid.UUID `json:"id"`
	OptionID   uuid.UUID `json:"option_id"`
	Qty        int       `json:"qty"`
	PriceDelta float64   `json:"price_delta"`
}

type FoodOrderItemDTO struct {
	ID                 uuid.UUID                `json:"id"`
	MenuItemID         uuid.UUID                `json:"menu_item_id"`
	Quantity           int                      `json:"quantity"`
	UnitPrice          float64                  `json:"unit_price"`
	AddOnSubtotal      float64                  `json:"addon_subtotal"`
	LineTotal          float64                  `json:"line_total"`
	CreatedByCustomer  *uuid.UUID               `json:"created_by_customer"`
	Options            []FoodOrderItemOptionDTO `json:"options"`
}

type FoodOrderResponse struct {
	ID           uuid.UUID        `json:"id"`
	RestaurantID uuid.UUID        `json:"restaurant_id"`
	ReservationID *uuid.UUID      `json:"reservation_id"`
	ExpectedTime time.Time        `json:"expected_time"`
	Status       string           `json:"status"`
	Notes        string           `json:"notes"`
	TotalAmount  float64          `json:"total_amount"`
	Channel      string           `json:"channel"`

	Items []FoodOrderItemDTO `json:"items"`
}

type CreateFoodOrderItemOptionReq struct {
	OptionID uuid.UUID `json:"option_id" binding:"required"`
	Qty      int       `json:"qty" binding:"required,min=1"`
}

type CreateFoodOrderItemReq struct {
	MenuItemID uuid.UUID                     `json:"menu_item_id" binding:"required"`
	Quantity   int                           `json:"quantity" binding:"required,min=1"`
	Options    []CreateFoodOrderItemOptionReq `json:"options"`
}

type CreateFoodOrderReq struct {
	Channel      string     `json:"channel" binding:"required,oneof=web walk_in"`
	RestaurantID uuid.UUID  `json:"restaurant_id" binding:"required"`
	ReservationID *uuid.UUID `json:"reservation_id"`
	ExpectedTime time.Time  `json:"expected_time" binding:"required"`
	Notes        string     `json:"notes"`
	Items        []CreateFoodOrderItemReq `json:"items" binding:"required,dive"`
}

type CreateFoodOrderResp struct {
	OrderID      uuid.UUID `json:"order_id"`
	Status       string    `json:"status"`
	TotalAmount  float64   `json:"total_amount"`
	ExpectedTime time.Time `json:"expected_time"`
}


// ===== Append Items to existing order =====
type AppendItemsReq struct {
	Items []CreateFoodOrderItemReq `json:"items" binding:"required,dive"`
}

type AppendItemsResp struct {
	OrderID     uuid.UUID `json:"order_id"`
	AddedCount  int       `json:"added_count"`
	TotalAmount float64   `json:"total_amount"`
}

// ===== Remove single item from order =====
type RemoveItemResp struct {
	OrderID     uuid.UUID `json:"order_id"`
	RemovedItem uuid.UUID `json:"removed_item"`
	TotalAmount float64   `json:"total_amount"`
}

// ===== Attach / List customers in order =====
type AttachCustomerReq struct {
	CustomerID uuid.UUID `json:"customer_id" binding:"required"`
}

type OrderCustomerDTO struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Role       string    `json:"role"` // owner | contributor
}

type ListCustomersResp struct {
	OrderID   uuid.UUID         `json:"order_id"`
	Customers []OrderCustomerDTO `json:"customers"`
}