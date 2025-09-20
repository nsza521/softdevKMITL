package models

import (

	"github.com/google/uuid"
)

type MenuType struct {
	Base
	Type        	string  	`gorm:"not null"`
	RestaurantID 	uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Restaurant  	Restaurant  `gorm:"foreignKey:RestaurantID;references:ID"`
	MenuTags  		[]MenuTag 	`gorm:"foreignKey:MenuTypeID"`
	MenuItems 		[]MenuItem 	`gorm:"many2many:menu_tags;joinForeignKey:MenuTypeID;joinReferences:MenuItemID"`
}

type MenuItem struct {
	Base
	RestaurantID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string     `gorm:"not null"`
	Price       float64    `gorm:"not null"`
	MenuPic     *string    `gorm:"type:text"`
	TimeTaken   int        `gorm:"type:int;default:1"`
	Description string     `gorm:"type:text"`

	// ความสัมพันธ์กับ MenuType เดิม (หมวดหมู่)
	MenuTags  []MenuTag  `gorm:"foreignKey:MenuItemID"`
	MenuTypes []MenuType `gorm:"many2many:menu_tags;joinForeignKey:MenuItemID;joinReferences:MenuTypeID"`

	// ความสัมพันธ์กับ Add-on groups
	AddOnGroups []MenuAddOnGroup `gorm:"many2many:menu_item_addon_groups;joinForeignKey:MenuItemID;joinReferences:AddOnGroupID"`
}

type MenuTag struct {
	Base
	MenuItemID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuTypeID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuItem   MenuItem  `gorm:"foreignKey:MenuItemID;references:ID"`
	MenuType   MenuType  `gorm:"foreignKey:MenuTypeID;references:ID"`
}

// กลุ่ม Add-on (เช่น "เลือกเส้น", "ระดับความเผ็ด")
type MenuAddOnGroup struct {
	Base
	RestaurantID uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Name      string `gorm:"not null"`
	Required  bool   `gorm:"not null;default:false"`
	MinSelect *int   `gorm:"type:int"`
	MaxSelect *int   `gorm:"type:int"`
	AllowQty  bool   `gorm:"not null;default:false"`

	Options []MenuAddOnOption `gorm:"foreignKey:GroupID"`
}

// ตัวเลือกใน Add-on group (เช่น "เส้นเล็ก", "บะหมี่", "ไข่ลวก")
type MenuAddOnOption struct {
	Base
	GroupID    uuid.UUID `gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name       string    `gorm:"not null"`
	PriceDelta float64   `gorm:"not null;default:0"`
	IsDefault  bool      `gorm:"not null;default:false"`
	MaxQty     *int      `gorm:"type:int"`
}

// ตารางกลางเชื่อม MenuItem ↔ AddOnGroup
type MenuItemAddOnGroup struct {
	MenuItemID   uuid.UUID `gorm:"type:char(36);not null;index"`
	AddOnGroupID uuid.UUID `gorm:"type:char(36);not null;index"`
}