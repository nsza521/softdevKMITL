package models

import (
	// "time"
	"github.com/google/uuid"
)

type Restaurant struct {
	Base
	Username      string    `gorm:"not null;unique"`
	Password      string    `gorm:"not null"`
	Email         string    `gorm:"not null;unique"`
	// OpenTime      time.Time `gorm:"not null"`
	// CloseTime     time.Time `gorm:"not null"`
	Status 	  	  string    `gorm:"default:'closed'"` // e.g., "open", "closed", "renovation"
	WalletBalance float32   `gorm:"default:0"`
	ProfilePic    *string   `gorm:"type:text"`

	// BankAccountID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Name          *string   `gorm:"type:varchar(255)"`
	MenuType      *string   `gorm:"type:varchar(64)"`
}

type BankAccount struct {
	Base
	UserID  		uuid.UUID 	`gorm:"type:char(36);not null;unique;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BankName      	string 		`gorm:"not null"`
	AccountNumber 	string 		`gorm:"not null"`
	AccountName   	string 		`gorm:"not null"`
	AccountBalance 	float32 	`gorm:"default:0"`
}

type RestaurantAddOn struct {
    ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
    RestaurantID uuid.UUID `gorm:"type:char(36);index"`
    ItemName     string    `gorm:"type:varchar(255);index"`
}