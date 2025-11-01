package dto

type OrderHistoryResponse struct {
	Date   string               `json:"date"`
	Orders []OrderHistoryOrder  `json:"orders"`
}

type OrderHistoryOrder struct {
	OrderID     string               `json:"order_id"`
	Status      string               `json:"status"`
	Channel     string               `json:"channel"`
	Note        string               `json:"note"`
	TotalAmount float64              `json:"total_amount"`
	OrderTime   string               `json:"order_time"`
	Items       []OrderHistoryItem   `json:"items"`
}

type OrderHistoryItem struct {
	MenuName   string                    `json:"menu_name"`
	Quantity   int                       `json:"quantity"`
	UnitPrice  float64                   `json:"unit_price"`
	Subtotal   float64                   `json:"subtotal"`
	Options    []OrderHistoryItemOption  `json:"options"`
}

type OrderHistoryItemOption struct {
	OptionName string  `json:"option_name"`
	PriceDelta float64 `json:"price_delta"`
}
