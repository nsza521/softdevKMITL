package menu

import "github.com/google/uuid"

// ===== Group DTO =====
type CreateAddOnGroupRequest struct {
	Name      string `json:"name" binding:"required"`
	Required  bool   `json:"required"`
	MinSelect *int   `json:"min_select"`
	MaxSelect *int   `json:"max_select"`
	AllowQty  bool   `json:"allow_qty"`
}

type UpdateAddOnGroupRequest struct {
	Name      *string `json:"name"`
	Required  *bool   `json:"required"`
	MinSelect *int    `json:"min_select"`
	MaxSelect *int    `json:"max_select"`
	AllowQty  *bool   `json:"allow_qty"`
}

type AddOnGroupResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Required   bool      `json:"required"`
	MinSelect  *int      `json:"min_select"`
	MaxSelect  *int      `json:"max_select"`
	AllowQty   bool      `json:"allow_qty"`
	Restaurant uuid.UUID `json:"restaurant_id"`
}

// ===== Option DTO =====
type CreateAddOnOptionRequest struct {
	Name       string   `json:"name" binding:"required"`
	PriceDelta float64  `json:"price_delta"`
	IsDefault  bool     `json:"is_default"`
	MaxQty     *int     `json:"max_qty"`
}

type UpdateAddOnOptionRequest struct {
	Name       *string  `json:"name"`
	PriceDelta *float64 `json:"price_delta"`
	IsDefault  *bool    `json:"is_default"`
	MaxQty     *int     `json:"max_qty"`
}

type AddOnOptionResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceDelta float64   `json:"price_delta"`
	IsDefault  bool      `json:"is_default"`
	MaxQty     *int      `json:"max_qty"`
	GroupID    uuid.UUID `json:"group_id"`
}
