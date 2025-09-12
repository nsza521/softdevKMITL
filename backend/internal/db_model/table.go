package models

import (
	"time"

	"github.com/google/uuid"
)

type Table struct {
	Base
	PeopleNum int `gorm:"not null"`
	Row	  string `gorm:"not null"` // A B C ...
	Col	  string `gorm:"not null"` // 1 2 3 ...
}

type Timeslot struct {
	Base
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
}

type TableTimeslot struct {
	Base
	TableID    uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TimeslotID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status     string    `gorm:"not null; default:'available'"` // e.g., "available", "reserved"
}

type TableReservation struct {
	Base
	TableTimeslotID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
