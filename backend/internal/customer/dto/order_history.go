package dto

type OrderHistoryDay struct {
	Date   string             `json:"date"`
	Orders []OrderHistoryItem `json:"orders"`
}

type OrderHistoryItem struct {
	OrderID     string                 `json:"order_id"`
	Channel     string                 `json:"channel"`
	Note        string                 `json:"note"`
	TotalAmount float64                  `json:"total_amount"`
	OrderTime   string                 `json:"order_time"`
	Items       []OrderHistoryLineItem `json:"items"`
}

type OrderHistoryLineItem struct {
	MenuName  string                        `json:"menu_name"`
	Quantity  int                           `json:"quantity"`
	UnitPrice float64                         `json:"unit_price"`
	Subtotal  float64                         `json:"subtotal"`
	Options   []OrderHistoryLineItemOption  `json:"options"`
}

type OrderHistoryLineItemOption struct {
	OptionName string `json:"option_name"`
	PriceDelta float64  `json:"price_delta"`
}
