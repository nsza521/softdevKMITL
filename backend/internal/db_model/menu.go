package models

import (
	"time"

	"github.com/google/uuid"
)

// ===== เมนู & หมวด =====

type MenuType struct {
	Base
	Type         string    `gorm:"not null"`
	RestaurantID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// ความสัมพันธ์เดิม
	MenuTags  []MenuTag  `gorm:"foreignKey:MenuTypeID"`
	MenuItems []MenuItem `gorm:"many2many:menu_tags;joinForeignKey:MenuTypeID;joinReferences:MenuItemID"`

	// ใหม่: ผูก Add-on ระดับ "หมวด"
	AddOnGroups []MenuAddOnGroup `gorm:"many2many:menu_type_add_on_groups;joinForeignKey:MenuTypeID;joinReferences:AddOnGroupID"`
}

type MenuItem struct {
	Base
	RestaurantID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name         string    `gorm:"not null"`
	Price        float64   `gorm:"not null"`
	MenuPic      *string   `gorm:"type:text"`
	TimeTaken    int       `gorm:"type:int;default:1"`
	Description  string    `gorm:"type:text"`

	// ความสัมพันธ์กับ Type/Tag เดิม
	MenuTags  []MenuTag  `gorm:"foreignKey:MenuItemID"`
	MenuTypes []MenuType `gorm:"many2many:menu_tags;joinForeignKey:MenuItemID;joinReferences:MenuTypeID"`

	// (คงไว้เป็น "override เฉพาะเมนู") เช่น เพิ่ม/ลด group บางอัน
	AddOnGroups []MenuAddOnGroup `gorm:"many2many:menu_item_add_on_groups;joinForeignKey:MenuItemID;joinReferences:AddOnGroupID"`
}

type MenuTag struct {
	Base
	MenuItemID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuTypeID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuItem   MenuItem  `gorm:"foreignKey:MenuItemID;references:ID"`
	MenuType   MenuType  `gorm:"foreignKey:MenuTypeID;references:ID"`
}

// ===== Add-on =====

// กลุ่ม Add-on (เช่น "เลือกเส้น", "ความเผ็ด")
type MenuAddOnGroup struct {
	Base
	RestaurantID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name         string    `gorm:"not null"`
	Required     bool      `gorm:"not null;default:false"`
	MinSelect    *int      `gorm:"type:int"`
	MaxSelect    *int      `gorm:"type:int"`
	AllowQty     bool      `gorm:"not null;default:false"`

	// ปิดทั้งกลุ่ม (เช่น ของหมดทั้งกลุ่ม)
	IsAvailable      bool       `gorm:"not null;default:true"`
	OutOfStockUntil  *time.Time `gorm:""` // ถ้ามีแผนเปิดอัตโนมัติ
	OutOfStockNote   *string    `gorm:"type:text"`

	Options []MenuAddOnOption `gorm:"foreignKey:GroupID"`
}

// ตัวเลือกในกลุ่ม (เช่น เส้นเล็ก/บะหมี่)
type MenuAddOnOption struct {
	Base
	GroupID          uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name             string    `gorm:"not null"`
	PriceDelta       float64   `gorm:"not null;default:0"`
	IsDefault        bool      `gorm:"not null;default:false"`
	MaxQty           *int      `gorm:"type:int"`

	// ของหมดระดับ "ตัวเลือก"
	IsAvailable      bool       `gorm:"not null;default:true"`
	StockQty         *int       `gorm:"type:int"` // ถ้าอยากตัดสต็อก
	OutOfStockUntil  *time.Time `gorm:""`
	OutOfStockNote   *string    `gorm:"type:text"`
}

// ===== Join Tables =====

type MenuTypeAddOnGroup struct {
	MenuTypeID   uuid.UUID `gorm:"type:char(36);not null;index"`
	AddOnGroupID uuid.UUID `gorm:"type:char(36);not null;index"`
}

type MenuItemAddOnGroup struct {
	MenuItemID   uuid.UUID `gorm:"type:char(36);not null;index"`
	AddOnGroupID uuid.UUID `gorm:"type:char(36);not null;index"`
}
