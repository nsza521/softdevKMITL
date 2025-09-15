package menu

import "github.com/google/uuid"

// ใช้ตอนสร้าง MenuType ใหม่
type CreateMenuTypeRequest struct {
	Type string `json:"type" binding:"required"`
}

// ใช้ตอนแก้ไข MenuType
// ใช้ pointer เพื่อให้แยกได้ว่า "ไม่ส่ง field" หรือ "ส่งมาเป็นค่าว่าง"
type UpdateMenuTypeRequest struct {
	Type *string `json:"type,omitempty"`
}

// ใช้ตอบกลับ client
type MenuTypeResponse struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Type         string    `json:"type"`
}
