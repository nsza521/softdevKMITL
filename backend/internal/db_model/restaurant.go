package models

import (
	"time"
)

type Restaurant struct {
	Base
	Username		string 		`gorm:"not null;unique"`
	Password		string 		`gorm:"not null"`
	Email			string 		`gorm:"not null;unique"`
	OpenTime		time.Time 	`gorm:"not null"`
	CloseTime		time.Time 	`gorm:"not null"`
	WalletBalance	float32 	`gorm:"default:0"`
	ProfilePic		*string 	`gorm:"not null"`
}