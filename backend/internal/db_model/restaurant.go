package models
type Restaurant struct {
	Base
	Username  string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"not null;unique"`
	OpenTime  OnlyTime  `gorm:"column:open_time;type:TIME;not null"`
	CloseTime OnlyTime  `gorm:"column:close_time;type:TIME;not null"`
	WalletBalance float32 `gorm:"default:0"`
	ProfilePic *string  `gorm:"not null"`
}
