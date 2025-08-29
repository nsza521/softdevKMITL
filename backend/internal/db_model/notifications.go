package models

type Notifications struct {
	Base
	Title	string	`gorm:"not null"`
	Content	string	`gorm:"not null"`
}