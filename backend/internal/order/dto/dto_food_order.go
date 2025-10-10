package dto

import "github.com/google/uuid"

type CreateFoodOrderReq struct {
	ReservationID *uuid.UUID            `json:"reservation_id"`           // <-- เพิ่มบรรทัดนี้
	Items         []CreateFoodOrderItem `json:"items" binding:"required,min=1"`
	Note          *string               `json:"note"`
}

type CreateFoodOrderItem struct {
	MenuItemID uuid.UUID               `json:"menu_item_id" binding:"required"`
	CustomerID  *uuid.UUID              `json:"customer_id"`
	Quantity   int                     `json:"quantity" binding:"required,min=1"`
	Note       *string                 `json:"note"`
	Selections []CreateFoodOrderSelect `json:"selections"`
}

type CreateFoodOrderSelect struct {
	GroupID  uuid.UUID `json:"group_id" binding:"required"`
	OptionID uuid.UUID `json:"option_id" binding:"required"`
	Qty      int       `json:"qty"` // ถ้า group.allow_qty=false ฝั่ง BE จะบังคับเป็น 1
}

type CreateFoodOrderResp struct {
	OrderID     uuid.UUID `json:"order_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
}
