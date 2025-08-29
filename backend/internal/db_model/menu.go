package models

import (

	// "github.com/google/uuid"
)

type Menu struct {
	Base
	Type        	string  	`gorm:"not null"`
	RestaurantID 	string 	`gorm:"foreignKey:RestaurantID;not null;type:char(36)"`
	MenuItems  		[]MenuItem 	`gorm:"foreignKey:MenuID;not null"`
}

type MenuItem struct {
	Base
	Name       string    `gorm:"not null"`
	Price      float64   `gorm:"not null"`
	MenuID     string    `gorm:"foreignKey:MenuID;not null;type:char(36)"`
	MenuPic    *string   `gorm:"not null"`
	TimeTaken  int 		 `gorm:"not null"`
	// Description string    `gorm:"not null"`
}

// type MenuItem struct {
// 	Base
// 	RestaurantID int       `gorm:"foreignKey:RestaurantID;not null"`
// 	Name        string    `gorm:"not null"`
// 	// Description string    `gorm:"not null"`
// 	Price      float64   `gorm:"not null"`
// 	MenuPic    *string   `gorm:"not null"`
// 	TimeTaken  time.Time `gorm:"not null"`
// }
