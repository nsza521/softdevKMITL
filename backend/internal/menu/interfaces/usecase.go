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
	// MenuTypeIDs []uuid.UUID    `json:"menu_type_ids"`
	Types       []MenuTypeBrief`json:"types"` // 👈 รายละเอียด tag ของร้านนั้น
}

type CreateMenuItemRequest struct {
	Name        string      `json:"name" binding:"required"`
	Price       float64     `json:"price" binding:"required"`
	MenuPic     *string     `json:"menu_pic"`
	TimeTaken   int         `json:"time_taken"`
	Description string      `json:"description"`
	MenuTypeIDs []uuid.UUID `json:"menu_type_ids" binding:"required,min=1"`
}

type UpdateMenuItemRequest struct {
	Name        *string      `json:"name"`
	Price       *float64     `json:"price"`
	MenuPic     **string     `json:"menu_pic"` // null=เคลียร์, ไม่ส่ง=ไม่แตะ
	TimeTaken   *int         `json:"time_taken"`
	Description *string      `json:"description"`
	MenuTypeIDs *[]uuid.UUID `json:"menu_type_ids"` // nil=ไม่แตะ, []=ล้างทั้งหมด
}

type MenuUsecase interface {
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]MenuItemBrief, error)
	CheckRestaurantExists(ctx context.Context, restaurantID uuid.UUID) error

	// สร้าง/แก้ไข/ลบ MenuItem (และผูก/แก้ไข/ลบ MenuType ผ่าน MenuTag)
	CreateMenuItem(ctx context.Context, restaurantID uuid.UUID, in *CreateMenuItemRequest) (*MenuItemBrief, error)
	UpdateMenuItem(ctx context.Context, menuItemID uuid.UUID, in *UpdateMenuItemRequest) (*MenuItemBrief, error)
	DeleteMenuItem(ctx context.Context, menuItemID uuid.UUID) error
}
