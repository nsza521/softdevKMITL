package menu

import "github.com/google/uuid"

type MenuTypeBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type AddOnOptionDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceDelta float64   `json:"price_delta"`
	IsDefault  bool      `json:"is_default"`
	MaxQty     *int      `json:"max_qty"`
}

type AddOnGroupDTO struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Required  bool            `json:"required"`
	MinSelect *int            `json:"min_select"`
	MaxSelect *int            `json:"max_select"`
	AllowQty  bool            `json:"allow_qty"`
	// แหล่งที่มา (ช่วย debug)
	From string `json:"from"` // "type" | "item" | "merged"
	// ตัวเลือก
	Options []AddOnOptionDTO `json:"options"`
}

type MenuItemDetailResponse struct {
	ID           uuid.UUID      `json:"id"`
	RestaurantID uuid.UUID      `json:"restaurant_id"`
	Name         string         `json:"name"`
	Price        float64        `json:"price"`
	MenuPic      *string        `json:"menu_pic"`
	TimeTaken    int            `json:"time_taken"`
	Description  string         `json:"description"`

	Types  []MenuTypeBrief `json:"types"`
	AddOns []AddOnGroupDTO  `json:"addons"`
}
