package models

import (
	"time"
	// "github.com/google/uuid"
)

type Table struct {
	Base
	PeopleNum int `gorm:"not null"`
}

type TableReservation struct {
	Base
	TableID       string    `gorm:"foreignKey:TableID;not null;type:char(36)"`
	CustomerID    string    `gorm:"foreignKey:CustomerID;not null;type:char(36)"`
	ReservePeople int       `gorm:"not null"`
	StartTime     time.Time `gorm:"not null"`
	EndTime       time.Time `gorm:"not null"`
	Type          string    `gorm:"not null"` // e.g., "random", "not random"
	Status        string    `gorm:"not null"` // e.g., "pending", "confirmed", "canceled"
}

type TableReservationMembers struct {
	ReservationID 	string `gorm:"foreignKey:ReservationID;not null;type:char(36)"`
	CustomerID   	string `gorm:"foreignKey:CustomerID;not null;type:char(36)"`
}