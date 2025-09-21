package models

import (
	"github.com/google/uuid"
)

type TopupHistory struct {
	Base
	UserID    uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount    int    	`gorm:"amount; not null"`
	Status    string 	`gorm:"status; not null"` // e.g., "pending", "completed", "failed"
}
