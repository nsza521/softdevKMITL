// internal/order/dto/detail_for_restaurant.go
package dto

import (
	"time"
	"github.com/google/uuid"
)

type OrderDetailForRestaurantResp struct {
	OrderID         uuid.UUID          `json:"order_id"`
	Status          string             `json:"status"`
	OrderDate       time.Time          `json:"order_date"`
	ExpectedReceive *time.Time         `json:"expected_receive_time,omitempty"`
	ReservationID   *uuid.UUID         `json:"reservation_id,omitempty"`
	Note            *string            `json:"note,omitempty"`

	// 👇 เพิ่มตามสคีมา table_reservations
	TableLabel    *string    `json:"table_label,omitempty"`     // เช่น "A5" จาก Row+Col
	TimeslotStart *time.Time `json:"timeslot_start,omitempty"`  // จาก timeslots.start_time
	TimeslotEnd   *time.Time `json:"timeslot_end,omitempty"`    // จาก timeslots.end_time

	Items []OrderKitchenItem `json:"items"`
}

type OrderKitchenItem struct {
	OrderItemID   uuid.UUID                 `json:"order_item_id"`
	MenuItemID    uuid.UUID                 `json:"menu_item_id"`
	MenuName      string                    `json:"menu_name"`
	Quantity      int                       `json:"quantity"`
	UnitPrice     float64                   `json:"unit_price"`     // snapshot
	LineSubtotal  float64                   `json:"line_subtotal"`  // snapshot
	Note          *string                   `json:"note,omitempty"` // เช่น “ไม่ผัก”
	Options       []OrderKitchenItemOption  `json:"options"`
}

type OrderKitchenItemOption struct {
	GroupName   string  `json:"group_name"`
	OptionName  string  `json:"option_name"`
	Qty         int     `json:"qty"`
}
