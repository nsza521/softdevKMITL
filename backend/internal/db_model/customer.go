package models

type Customer struct {
	Base
	Username    	string 		`gorm:"not null;unique"`
	Password    	string 		`gorm:"not null"`
	Email       	string 		`gorm:"not null;unique"`
	FirstName   	string 		`gorm:"not null"`
	LastName    	string 		`gorm:"not null"`
	WalletBalance	float32 	`gorm:"default:0"`
	ProfilePic  	*string 	`gorm:"type:text"`
}