package dto

import "time"

type QueueItemOptionDTO struct {
	OptionName string  `json:"option_name"`
	GroupName  string  `json:"group_name"`
	PriceDelta float64 `json:"price_delta"`
	Qty        int     `json:"qty"`
}

type QueueItemDTO struct {
	ID           string               `json:"id"`
	MenuName     string               `json:"menu_name"`
	Quantity     int                  `json:"quantity"`
	UnitPrice    float64              `json:"unit_price"`
	Subtotal     float64              `json:"subtotal"`
	TimeTakenMin int                  `json:"time_taken_min"`
	Note         *string              `json:"note"`
	Options      []QueueItemOptionDTO `json:"options"`
}

type QueueOrderDTO struct {
	ID              string         `json:"id"`
	Status          string         `json:"status"`
	Channel         string         `json:"channel"`
	OrderDate       time.Time      `json:"order_date"`
	ExpectedReceive *time.Time     `json:"expected_receive"`
	TotalAmount     float64        `json:"total_amount"`
	Note            *string        `json:"note"`

	Items []QueueItemDTO `json:"items"`
}

type QueueResponse struct {
	Orders []QueueOrderDTO `json:"orders"`
}
