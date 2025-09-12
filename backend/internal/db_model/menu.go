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
	Name       string    `gorm:"not null"`
	Price      float64   `gorm:"not null"`
	MenuPic    *string   `gorm:"type:text"`
	TimeTaken  int 		 `gorm:"type:int;default:1"` // in minutes
	MenuTags   []MenuTag  `gorm:"foreignKey:MenuItemID"`
	MenuTypes  []MenuType `gorm:"many2many:menu_tags;joinForeignKey:MenuItemID;joinReferences:MenuTypeID"`
	Description string    `gorm:"type:text"`
}

type MenuTag struct {
	Base
	MenuItemID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuTypeID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuItem   MenuItem  `gorm:"foreignKey:MenuItemID;references:ID"`
	MenuType   MenuType  `gorm:"foreignKey:MenuTypeID;references:ID"`
}
