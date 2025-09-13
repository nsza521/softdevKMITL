// internal/menu/interfaces/usecase.go
package interfaces

import (
	"context"
	"github.com/google/uuid"
)

type MenuTypeBrief struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
}

// เดิมมี MenuItemBrief อยู่แล้ว → เพิ่ม fields "types"
type MenuItemBrief struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Price       float64        `json:"price"`
	MenuPic     *string        `json:"menu_pic"`
	TimeTaken   int            `json:"time_taken"`
	Description string         `json:"description"`
	MenuTypeIDs []uuid.UUID    `json:"menu_type_ids"`
	Types       []MenuTypeBrief`json:"types"` // 👈 รายละเอียด tag ของร้านนั้น
}

type MenuUsecase interface {
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]MenuItemBrief, error)
	CheckRestaurantExists(ctx context.Context, restaurantID uuid.UUID) error
}
