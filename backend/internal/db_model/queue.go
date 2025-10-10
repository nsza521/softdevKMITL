package models

import (

	"github.com/google/uuid"
)

type ReservaTableQueue struct {
	Base
	CustomerID    uuid.UUID `gorm:"type:uuid;not null"`
	Customer      Customer  `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReservePeople int       `gorm:"not null"`
	Status        string    `gorm:"type:varchar(20);not null;default:'pending'"` // e.g., "pending", "seated", "canceled"
}

type FoodOrderQueue struct {
	Base
	CustomerID uuid.UUID `gorm:"type:uuid;not null"`
	Customer   Customer  `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MenuID     uuid.UUID `gorm:"type:uuid;not null"`
	Menu       MenuItem      `gorm:"foreignKey:MenuID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Quantity   int       `gorm:"not null;default:1"`
	Status     string    `gorm:"type:varchar(20);not null;default:'pending'"` // e.g., "pending", "preparing", "ready", "served", "canceled"
}