package models

import (
	"time"
	"github.com/google/uuid"
)

type Table struct {
	Base
	PeopleNum int `gorm:"not null"`
}

type TimeSlot struct {
	Base
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
}

type TableTimeSlot struct {
	Base
	TableID    uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TimeSlotID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status    string    `gorm:"not null; default:'available'"` // e.g., "available", "reserved"
}

type TableReservation struct {
	Base
	TableTimeSlotID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID      uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReservePeople   int       `gorm:"not null"`
	// StartTime       time.Time `gorm:"not null"`
	// EndTime         time.Time `gorm:"not null"`
	Type            string    `gorm:"not null"` // e.g., "random", "not random"
	Status          string    `gorm:"not null"` // e.g., "pending", "confirmed", "canceled"
}

type TableReservationMembers struct {
	Base
	ReservationID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CustomerID    uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
